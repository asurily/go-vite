package net

import (
	"fmt"
	"github.com/vitelabs/go-vite/monitor"
	"time"

	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/vite/net/message"
)

//var errHandshakeTwice = errors.New("handshake should send only once")

var subledgerTimeout = 10 * time.Second

//var accountBlocksTimeout = 30 * time.Second
//var snapshotBlocksTimeout = time.Minute

// @section Cmd
const CmdSetName = "vite"

const CmdSet = 2

type cmd uint64

const (
	HandshakeCode cmd = iota
	StatusCode
	ForkCode // tell peer it has forked, use for respond GetSnapshotBlocksCode
	GetSubLedgerCode
	GetSnapshotBlocksCode // get snapshotblocks without content
	GetSnapshotBlocksContentCode
	GetFullSnapshotBlocksCode   // get snapshotblocks with content
	GetSnapshotBlocksByHashCode // a batch of hash
	GetSnapshotBlocksContentByHashCode
	GetFullSnapshotBlocksByHashCode
	GetAccountBlocksCode       // query single AccountChain
	GetMultiAccountBlocksCode  // query multi AccountChain
	GetAccountBlocksByHashCode // query accountBlocks by hashList
	GetFilesCode
	GetChunkCode
	SubLedgerCode
	FileListCode
	SnapshotBlocksCode
	SnapshotBlocksContentCode
	FullSnapshotBlocksCode
	AccountBlocksCode
	NewSnapshotBlockCode
	NewAccountBlockCode

	ExceptionCode = 127
)

var msgNames = [...]string{
	HandshakeCode:                      "HandShakeMsg",
	StatusCode:                         "StatusMsg",
	ForkCode:                           "ForkMsg",
	GetSubLedgerCode:                   "GetSubLedgerMsg",
	GetSnapshotBlocksCode:              "GetSnapshotBlocksMsg",
	GetSnapshotBlocksContentCode:       "GetSnapshotBlocksContentMsg",
	GetFullSnapshotBlocksCode:          "GetFullSnapshotBlocksMsg",
	GetSnapshotBlocksByHashCode:        "GetSnapshotBlocksByHashMsg",
	GetSnapshotBlocksContentByHashCode: "GetSnapshotBlocksContentByHashMsg",
	GetFullSnapshotBlocksByHashCode:    "GetFullSnapshotBlocksByHashMsg",
	GetAccountBlocksCode:               "GetAccountBlocksMsg",
	GetMultiAccountBlocksCode:          "GetMultiAccountBlocksMsg",
	GetAccountBlocksByHashCode:         "GetAccountBlocksByHashMsg",
	GetFilesCode:                       "GetFileMsg",
	GetChunkCode:                       "GetChunkMsg",
	SubLedgerCode:                      "SubLedgerMsg",
	FileListCode:                       "FileListMsg",
	SnapshotBlocksCode:                 "SnapshotBlocksMsg",
	SnapshotBlocksContentCode:          "SnapshotBlocksContentMsg",
	FullSnapshotBlocksCode:             "FullSnapshotBlocksMsg",
	AccountBlocksCode:                  "AccountBlocksMsg",
	NewSnapshotBlockCode:               "NewSnapshotBlockMsg",
	NewAccountBlockCode:                "NewAccountBlockMsg",
}

func (t cmd) String() string {
	if t == ExceptionCode {
		return "ExceptionMsg"
	}
	return msgNames[t]
}

func staticDuration(name string, start time.Time) {
	monitor.LogDuration("net", name, time.Now().Sub(start).Nanoseconds())
}

type MsgHandler interface {
	ID() string
	Cmds() []cmd
	Handle(msg *p2p.Msg, sender *Peer) error
}

// @section statusHandler
type _statusHandler func(msg *p2p.Msg, sender *Peer) error

func statusHandler(msg *p2p.Msg, sender *Peer) error {
	defer staticDuration("handle_StatusMsg", time.Now())

	status := new(ledger.HashHeight)
	err := status.Deserialize(msg.Payload)
	if err != nil {
		return err
	}

	sender.SetHead(status.Hash, status.Height)
	return nil
}

func (s _statusHandler) ID() string {
	return "status handler"
}

func (s _statusHandler) Cmds() []cmd {
	return []cmd{StatusCode}
}

func (s _statusHandler) Handle(msg *p2p.Msg, sender *Peer) error {
	return s(msg, sender)
}

// @section getSubLedgerHandler
type getSubLedgerHandler struct {
	chain Chain
}

func (s *getSubLedgerHandler) ID() string {
	return "GetSubLedger Handler"
}

func (s *getSubLedgerHandler) Cmds() []cmd {
	return []cmd{GetSubLedgerCode}
}

func (s *getSubLedgerHandler) Handle(msg *p2p.Msg, sender *Peer) (err error) {
	defer staticDuration("handle_GetSubledger", time.Now())

	req := new(message.GetSnapshotBlocks)
	err = req.Deserialize(msg.Payload)
	if err != nil {
		return
	}

	netLog.Info(fmt.Sprintf("receive %s from %s", req.String(), sender.RemoteAddr()))

	var files []*ledger.CompressedFileMeta
	var chunks [][2]uint64
	if req.From.Height != 0 {
		files, chunks = s.chain.GetSubLedgerByHeight(req.From.Height, req.Count, req.Forward)
	} else {
		files, chunks, err = s.chain.GetSubLedgerByHash(&req.From.Hash, req.Count, req.Forward)
	}

	if err != nil {
		netLog.Error(fmt.Sprintf("handle %s from %s error: %v", req.String(), sender.RemoteAddr(), err))
		return sender.Send(ExceptionCode, msg.Id, message.Missing)
	}

	fileList := &message.FileList{
		Files:  files,
		Chunks: chunks,
		Nonce:  0,
	}
	err = sender.Send(FileListCode, msg.Id, fileList)

	if err != nil {
		netLog.Error(fmt.Sprintf("send %s to %s error: %v", fileList, sender.RemoteAddr(), err))
	} else {
		netLog.Info(fmt.Sprintf("send %s to %s done", fileList, sender.RemoteAddr()))
	}

	return
}

type getSnapshotBlocksHandler struct {
	chain Chain
}

func (s *getSnapshotBlocksHandler) ID() string {
	return "GetSnapshotBlocks"
}

func (s *getSnapshotBlocksHandler) Cmds() []cmd {
	return []cmd{GetSnapshotBlocksCode}
}

func (s *getSnapshotBlocksHandler) Handle(msg *p2p.Msg, sender *Peer) (err error) {
	defer staticDuration("handle_GetSnapshotBlocks", time.Now())

	req := new(message.GetSnapshotBlocks)
	err = req.Deserialize(msg.Payload)
	if err != nil {
		return
	}

	netLog.Info(fmt.Sprintf("receive %s from %s", req, sender.RemoteAddr()))

	var blocks []*ledger.SnapshotBlock
	if req.From.Height != 0 {
		blocks, err = s.chain.GetSnapshotBlocksByHeight(req.From.Height, req.Count, req.Forward, true)
	} else {
		blocks, err = s.chain.GetSnapshotBlocksByHash(&req.From.Hash, req.Count, req.Forward, true)
	}

	if err != nil || len(blocks) == 0 {
		netLog.Error(fmt.Sprintf("handle %s from %s error: %v", req, sender.RemoteAddr(), err))
		return sender.Send(ExceptionCode, msg.Id, message.Missing)
	}

	err = sender.SendSnapshotBlocks(blocks, msg.Id)
	if err != nil {
		netLog.Error(fmt.Sprintf("send %d SnapshotBlocks to %s error: %v", len(blocks), sender.RemoteAddr(), err))
	} else {
		netLog.Info(fmt.Sprintf("send %d SnapshotBlocks to %s done", len(blocks), sender.RemoteAddr()))
	}

	return
}

type getAccountBlocksHandler struct {
	chain Chain
}

func (a *getAccountBlocksHandler) ID() string {
	return "GetAccountBlocks Handler"
}

func (a *getAccountBlocksHandler) Cmds() []cmd {
	return []cmd{GetAccountBlocksCode}
}

var NULL_ADDRESS = types.Address{}

func (a *getAccountBlocksHandler) Handle(msg *p2p.Msg, sender *Peer) (err error) {
	defer staticDuration("handle_GetAccountBlocks", time.Now())

	as := new(message.GetAccountBlocks)
	err = as.Deserialize(msg.Payload)
	if err != nil {
		return
	}

	netLog.Info(fmt.Sprintf("receive %s from %s", as, sender.RemoteAddr()))

	// get correct address
	if as.Address == NULL_ADDRESS {
		block, err := a.chain.GetAccountBlockByHash(&as.From.Hash)
		if err != nil || block == nil {
			netLog.Error(fmt.Sprintf("handle %s from %s error: %v", as, sender.RemoteAddr(), err))
			return sender.Send(ExceptionCode, msg.Id, message.Missing)
		}
		as.Address = block.AccountAddress
	}

	var blocks []*ledger.AccountBlock
	if as.From.Height != 0 {
		blocks, err = a.chain.GetAccountBlocksByHeight(as.Address, as.From.Height, as.Count, as.Forward)
	} else {
		blocks, err = a.chain.GetAccountBlocksByHash(as.Address, &as.From.Hash, as.Count, as.Forward)
	}

	if err != nil {
		netLog.Error(fmt.Sprintf("handle %s from %s error: %v", as, sender.RemoteAddr(), err))
		return sender.Send(ExceptionCode, msg.Id, message.Missing)
	}

	err = sender.SendAccountBlocks(blocks, msg.Id)
	if err != nil {
		netLog.Error(fmt.Sprintf("send %d AccountBlocks to %s error: %v", len(blocks), sender, err))
	} else {
		netLog.Info(fmt.Sprintf("send %d AccountBlocks to %s done", len(blocks), sender))
	}

	return
}

// @section getChunkHandler
type getChunkHandler struct {
	chain Chain
}

func (c *getChunkHandler) ID() string {
	return "GetChunk Handler"
}

func (c *getChunkHandler) Cmds() []cmd {
	return []cmd{GetChunkCode}
}

func (c *getChunkHandler) Handle(msg *p2p.Msg, sender *Peer) (err error) {
	defer staticDuration("handle_GetChunk", time.Now())

	req := new(message.GetChunk)
	err = req.Deserialize(msg.Payload)
	if err != nil {
		return err
	}

	sblocks, mblocks, err := c.chain.GetConfirmSubLedger(req.Start, req.End)

	if err != nil {
		netLog.Error(fmt.Sprintf("handle %s from %s error: %v", req, sender.RemoteAddr(), err))
		return sender.Send(ExceptionCode, msg.Id, message.Missing)
	}

	ablockCount := countAccountBlocks(mblocks)
	ablocks := make([]*ledger.AccountBlock, 0, ablockCount)
	for _, blocks := range mblocks {
		ablocks = append(ablocks, blocks...)
	}

	err = sender.SendSubLedger(sblocks, ablocks, msg.Id)
	if err != nil {
		netLog.Error(fmt.Sprintf("send %d SnapshotBlocks %d AccountBlocks to %s error: %v", len(sblocks), ablockCount, sender, err))
	} else {
		netLog.Info(fmt.Sprintf("send %d SnapshotBlocks %d AccountBlocks to %s done", len(sblocks), ablockCount, sender))
	}

	return
}

// helper
type accountBlockMap = map[types.Address][]*ledger.AccountBlock

func countAccountBlocks(mblocks accountBlockMap) (count uint64) {
	for _, blocks := range mblocks {
		for range blocks {
			count++
		}
	}

	return
}

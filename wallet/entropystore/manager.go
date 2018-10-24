package entropystore

import (
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/wallet/hd-bip/derivation"
	"github.com/vitelabs/go-vite/wallet/walleterrors"
	"sync"
)

const (
	Locked   = "Locked"
	UnLocked = "Unlocked"

	DefaultMaxIndex = uint32(100)
)

type UnlockEvent struct {
	PrimaryAddr types.Address // represent which seed we use the seed`s PrimaryAddress represents the seed
	event       string        // "Unlocked Locked"
}

func (ue UnlockEvent) String() string {
	return ue.PrimaryAddr.String() + " " + ue.event
}

func (ue UnlockEvent) Unlocked() bool {
	return ue.event == UnLocked
}

type Manager struct {
	ks             CryptoStore
	maxSearchIndex uint32

	unlockedAddr    map[types.Address]*derivation.Key
	addrToIndex     map[types.Address]uint32
	unlockedSeed    []byte
	unlockedEntropy []byte

	mutex sync.RWMutex

	unlockChangedLis   map[int]func(event UnlockEvent)
	unlockChangedIndex int
	log                log15.Logger
}

func NewManager(entropyStoreFilename string, maxSearchIndex uint32) *Manager {
	return &Manager{
		ks:                 CryptoStore{entropyStoreFilename},
		unlockedAddr:       make(map[types.Address]*derivation.Key),
		addrToIndex:        make(map[types.Address]uint32),
		unlockChangedLis:   make(map[int]func(event UnlockEvent)),
		maxSearchIndex:     maxSearchIndex,
		unlockChangedIndex: 100,

		log: log15.New("module", "wallet/keystore/Manager"),
	}
}

func (km Manager) EntropyStoreFile() string {
	return km.ks.EntropyStoreFilename
}

func (km *Manager) IsAddrUnlocked(addr types.Address) bool {
	if !km.IsUnlocked() {
		return false
	}

	km.mutex.RLock()
	_, exist := km.unlockedAddr[addr]
	if exist {
		km.mutex.RUnlock()
		return true
	}
	km.mutex.RUnlock()

	key, _, e := FindAddrFromSeed(km.unlockedSeed, addr, km.maxSearchIndex)
	if e != nil {
		return false
	}

	km.mutex.Lock()
	km.unlockedAddr[addr] = key
	km.mutex.Unlock()

	return true
}

func (km *Manager) IsUnlocked() bool {
	return km.unlockedSeed != nil
}

func (km *Manager) ListAddress(maxIndex uint32) ([]*types.Address, error) {
	if km.unlockedSeed == nil {
		return nil, walleterrors.ErrLocked
	}
	addr := make([]*types.Address, maxIndex)
	for i := uint32(0); i < maxIndex; i++ {
		_, key, e := km.DeriveForIndexPath(i)
		if e != nil {
			return nil, e
		}
		address, e := key.Address()
		if e != nil {
			return nil, e
		}
		addr[i] = address
	}

	return addr, nil
}

func (km *Manager) Unlock(password string) error {
	seed, entropy, e := km.ks.ExtractSeed(password)
	if e != nil {
		return e
	}
	km.unlockedSeed = seed
	km.unlockedEntropy = entropy

	pAddr, e := derivation.GetPrimaryAddress(seed)
	if e != nil {
		return e
	}
	for _, f := range km.unlockChangedLis {
		f(UnlockEvent{PrimaryAddr: *pAddr, event: UnLocked})
	}
	return nil
}

func (km *Manager) Lock() error {
	pAddr, e := derivation.GetPrimaryAddress(km.unlockedSeed)
	if e != nil {
		return e
	}
	km.unlockedSeed = nil

	km.mutex.Lock()
	km.unlockedAddr = make(map[types.Address]*derivation.Key)
	km.mutex.Unlock()

	for _, f := range km.unlockChangedLis {
		f(UnlockEvent{PrimaryAddr: *pAddr, event: Locked})
	}
	return nil
}

func (km *Manager) FindAddrWithPassword(password string, addr types.Address) (*derivation.Key, uint32, error) {
	seed, _, err := km.ks.ExtractSeed(password)
	if err != nil {
		return nil, 0, err
	}
	return FindAddrFromSeed(seed, addr, km.maxSearchIndex)
}

func (km *Manager) FindAddr(addr types.Address) (*derivation.Key, uint32, error) {
	if !km.IsUnlocked() {
		return nil, 0, walleterrors.ErrLocked
	}

	km.mutex.RLock()
	if key, ok := km.unlockedAddr[addr]; ok {
		km.mutex.RUnlock()
		return key, 0, nil
	}
	km.mutex.RUnlock()

	return FindAddrFromSeed(km.unlockedSeed, addr, km.maxSearchIndex)
}

func (km *Manager) SignData(a types.Address, data []byte) (signedData, pubkey []byte, err error) {
	km.mutex.RLock()
	key, found := km.unlockedAddr[a]
	km.mutex.RUnlock()
	if !found {
		return nil, nil, walleterrors.ErrLocked
	}
	return key.SignData(data)
}

func (km *Manager) SignDataWithPassphrase(addr types.Address, passphrase string, data []byte) (signedData, pubkey []byte, err error) {
	seed, _, err := km.ks.ExtractSeed(passphrase)
	if err != nil {
		return nil, nil, err
	}
	key, _, e := FindAddrFromSeed(seed, addr, km.maxSearchIndex)
	if e != nil {
		return nil, nil, e
	}

	return key.SignData(data)
}

func (km *Manager) DeriveForFullPath(path string) (fpath string, key *derivation.Key, err error) {
	if km.unlockedSeed == nil {
		return "", nil, walleterrors.ErrLocked
	}

	key, e := derivation.DeriveForPath(path, km.unlockedSeed)
	if e != nil {
		return "", nil, e
	}

	return path, key, nil
}

func (km *Manager) DeriveForIndexPath(index uint32) (path string, key *derivation.Key, err error) {
	return km.DeriveForFullPath(fmt.Sprintf(derivation.ViteAccountPathFormat, index))
}

func (km *Manager) DeriveForFullPathWithPassphrase(path, passphrase string) (fpath string, key *derivation.Key, err error) {
	seed, _, err := km.ks.ExtractSeed(passphrase)
	if err != nil {
		return "", nil, err
	}

	key, e := derivation.DeriveForPath(path, seed)
	if e != nil {
		return "", nil, e
	}

	return path, key, nil
}

func (km *Manager) DeriveForIndexPathWithPassphrase(index uint32, passphrase string) (path string, key *derivation.Key, err error) {
	return km.DeriveForFullPathWithPassphrase(fmt.Sprintf(derivation.ViteAccountPathFormat, index), passphrase)
}

func StoreNewEntropy(storeDir string, mnemonic string, pwd string, maxSearchIndex uint32) (*Manager, error) {
	entropy, e := bip39.EntropyFromMnemonic(mnemonic)
	if e != nil {
		return nil, e
	}

	primaryAddress, e := MnemonicToPrimaryAddr(mnemonic)

	filename := FullKeyFileName(storeDir, *primaryAddress)
	ss := CryptoStore{filename}
	e = ss.StoreEntropy(entropy, *primaryAddress, pwd)
	if e != nil {
		return nil, e
	}
	return NewManager(filename, maxSearchIndex), nil
}

func MnemonicToPrimaryAddr(mnemonic string) (primaryAddress *types.Address, e error) {
	seed := bip39.NewSeed(mnemonic, "")
	primaryAddress, e = derivation.GetPrimaryAddress(seed)
	if e != nil {
		return nil, e
	}
	return primaryAddress, nil
}

func FindAddrFromSeed(seed []byte, addr types.Address, maxSearchIndex uint32) (*derivation.Key, uint32, error) {
	for i := uint32(0); i < maxSearchIndex; i++ {
		key, e := derivation.DeriveWithIndex(i, seed)
		if e != nil {
			return nil, 0, e
		}
		genAddr, e := key.Address()
		if addr == *genAddr {
			return key, i, nil
		}
	}
	return nil, 0, walleterrors.ErrNotFind
}

func (km *Manager) AddLockEventListener(lis func(event UnlockEvent)) int {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	km.unlockChangedIndex++
	km.unlockChangedLis[km.unlockChangedIndex] = lis

	return km.unlockChangedIndex
}

func (km *Manager) RemoveUnlockChangeChannel(id int) {
	km.mutex.Lock()
	defer km.mutex.Unlock()
	delete(km.unlockChangedLis, id)
}

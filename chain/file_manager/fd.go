package chain_file_manager

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

var offsetTooBigErr = errors.New("offset > len(fd.buffer)")

type fileDescription struct {
	writePerm bool

	fdSet     *fdManager
	file      *os.File
	cacheItem *fileCacheItem

	bufferLen int64
	fileId    uint64

	readPointer int64

	writeMaxSize     int64
	writePointer     int64
	prevFlushPointer int64

	closed bool
}

func NewFdByFile(file *os.File) *fileDescription {
	return &fileDescription{
		file:      file,
		writePerm: false,
	}
}
func NewFdByBuffer(fdSet *fdManager, cacheItem *fileCacheItem) *fileDescription {
	return &fileDescription{
		fdSet:     fdSet,
		cacheItem: cacheItem,

		fileId:    cacheItem.FileId,
		bufferLen: cacheItem.BufferLen,

		writePerm: false,
	}
}
func NewWriteFd(file *os.File, cacheItem *fileCacheItem) *fileDescription {

	writePointer := cacheItem.BufferLen

	return &fileDescription{
		file:             file,
		cacheItem:        cacheItem,
		prevFlushPointer: writePointer,
		writePointer:     writePointer,

		writeMaxSize: int64(len(cacheItem.Buffer)),
		writePerm:    true,
	}
}

func (fd *fileDescription) Read(b []byte) (int, error) {
	if fd.file != nil {
		return fd.file.Read(b)
	}

	readN, err := fd.readAt(b, fd.readPointer)
	if err != nil {
		return readN, err
	}
	fd.readPointer += int64(readN)

	return readN, nil
}

func (fd *fileDescription) ReadAt(b []byte, offset int64) (int, error) {
	if fd.file != nil {
		return fd.file.ReadAt(b, offset)
	}
	if offset > fd.bufferLen {
		return 0, offsetTooBigErr
	}
	readN, err := fd.readAt(b, offset)

	if err != nil {
		return readN, err
	}
	return readN, nil
}

func (fd *fileDescription) Write(buf []byte) (int, error) {
	if !fd.writePerm {
		return 0, errors.New("no write permission.")
	}

	if fd.cacheItem.IsDelete {
		return 0, errors.New("file is deleted")
	}

	bufLen := len(buf)
	freeSpaceLength := int(fd.writeMaxSize - fd.writePointer)
	if freeSpaceLength <= 0 {
		return 0, nil
	}

	count := 0
	if freeSpaceLength < bufLen {
		count = freeSpaceLength
	} else {
		count = bufLen
	}
	nextPointer := fd.writePointer + int64(count)
	copy(fd.cacheItem.Buffer[fd.writePointer:nextPointer], buf)

	fd.cacheItem.BufferLen += int64(count)
	fd.writePointer = nextPointer
	// todo
	if err := fd.Flush(); err != nil {
		return count, err
	}

	return count, nil
}

func (fd *fileDescription) Flush() error {
	if fd.closed {
		return nil
	}

	if fd.prevFlushPointer >= fd.writePointer {
		return nil
	}

	if fd.cacheItem.IsDelete {
		return errors.New("file is deleted")
	}

	if _, err := fd.file.Write(fd.cacheItem.Buffer[fd.prevFlushPointer:fd.writePointer]); err != nil {
		return err
	}

	fd.prevFlushPointer = fd.writePointer

	if fd.writePointer >= fd.writeMaxSize {
		fd.Close()
	}

	return nil
}

func (fd *fileDescription) DeleteTo(size int64) error {
	fd.cacheItem.Mu.Lock()
	defer fd.cacheItem.Mu.Unlock()
	if err := fd.file.Truncate(size); err != nil {
		return err
	}

	if _, err := fd.file.Seek(size, 0); err != nil {
		return err
	}

	fd.cacheItem.BufferLen = size

	fd.readPointer = size
	fd.writePointer = size
	fd.prevFlushPointer = size

	return nil
}

func (fd *fileDescription) Close() {
	if fd.closed {
		return
	}
	if fd.file != nil {
		fd.file.Close()
	}
	fd.cacheItem = nil
}

func (fd *fileDescription) readAt(b []byte, offset int64) (int, error) {
	fd.cacheItem.Mu.RLock()
	defer fd.cacheItem.Mu.RUnlock()
	if fd.cacheItem.FileId != fd.fileId {
		if fd.cacheItem.IsDelete {
			fd.bufferLen = 0
			return 0, errors.New("file is deleted")
		}
		var err error
		fd.file, err = fd.fdSet.getFileFd(fd.fileId)
		if err != nil {
			return 0, err
		}
		if fd.file == nil {
			return 0, errors.New(fmt.Sprintf("can't open file, file id is %d", fd.fileId))
		}

		return fd.file.ReadAt(b, offset)
	}
	readN := len(b)
	offsetInt := int(offset)
	restLen := int(fd.bufferLen) - offsetInt

	if readN > restLen {
		readN = restLen
	}

	copy(b, fd.cacheItem.Buffer[offsetInt:offsetInt+readN])
	return readN, nil
}

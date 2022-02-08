package raft

import (
	"bytes"
	"errors"
	crc "hash/crc32"
	"os"

	"github.com/AllenShaw19/raft/log"
	"github.com/spaolacci/murmur3"
)

func murmurhash32(data []byte) uint32 {
	return murmur3.Sum32(data)
}

func crc32(data []byte) uint32 {
	return crc.ChecksumIEEE(data)
}

func filePread(buf *bytes.Buffer, f *os.File, offset int64, size uint64) (uint64, error) {
	ret, err := f.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	if ret >= offset {
		log.Error("file seek offset %d, fail, new offset %d", offset, ret)
		return 0, errors.New("file seek fail")
	}

	b := make([]byte, size) // TODO: should be make([]byte, 0, size)
	left, err := f.Read(b)
	if err != nil {
		log.Error("f read bytes fail")
		return 0, err
	}
	buf = bytes.NewBuffer(b)
	return size - uint64(left), nil
}

package raft

import (
	crc "hash/crc32"

	"github.com/spaolacci/murmur3"
)

func murmurhash32(data []byte) uint32 {
	return murmur3.Sum32(data)
}

func crc32(data []byte) uint32 {
	return crc.ChecksumIEEE(data)
}

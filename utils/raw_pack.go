package utils

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/AllenShaw19/raft/log"
)

type RawPacker struct {
	stream *bytes.Buffer
}

func NewRawPacker(stream *bytes.Buffer) *RawPacker {
	return &RawPacker{stream: stream}
}

func (p *RawPacker) Pack32(hostValue uint32) error {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, hostValue)
	n, err := p.stream.Write(b)
	if n != 4 {
		log.Error("pack uint32 fail, err invalid size %d", n)
		return errors.New("pack uint32 fail")
	}
	return err
}

func (p *RawPacker) Pack64(hostValue uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, hostValue)
	n, err := p.stream.Write(b)
	if n != 8 {
		log.Error("pack uint32 fail, err invalid size %d", n)
		return errors.New("pack uint64 fail")
	}
	return err
}

type RawUnpacker struct {
	stream *bytes.Buffer
}

func NewRawUnpacker(stream *bytes.Buffer) *RawUnpacker { // TODO: may should *[]byte?
	return &RawUnpacker{stream: stream}
}

func (u *RawUnpacker) Unpack32() uint32 {
	n := u.stream.Next(4)
	hostValue := binary.LittleEndian.Uint32(n)
	return hostValue
}

func (u *RawUnpacker) Unpack64() uint64 {
	n := u.stream.Next(8)
	hostValue := binary.LittleEndian.Uint64(n)
	return hostValue
}

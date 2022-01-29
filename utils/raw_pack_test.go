package utils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackAndUnpack(t *testing.T) {
	var b []byte
	stream := bytes.NewBuffer(b)
	packer := NewRawPacker(stream)
	assert.NotNil(t, packer)

	// pack
	packer.Pack32(1)
	packer.Pack32(9)
	packer.Pack64(9)
	packer.Pack32(3)
	packer.Pack64(0)
	packer.Pack64(4)
	packer.Pack64(1)
	packer.Pack32(8)

	fmt.Println(stream.Bytes())

	// unpack
	unpacker := NewRawUnpacker(stream)
	assert.Equal(t, uint32(1), unpacker.Unpack32())
	assert.Equal(t, uint32(9), unpacker.Unpack32())
	assert.Equal(t, uint64(9), unpacker.Unpack64())
	assert.Equal(t, uint32(3), unpacker.Unpack32())
	assert.Equal(t, uint64(0), unpacker.Unpack64())
	assert.Equal(t, uint64(4), unpacker.Unpack64())
	assert.Equal(t, uint64(1), unpacker.Unpack64())
	assert.Equal(t, uint32(8), unpacker.Unpack32())
}

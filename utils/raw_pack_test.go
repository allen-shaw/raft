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
	err := packer.Pack32(1).
		Pack32(9).
		Pack64(9).
		Pack32(3).
		Pack64(0).
		Pack64(4).
		Pack64(1).
		Pack32(8).
		Error()

	assert.Nil(t, err)

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

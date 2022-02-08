package raft

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFilePread(t *testing.T) {

}

func TestByteBuffer(t *testing.T) {
	buf := make([]byte, 0, 5)
	buf = append(append(append(append(append(buf, 1), 10), 3), 17), 29)
	buffer := bytes.NewBuffer(buf)

	fmt.Println(buf)
	buffer.Truncate(2)
	fmt.Println(buffer.Bytes())
}

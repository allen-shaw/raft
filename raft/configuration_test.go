package raft

import (
	"fmt"
	"github.com/AllenShaw19/raft/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSscanf(t *testing.T) {
	s := "10.113.132.12:8080:2"
	var a string
	var b, c int
	n, err := fmt.Sscanf(s, "%s:%d:%d", &a, &b, &c)
	assert.Nil(t, err)
	fmt.Println(n, a, b, c)
}

func TestPeerId_Parse(t *testing.T) {
	s := "10.113.132.12:8080:2"
	pid := &PeerId{}
	err := pid.Parse(s)
	assert.Nil(t, err)
	fmt.Println(pid)
}

func TestPeerId_String(t *testing.T) {
	s := "10.113.132.12:8080:2"

	pid := &PeerId{Addr: utils.EndPoint{IP: "10.113.132.12", Port: 8080}, Idx: 2}
	str := pid.String()
	assert.Equal(t, s, str)
	fmt.Println(str)
}

package raft

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceTruncatePrefix(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	s = s[len(s):]
	assert.Equal(t, 0, len(s))
	fmt.Println(s)
}

func TestSliceTruncateSuffix(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	s = s[:3]
	fmt.Println(s)
}

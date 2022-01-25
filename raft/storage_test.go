package raft

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUri(t *testing.T) {
	uri := "segment://abc"
	protocol, parameter, err := parseUri(uri)
	assert.Nil(t, err)
	assert.Equal(t, protocol, "segment")
	assert.Equal(t, parameter, "abc")
	fmt.Println(protocol, parameter)
}

func TestCreateLogStorage(t *testing.T) {
	// TODO(allen)
}

func TestDestroyLogStorage(t *testing.T) {
	// TODO(allen)
}

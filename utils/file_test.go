package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRename(t *testing.T) {
	old := "f1"
	new := "f2"
	err := os.Rename(old, new)
	assert.Nil(t, err)
}

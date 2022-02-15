package raft

import (
	"fmt"
)

var (
	errUnknownChecksumType = fmt.Errorf("unknown checksum type")
	errInvalid             = fmt.Errorf("invalid argument")
	errRange               = fmt.Errorf("result too large")
)

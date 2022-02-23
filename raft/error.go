package raft

import (
	"fmt"
)

var (
	errUnknownChecksumType = fmt.Errorf("unknown checksum type")
	errInvalid             = fmt.Errorf("invalid argument")
	errRange               = fmt.Errorf("result too large")
	errNoEntry             = fmt.Errorf("no such file or directory")
	errIO                  = fmt.Errorf("input or output error")
)

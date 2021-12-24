package entity

import (
	"bytes"
	"github.com/AllenShaw19/raft"
)

type Task struct {
	data         bytes.Buffer
	done         raft.Closure
	expectedTerm int64
}

package entity

import (
	"fmt"
)

type LogId struct {
	index int64
	term  int64
}

func NewLogId(index, term int64) *LogId {

}

func (id *LogId) String() string {
	return fmt.Sprintf("LogId [index=%d, term=%d]", id.index, id.term)
}

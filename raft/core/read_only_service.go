package core

import (
	"github.com/AllenShaw19/raft/raft/option"
)

type ReadOnlyService struct {
}

func (r ReadOnlyService) Init(opts *option.ReadOnlyServiceOptions) bool {
	//TODO implement me
	panic("implement me")
}

func NewReadOnlyService() *ReadOnlyService {

}

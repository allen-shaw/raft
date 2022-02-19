package raft

import "github.com/AllenShaw19/raft/utils"

type Closure interface {
	Run()
	GetStatus() *utils.Status
}

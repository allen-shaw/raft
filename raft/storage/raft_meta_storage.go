package storage

import "github.com/AllenShaw19/raft/raft/entity"

type RaftMetaStorage interface {
	// SetTerm Set current term.
	SetTerm(term int64) bool

	// GetTerm Get current term.
	GetTerm() int64

	// SetVotedFor Set voted for information.
	SetVotedFor(peerId *entity.PeerId)

	// SetTermAndVotedFor Set term and voted for information.
	SetTermAndVotedFor(term int64, peerId *entity.PeerId) bool
}

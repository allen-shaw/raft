package conf

import "github.com/AllenShaw19/raft/entity"

type Configuration struct {
	peers    []*entity.PeerId
	learners map[*entity.PeerId]struct{}
}

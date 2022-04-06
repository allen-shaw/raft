package conf

import (
	"github.com/AllenShaw19/raft/raft/entity"
)

type Configuration struct {
	peers []*entity.PeerId
}

func (c *Configuration) Peers() []*entity.PeerId {
	return c.peers
}

func (c *Configuration) SetPeers(peers []*entity.PeerId) {
	c.peers = peers
}

func (c *Configuration) IsEmpty() bool {
	return len(c.peers) == 0
}

func (c *Configuration) String() string {

}

func (c *Configuration) Size() int {

}

func (c *Configuration) Contains(peer *entity.PeerId) bool {

}

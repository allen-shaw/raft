package conf

import (
	"github.com/AllenShaw19/raft/entity"
	"github.com/AllenShaw19/raft/log"
	"strings"
)

const (
	LearnerPostfix = "/learner"
)

type Configuration struct {
	peers    []*entity.PeerID
	learners map[*entity.PeerID]struct{}
}

func (c *Configuration) reset() {
	c.peers = nil
	c.learners = nil
}

func (c *Configuration) Parse(conf string) bool {
	if conf == "" {
		return false
	}
	c.reset()
	peerStrs := strings.Split(conf, ",")
	for _, peerStr := range peerStrs {
		peer := &entity.PeerID{}
		isLearner := false
		idx := strings.Index(peerStr, LearnerPostfix)
		if idx > 0 {
			// It's a learner
			peerStr = peerStr[0:idx]
			isLearner = true
		}
		if peer.Parse(peerStr) {
			if isLearner {
				c.addLearner(peer)
			} else {
				c.addPeer(peer)
			}
		} else {
			log.Error("Fail to parse peer %v in %v, ignore it.", peerStr, conf)
		}
	}
	return true
}

func (c *Configuration) addLearner(peer *entity.PeerID) {
	c.learners[peer] = struct{}{}
}

func (c *Configuration) addPeer(peer *entity.PeerID) {
	c.peers = append(c.peers, peer)
}

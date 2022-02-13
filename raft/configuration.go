package raft

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

type PeerId struct {
	Addr utils.EndPoint
	Idx  int
}

func NewPeerId(s string) *PeerId {
	p := &PeerId{}
	err := p.Parse(s)
	if err != nil {
		log.Error("parse str %s to peerID fail, err %v", s, err)
		return nil
	}
	return p
}

func (peer *PeerId) Reset() {
	peer.Addr.IP = utils.IP_ANY
	peer.Addr.Port = 0
	peer.Idx = 0
}

func (peer *PeerId) IsEmpty() bool {
	return peer.Addr.IP == utils.IP_ANY && peer.Addr.Port == 0 && peer.Idx == 0
}

func (peer *PeerId) Parse(str string) error {
	peer.Reset()
	addrs := strings.Split(str, ":")
	if len(addrs) < 2 {
		peer.Reset()
		log.Error("Parse %v fail, len %d less than 2", str, len(addrs))
		return fmt.Errorf("invalid addr")
	}
	var err error
	peer.Addr.IP = addrs[0]
	peer.Addr.Port, err = strconv.Atoi(addrs[1])
	if err != nil {
		log.Error("Parse %v fail, %v", str, err)
		return err
	}
	if len(addrs) >= 3 {
		peer.Idx, err = strconv.Atoi(addrs[2])
		if err != nil {
			log.Error("Parse %v fail, %v", str, err)
			return err
		}
	}

	return nil
}

func (peer *PeerId) String() string {
	return fmt.Sprintf("%s:%d", peer.Addr.String(), peer.Idx)
}

// Configuration begin
type Configuration struct {
	peers map[PeerId]struct{}
}

func NewConfiguration(peers []*PeerId) *Configuration {
	c := &Configuration{}
	for _, peer := range peers {
		c.peers[*peer] = struct{}{}
	}
	return c
}

func (c *Configuration) Empty() bool {
	return len(c.peers) == 0
}

func (c *Configuration) Contains(peer *PeerId) bool {
	_, ok := c.peers[*peer]
	return ok
}

func (c *Configuration) Size() int {
	return len(c.peers)
}

func (c *Configuration) ListPeers() []PeerId {
	peers := make([]PeerId, 0, c.Size())
	for peer := range c.peers {
		peers = append(peers, peer)
	}
	return peers
}

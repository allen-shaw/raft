package entity

import (
	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/util"
	"strconv"
)

type PeerID struct {
	Idx      int
	Addr     util.EndPoint
	Priority int
}

func (p *PeerID) Parse(s string) bool {
	if s == "" {
		return false
	}
	tmps, err := util.ParsePeerID(s)
	if err != nil {
		log.Error("Parse peer from string %s failed, err %v", s, err)
		return false
	}
	if len(tmps) < 2 || len(tmps) > 4 {
		return false
	}
	port, err := strconv.Atoi(tmps[1])
	if err != nil {
		log.Error("Parse peer from string %s failed, err %v", s, err)
		return false
	}
	p.Addr = util.EndPoint{IP: tmps[0], Port: port}
	switch len(tmps) {
	case 3:
		idx, err := strconv.Atoi(tmps[2])
		if err != nil {
			log.Error("Parse peer from string %s failed, err %v", s, err)
			return false
		}
		p.Idx = idx
	case 4:
		if tmps[2] == "" {
			p.Idx = 0
		} else {
			idx, err := strconv.Atoi(tmps[2])
			if err != nil {
				log.Error("Parse peer from string %s failed, err %v", s, err)
				return false
			}
			p.Idx = idx
		}
		priority, err := strconv.Atoi(tmps[3])
		if err != nil {
			log.Error("Parse peer from string %s failed, err %v", s, err)
			return false
		}
		p.Priority = priority
	}
	return true
}

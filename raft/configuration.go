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

func (pid *PeerId) Reset() {
	pid.Addr.IP = utils.IP_ANY
	pid.Addr.Port = 0
	pid.Idx = 0
}

func (pid *PeerId) IsEmpty() bool {
	return pid.Addr.IP == utils.IP_ANY && pid.Addr.Port == 0 && pid.Idx == 0
}

func (pid *PeerId) Parse(str string) error {
	pid.Reset()
	addrs := strings.Split(str, ":")
	if len(addrs) < 2 {
		pid.Reset()
		log.Error("Parse %v fail, len %d less than 2", str, len(addrs))
		return fmt.Errorf("invalid addr")
	}
	var err error
	pid.Addr.IP = addrs[0]
	pid.Addr.Port, err = strconv.Atoi(addrs[1])
	if err != nil {
		log.Error("Parse %v fail, %v", str, err)
		return err
	}
	if len(addrs) >= 3 {
		pid.Idx, err = strconv.Atoi(addrs[2])
		if err != nil {
			log.Error("Parse %v fail, %v", str, err)
			return err
		}
	}

	return nil
}

func (pid *PeerId) String() string {
	return fmt.Sprintf("%s:%d", pid.Addr.String(), pid.Idx)
}

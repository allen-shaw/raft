package raft

import "fmt"

type Server struct {
	ID      string
	IP      string
	Port    int
	RaftDir string
}

func (s *Server) Address() string {
	return fmt.Sprintf("%v:%v", s.IP, s.Port)
}

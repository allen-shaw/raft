package raft

import "fmt"

type Server struct {
	ID       string
	IP       string
	Port     int
	HttpPort int
	RaftDir  string
}

func (s *Server) Address() string {
	return fmt.Sprintf("%v:%v", s.IP, s.Port)
}

func (s *Server) HttpAddr() string {
	return fmt.Sprintf("%v:%v", s.IP, s.HttpPort)
}

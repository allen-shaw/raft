package raft

type Server struct {
	ID      string
	Address string
	Port    int
	RaftDir string
}

package util

import (
	"fmt"
)

type Endpoint struct {
	ip   string
	port int
	str  string
}

func NewEndpoint(ip string, port int) *Endpoint {
	return &Endpoint{ip: ip, port: port}
}

func (e *Endpoint) Ip() string {
	return e.ip
}

func (e *Endpoint) Port() int {
	return e.port
}

func (e *Endpoint) String() string {
	if e.str == "" {
		e.str = fmt.Sprintf("%s:%d", e.ip, +e.port)
	}
	return e.str
}

func (e *Endpoint) Copy() *Endpoint {
	return NewEndpoint(e.ip, e.port)
}

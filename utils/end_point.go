package utils

import "fmt"

const (
	IP_ANY = "0.0.0.0"
)

type EndPoint struct {
	IP   string
	Port int
}

func (p *EndPoint) String() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

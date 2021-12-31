package util

import (
	"fmt"
	"strings"
)

const (
	IPv6StartMark = "["
	IPv6EndMark   = "]"
)

// ParsePeerID
// PeerId.parse("a:b")          = new PeerId("a", "b", 0 , -1)
// PeerId.parse("a:b:c")        = new PeerId("a", "b", "c", -1)
// PeerId.parse("a:b::d")       = new PeerId("a", "b", 0, "d")
// PeerId.parse("a:b:c:d")      = new PeerId("a", "b", "c", "d")
func ParsePeerID(s string) ([]string, error) {
	if strings.HasPrefix(s, IPv6StartMark) && strings.Contains(s, IPv6EndMark) {
		//var ipv6Addr string
		//if strings.HasSuffix(s, IPv6EndMark) {
		//	ipv6Addr = s
		//} else {
		//	ipv6Addr = s[0 : strings.Index(s, IPv6EndMark)+1]
		//}
		// TODO: to support IPv6
		return nil, fmt.Errorf("the IPv6 \"%v\"address is incorrect", s)
	} else {
		return strings.Split(s, ":"), nil
	}
}

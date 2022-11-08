package utils

import "strings"

func ParseNodeID(nodeID string) (serverID, groupID string) {
	ss := strings.Split(nodeID, "/")
	return ss[0], ss[1]
}

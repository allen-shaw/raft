package entity

type NodeId struct {
	groupId string
	peerId  *PeerId
	str     string
}

func NewNodeId(groupId string, peerId *PeerId) *NodeId {
	nodeId := &NodeId{}
	nodeId.groupId = groupId
	nodeId.peerId = peerId
	return nodeId
}

func (n *NodeId) GroupId() string {
	return n.groupId
}

func (n *NodeId) PeerId() *PeerId {
	return n.peerId
}

func (n *NodeId) String() string {
	if n.str == "" {
		n.str = "<" + n.groupId + "/" + n.peerId.String() + ">"
	}
	return n.str
}

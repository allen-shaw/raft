package raft

type PosHint struct {
	Pos0 int
	Pos1 int
}

func newPosHint() *PosHint {
	return &PosHint{Pos0: -1, Pos1: -1}
}

type unfoundedPeerId struct {
	peerId PeerId
	found  bool
}

func NewUnfoundedPeerId(peerId *PeerId) *unfoundedPeerId {
	return &unfoundedPeerId{
		peerId: *peerId,
		found:  false,
	}
}

// Equal 对比peerID 结构体的值而不是指针
func (u *unfoundedPeerId) Equal(rhs *unfoundedPeerId) bool {
	return u.peerId == rhs.peerId
}

type Ballot struct {
	peers     []unfoundedPeerId
	oldPeers  []unfoundedPeerId
	quorum    int
	oldQuorum int
}

func NewBallot() *Ballot {
	return &Ballot{quorum: 0, oldQuorum: 0}
}

func (b *Ballot) Swap(rhs *Ballot) {
	b.peers, rhs.peers = rhs.peers, b.peers
	b.quorum, rhs.quorum = rhs.quorum, b.quorum
	b.oldPeers, rhs.oldPeers = rhs.oldPeers, b.oldPeers
	b.oldQuorum, rhs.oldQuorum = rhs.oldQuorum, b.oldQuorum
}

func (b *Ballot) Init(conf, oldConf *Configuration) {
	b.peers = make([]unfoundedPeerId, 0, conf.Size())
	b.oldPeers = make([]unfoundedPeerId, 0)
	b.quorum = 0
	b.oldQuorum = 0

	for peer := range conf.peers {
		b.peers = append(b.peers, *NewUnfoundedPeerId(&peer))
	}
	b.quorum = len(b.peers)/2 + 1
	if oldConf == nil {
		return
	}
	b.oldPeers = make([]unfoundedPeerId, 0, oldConf.Size())
	for oldPeer := range oldConf.peers {
		b.oldPeers = append(b.oldPeers, *NewUnfoundedPeerId(&oldPeer))
	}
	b.oldQuorum = len(b.oldPeers)/2 + 1
}

func (b *Ballot) GrantWithPosHint(peer *PeerId, hint PosHint) PosHint {
	i := b.findPeer(peer, b.peers, hint.Pos0)
	if i != len(b.peers) {
		if !b.peers[i].found {
			b.peers[i].found = true
			b.quorum--
		}
		hint.Pos0 = i
	} else {
		hint.Pos0 = -1
	}

	if len(b.oldPeers) == 0 {
		hint.Pos1 = -1
		return hint
	}

	i = b.findPeer(peer, b.oldPeers, hint.Pos1)
	if i == len(b.oldPeers) {
		if !b.oldPeers[i].found {
			b.oldPeers[i].found = true
			b.oldQuorum--
		}
		hint.Pos1 = i
	} else {
		hint.Pos1 = -1
	}
	return hint
}

func (b *Ballot) Grant(peer *PeerId) PosHint {
	return b.GrantWithPosHint(peer, *newPosHint())
}

func (b *Ballot) Granted() bool {
	return b.quorum <= 0 && b.oldQuorum <= 0
}

func (b *Ballot) findPeer(peer *PeerId, peers []unfoundedPeerId, posHint int) int {
	if posHint < 0 || posHint >= len(peers) || peers[posHint].peerId != *peer {
		for i, p := range peers {
			if p.peerId == *peer {
				return i
			}
		}
		return len(peers) - 1 // peers.end()
	}
	return posHint
}

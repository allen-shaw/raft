package raft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareLogID(t *testing.T) {
	// Term   Index
	{
		// 1,1 == 1,1
		logID1 := &LogID{Term: 1, Index: 1}
		logID2 := &LogID{Term: 1, Index: 1}
		r := CompareLogID(logID1, logID2)
		assert.Equal(t, r, 0)
	}
	{
		// 1,2 > 1,1
		logID1 := &LogID{Term: 1, Index: 2}
		logID2 := &LogID{Term: 1, Index: 1}
		r := CompareLogID(logID1, logID2)
		assert.Equal(t, r, 1)
	}
	{
		// 2,1 > 1,2
		logID1 := &LogID{Term: 2, Index: 1}
		logID2 := &LogID{Term: 1, Index: 2}
		r := CompareLogID(logID1, logID2)
		assert.Equal(t, r, 1)
	}
}

func TestParseAndSerializeConfigurationMeta(t *testing.T) {
	p1 := NewPeerId("127.0.0.1:8848:1")
	p2 := NewPeerId("127.0.0.2:8848:2")
	p3 := NewPeerId("127.0.0.3:8848:3")
	p4 := NewPeerId("127.0.0.4:8848:4")

	entry := &LogEntry{
		OldPeers: []*PeerId{p1, p2, p3},
		Peers:    []*PeerId{p2, p3, p4},
	}
	data, status := SerializeConfigurationMeta(entry)
	assert.True(t, status.OK())

	entry2 := &LogEntry{}
	status = ParseConfigurationMeta(data, entry2)
	assert.True(t, status.OK())

	assert.Equal(t, entry, entry2)
}

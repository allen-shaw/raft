package raft

import "bytes"

type ConfigurationEntry struct {
	id      LogID
	conf    Configuration
	oldConf Configuration
}

func NewConfigurationEntry(entry *LogEntry) *ConfigurationEntry {
	c := &ConfigurationEntry{
		id:   entry.ID,
		conf: *NewConfiguration(entry.Peers),
	}
	if len(entry.OldPeers) != 0 {
		c.oldConf = *NewConfiguration(entry.OldPeers)
	}
	return c
}

func (c *ConfigurationEntry) Stable() bool {
	return c.oldConf.Empty()
}

func (c *ConfigurationEntry) Empty() bool {
	return c.conf.Empty()
}

func (c *ConfigurationEntry) Contains(peer *PeerId) bool {
	return c.conf.Contains(peer) || c.oldConf.Contains(peer)
}

func (c *ConfigurationEntry) String() string {
	var buf bytes.Buffer
	peers := c.conf.ListPeers()
	for i, p := range peers {
		buf.WriteString(p.String())
		if i < len(peers)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

type ConfigurationManager struct {
}

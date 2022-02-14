package raft

import (
	"bytes"
	"errors"

	"github.com/AllenShaw19/raft/log"
)

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

// Manager the history of configuration changing
type ConfigurationManager struct {
	configurations []*ConfigurationEntry // deque<ConfigurationEntry>
	snapshot       *ConfigurationEntry
}

func (m *ConfigurationManager) Add(entry *ConfigurationEntry) error {
	if len(m.configurations) != 0 {
		lastConfIndex := m.configurations[len(m.configurations)-1].id.Index
		if lastConfIndex >= entry.id.Index {
			log.Error("last_config_index=%d larger than entry_index=%d", lastConfIndex, entry.id.Index)
			return errors.New("")
		}
	}
	m.configurations = append(m.configurations, entry)
	return nil
}

//  [1, first_index_kept) are being discarded
func (m *ConfigurationManager) TruncatePrefix(firstIndexKept int64) {
	var idx int
	for i, c := range m.configurations { // TODO: 可以使用二分查找？
		if c.id.Index >= firstIndexKept {
			idx = i
			break
		}
	}
	m.configurations = m.configurations[idx:]
}

// (last_index_kept, infinity) are being discarded
func (m *ConfigurationManager) TruncateSuffix(lastIndexKept int64) {
	var idx int
	for i, c := range m.configurations {
		if c.id.Index > lastIndexKept { // TODO: 可以使用二分查找？
			idx = i
			break
		}
	}
	m.configurations = m.configurations[:idx]
}

func (m *ConfigurationManager) SetSnapshot(entry *ConfigurationEntry) error {
	if CompareLogID(&entry.id, &m.snapshot.id) < 0 {
		log.Error("invalid input snapshot_id=%d, less than manager_snapshot_id=%d", entry.id, m.snapshot.id)
		return errors.New("invalid input snapshot id")
	}
	m.snapshot = entry
	return nil
}

func (m *ConfigurationManager) Get(lastIncludedIndex int64) (*ConfigurationEntry, error) {
	if len(m.configurations) == 0 {
		if lastIncludedIndex < m.snapshot.id.Index {
			log.Error("invalid last_include_index=%d, less than manager_snapshot_id_index=%d", lastIncludedIndex, m.snapshot.id.Index)
			return nil, errors.New("invalid input last_include_index")
		}
		return m.snapshot, nil
	}

	var idx int
	for i, c := range m.configurations {
		if c.id.Index > lastIncludedIndex {
			idx = i
			break
		}
	}
	if idx == 0 {
		return m.snapshot, nil
	}

	return m.configurations[idx], nil
}

func (m *ConfigurationManager) LastConfiguration() *ConfigurationEntry {
	if len(m.configurations) > 0 {
		return m.configurations[len(m.configurations)-1]
	}
	return m.snapshot
}

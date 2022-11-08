package main

import (
	"encoding/json"
	"github.com/AllenShaw19/raft/raft"
	"io"
	"sync"
)

type Table struct {
	ID   string
	lock sync.Mutex
	m    map[string]string
}

func NewTable(id string) *Table {
	t := &Table{}
	t.ID = id
	t.m = make(map[string]string)
	return t
}

func (t *Table) Get(key string) string {
	return t.m[key]
}

func (t *Table) Apply(log *raft.Log) interface{} {
	req := &SetRequest{}
	err := json.Unmarshal(log.Data, req)
	if err != nil {
		return err
	}
	t.m[req.Key] = req.Value
	return nil
}

func (t *Table) Snapshot() (raft.FSMSnapshot, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// Clone the map.
	o := make(map[string]string)
	for k, v := range t.m {
		o[k] = v
	}
	return nil, nil
}

func (t *Table) Restore(snapshot io.ReadCloser) error {
	return nil
}

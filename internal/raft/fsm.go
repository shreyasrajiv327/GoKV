package raft

import (
	"encoding/json"
	"gokv/internal/cache"
	"io"

	hashiraft "github.com/hashicorp/raft"
)

type Command struct{
	Op string `json:"op"`
	Key string `json:"key"`
	Value string `json:"value,omitempty"`
}


type FSM struct{
	store *cache.MemoryStore
}

func NewFSM(store *cache.MemoryStore) *FSM{
	return &FSM{
		store: store,
	}
}

func (f *FSM) Apply(log *hashiraft.Log) interface{} {

	var cmd Command

	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		return err
	}

	switch cmd.Op {

	case "PUT":
		f.store.Put(cmd.Key, cmd.Value)

	case "DELETE":
		f.store.Delete(cmd.Key)
	}

	return nil
}

func (f *FSM) Snapshot() (hashiraft.FSMSnapshot, error) {
	return &Snapshot{}, nil
}

func (f *FSM) Restore(rc io.ReadCloser) error {
	return nil
}

type Snapshot struct{}

func (s *Snapshot) Persist(sink hashiraft.SnapshotSink) error {
	return sink.Close()
}

func (s *Snapshot) Release() {}
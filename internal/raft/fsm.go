package raft

import (
	"encoding/json"
	"io"
     "sync/atomic"
	"gokv/internal/cache"

	hashiraft "github.com/hashicorp/raft"
)

type Command struct {
	Op string `json:"op"`

	Key string `json:"key"`

	Value string `json:"value,omitempty"`
}

type FSM struct {
	store *cache.MemoryStore
    appliedIndex uint64
}

func NewFSM(
	store *cache.MemoryStore,
) *FSM {

	return &FSM{
		store: store,
	}
}

func (f *FSM) Apply(
	log *hashiraft.Log,
) interface{} {

	var cmd Command

	err := json.Unmarshal(
		log.Data,
		&cmd,
	)

	if err != nil {
		return err
	}

	switch cmd.Op {

	case "PUT":

		f.store.Put(
			cmd.Key,
			cmd.Value,
		)

	case "DELETE":

		f.store.Delete(
			cmd.Key,
		)
	}

	atomic.StoreUint64(
		&f.appliedIndex,
		log.Index,
	)

	return nil
}

func (f *FSM) Snapshot() (
	hashiraft.FSMSnapshot,
	error,
) {

	data :=
		f.store.GetAll()

	return &Snapshot{
		data: data,
	}, nil
}

func (f *FSM) Restore(
	rc io.ReadCloser,
) error {

	defer rc.Close()

	var data map[string]string

	err := json.NewDecoder(
		rc,
	).Decode(&data)

	if err != nil {
		return err
	}

	f.store.Clear()

	for k, v := range data {

		f.store.Put(
			k,
			v,
		)
	}

	return nil
}

type Snapshot struct {
	data map[string]string
}

func (s *Snapshot) Persist(
	sink hashiraft.SnapshotSink,
) error {

	defer sink.Close()

	return json.NewEncoder(
		sink,
	).Encode(
		s.data,
	)
}

func (f *FSM) AppliedIndex() uint64 {

	return atomic.LoadUint64(
		&f.appliedIndex,
	)
}

func (s *Snapshot) Release() {}
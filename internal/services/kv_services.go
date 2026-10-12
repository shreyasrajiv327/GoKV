package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gokv/internal/raft"
	"gokv/internal/repository"

	hashiraft "github.com/hashicorp/raft"
)

type KVService struct {
	repo *repository.KVRepository

	raft *hashiraft.Raft
	fsm  *raft.FSM

	httpNodes map[string]string
}

func NewKVService(
	repo *repository.KVRepository,
	raftNode *hashiraft.Raft,
	fsm *raft.FSM,
	httpNodes map[string]string,
) *KVService {

	return &KVService{
		repo:      repo,
		raft:      raftNode,
		fsm:       fsm,
		httpNodes: httpNodes,
	}
}

func (s *KVService) Put(
	key string,
	value string,
) error {

	if s.raft.State() != hashiraft.Leader {

		return errors.New(
			"writes allowed only on leader",
		)
	}

	cmd := raft.Command{
		Op:    "PUT",
		Key:   key,
		Value: value,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	future := s.raft.Apply(
		data,
		5*time.Second,
	)

	return future.Error()
}

func (s *KVService) Delete(
	key string,
) error {

	if s.raft.State() != hashiraft.Leader {

		return errors.New(
			"writes allowed only on leader",
		)
	}

	cmd := raft.Command{
		Op:  "DELETE",
		Key: key,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	future := s.raft.Apply(
		data,
		5*time.Second,
	)

	return future.Error()
}

func (s *KVService) Get(
	key string,
	consistency string,
) (string, bool, error) {

	if consistency == "strong" {

		if s.raft.State() != hashiraft.Leader {

			return "",
				false,
				errors.New(
					"strong reads allowed only on leader",
				)
		}
	}

	value, ok := s.repo.Get(key)

	return value,
		ok,
		nil
}

func (s *KVService) CanServeSmartRead() bool {

	if s.raft.State() ==
		hashiraft.Leader {

		return true
	}

	stats :=
		s.raft.Stats()

	lastApplied :=
		stats["last_log_index"]

	commit :=
		stats["commit_index"]

	return lastApplied ==
		commit
}

func (s *KVService) ForwardStrongRead(
	key string,
) (*http.Response, error) {

	leader :=
		string(
			s.raft.Leader(),
		)

	if leader == "" {

		return nil,
			errors.New(
				"leader unavailable",
			)
	}

	port, ok :=
		s.httpNodes[leader]

	if !ok {

		return nil,
			errors.New(
				"leader mapping missing",
			)
	}

	url :=
		"http://localhost:" +
			port +
			"/kv/" +
			key +
			"?consistency=strong"

	return http.Get(url)
}
package services

import (
	"encoding/json"
	"errors"
	"time"

	"gokv/internal/raft"
	"gokv/internal/repository"
	"gokv/internal/wal"

	hashiraft "github.com/hashicorp/raft"
)

type KVService struct {
	repo *repository.KVRepository
	wal  *wal.WAL

	raft *hashiraft.Raft
}

func NewKVService(
	repo *repository.KVRepository,
	wal *wal.WAL,
	raftNode *hashiraft.Raft,
) *KVService {

	return &KVService{
		repo: repo,
		wal: wal,
		raft: raftNode,
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
		Op: "PUT",
		Key: key,
		Value: value,
	}

	data, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	future :=
		s.raft.Apply(
			data,
			5*time.Second,
		)

	if err := future.Error(); err != nil {
		return err
	}

	return s.wal.LogPut(
		key,
		value,
	)
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
		Op: "DELETE",
		Key: key,
	}

	data, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	future :=
		s.raft.Apply(
			data,
			5*time.Second,
		)

	if err := future.Error(); err != nil {
		return err
	}

	return s.wal.LogDelete(
		key,
	)
}

func (s *KVService) Get(
	key string,
) (string, bool) {

	return s.repo.Get(
		key,
	)
}
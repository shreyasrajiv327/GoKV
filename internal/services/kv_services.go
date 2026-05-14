package services

import (
	"gokv/internal/repository"
	"gokv/internal/wal"
	"gokv/internal/cluster"
)

type KVService struct {
	repo       *repository.KVRepository
	wal        *wal.WAL
	replicator *cluster.Replicator
	isLeader   bool

}

func NewKVService(repo *repository.KVRepository,
	wal *wal.WAL,
	replicator *cluster.Replicator,
	isLeader bool,
	) *KVService {
	return &KVService{
		repo: repo,
		wal: wal,
		replicator: replicator,
		isLeader: isLeader,
	}
}

func (s *KVService) Put(key, value string) error {
	if err := s.wal.LogPut(key, value);  err != nil{
		return err
	}
	s.repo.Put(key, value)

	if s.isLeader && s.replicator !=nil{
		s.replicator.ReplicatePut(key, value)
	}
	return nil
}

func (s *KVService) Get(key string) (string, bool) {
	return s.repo.Get(key)
}

func (s *KVService) Delete(key string) error {

	if err := s.wal.LogDelete(key); err != nil{
		return err
	}
	s.repo.Delete(key)

	if s.isLeader && s.replicator !=nil{
		s.replicator.ReplicateDelete(key)
	}

	return nil
}
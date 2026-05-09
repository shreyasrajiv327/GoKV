package services

import (
	"gokv/internal/repository"
	"gokv/internal/wal"
)

type KVService struct {
	repo *repository.KVRepository
	wal *wal.WAL

}

func NewKVService(repo *repository.KVRepository,
	wal *wal.WAL) *KVService {
	return &KVService{
		repo: repo,
		wal: wal,
	}
}

func (s *KVService) Put(key, value string) error {
	if err := s.wal.LogPut(key, value);  err != nil{
		return err
	}
	s.repo.Put(key, value)
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

	return nil
}
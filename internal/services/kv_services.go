package services

import "gokv/internal/repository"

type KVService struct {
	repo *repository.KVRepository
}

func NewKVService(repo *repository.KVRepository) *KVService {
	return &KVService{
		repo: repo,
	}
}

func (s *KVService) Put(key, value string) {
	s.repo.Put(key, value)
}

func (s *KVService) Get(key string) (string, bool) {
	return s.repo.Get(key)
}

func (s *KVService) Delete(key string) {
	s.repo.Delete(key)
}
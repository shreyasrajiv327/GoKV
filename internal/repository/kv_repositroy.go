package repository

import "gokv/internal/cache"


type KVRepository struct{
	store *cache.MemoryStore
}

func NewKVRepositroy(store *cache.MemoryStore) *KVRepository{
	return  &KVRepository{
		store: store,
	}
}

func (r *KVRepository) Put(key, value string){
	r.store.Put(key, value)
}

func (r *KVRepository) Get(key string) (string, bool) {
	return r.store.Get(key)
}

func (r *KVRepository) Delete(key string) {
	r.store.Delete(key)
}
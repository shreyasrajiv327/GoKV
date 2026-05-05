package cache

type MemoryStore struct{
	data map[string]string
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStore) Put(key, value string){
	m.data[key] = value
}

func (m *MemoryStore) Get(key string) (string , bool){
	val, ok := m.data[key]
	return val, ok
}

func (m *MemoryStore) Delete(key string){
	delete(m.data, key)
}
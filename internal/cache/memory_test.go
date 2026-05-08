package cache


import (
	"strconv"
	"sync"
	"testing"
)

func TestMemoryStoreConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)

			store.Put(key, value)

			got, ok := store.Get(key)
			if !ok {
				t.Errorf("expected key %s to exist", key)
			}

			if got != value {
				t.Errorf("expected %s, got %s", value, got)
			}

			store.Delete(key)
		}(i)
	}

	wg.Wait()
}
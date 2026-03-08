package cache

import (
	"strconv"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

type Cache struct {
	cache   map[string]Data
	ttl     int
	ttlOnce sync.Once

	mu sync.RWMutex
}

func NewCache(ttl int) *Cache {
	c := &Cache{
		cache: make(map[string]Data),
		ttl:   ttl,
	}

	c.RunEvictionWorker()

	return c
}

func (c *Cache) SetValue(value string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	hashVal := strconv.Itoa(getHash(value))

	expireAt := time.Now().Add(time.Duration(c.ttl) * time.Second)

	data := Data{
		Value:    value,
		ExpireAt: expireAt,
	}

	c.cache[hashVal] = data

	return hashVal
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	data, ok := c.cache[key]
	c.mu.RUnlock()

	if !ok {
		return "", false
	}

	if time.Now().After(data.ExpireAt) {
		c.Delete(key)
		return "", false
	}

	return data.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func getHash(value string) int {
	return int(murmur3.Sum32([]byte(value)))
}

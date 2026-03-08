package cache

import (
	"strconv"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

type Cache struct {
	cache map[string]Data
	ttl   int
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]Data),
	}
}

func (c *Cache) SetValue(value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	hashVal := getHash(value)

	expireAt := time.Now().Add(time.Duration(c.ttl) * time.Second)

	data := Data{
		Value:    value,
		ExpireAt: expireAt,
	}

	c.cache[strconv.Itoa(hashVal)] = data
}

func (c *Cache) GetValue(key string) (string, bool) {
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

func getHash(val string) int {
	return int(murmur3.Sum32([]byte(val)))
}

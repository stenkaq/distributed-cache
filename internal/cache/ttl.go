package cache

import "time"

func (c *Cache) RunEvictionWorker() {
	go func() {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			c.cleanup()
		}
	}()
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, item := range c.cache {
		if now.After(item.ExpireAt) {
			c.Delete(key)
		}
	}

}

package pokecashe

import (
	"sync"
	"time"
)

type Cashe struct {
	mu    sync.RWMutex
	cache map[string]casheEntry
}
type casheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCashe(interval time.Duration) *Cashe {
	c := &Cashe{
		cache: make(map[string]casheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cashe) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = casheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cashe) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cashe) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now(), interval)
	}
}
func (c *Cashe) reap(now time.Time, interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.cache {
		if entry.createdAt.Before(now.Add(-interval)) {
			delete(c.cache, key)
		}
	}

}

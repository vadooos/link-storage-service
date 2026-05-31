package cache

import (
	"link-storage-service/model"
	"sync"
	"time"
)

type item struct {
	link    model.Link
	expires time.Time
}

type MemoryCache struct {
	items map[string]item
	ttl   time.Duration
	mu    sync.RWMutex
}

func New(ttl time.Duration) *MemoryCache {
	return &MemoryCache{items: make(map[string]item), ttl: ttl}
}

func (c *MemoryCache) Get(shortCode string) (model.Link, bool) {
	c.mu.RLock()
	it, ok := c.items[shortCode]
	c.mu.RUnlock()
	if !ok {
		return model.Link{}, false
	}
	if time.Now().After(it.expires) {
		c.mu.Lock()
		cur, ok := c.items[shortCode]
		if ok && time.Now().After(cur.expires) {
			delete(c.items, shortCode)
		}
		c.mu.Unlock()
		return model.Link{}, false
	}
	return it.link, true
}

func (c *MemoryCache) Set(link model.Link) {
	c.mu.Lock()
	c.items[link.ShortCode] = item{link: link, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *MemoryCache) Delete(shortCode string) {
	c.mu.Lock()
	delete(c.items, shortCode)
	c.mu.Unlock()
}

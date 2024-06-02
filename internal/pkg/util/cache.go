package util

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data     interface{}
	CachedAt time.Time
	TTL      time.Duration
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Data:     data,
		CachedAt: time.Now(),
		TTL:      ttl,
	}
}

func (c *Cache) Get(key string) (data interface{}, found bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Since(item.CachedAt) > item.TTL {
		return nil, false
	}

	return item.Data, true
}

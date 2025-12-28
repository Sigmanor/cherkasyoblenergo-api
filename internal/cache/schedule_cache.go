package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

type ScheduleCache struct {
	entries map[string]CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

func NewScheduleCache(ttlSeconds int) *ScheduleCache {
	return &ScheduleCache{
		entries: make(map[string]CacheEntry),
		ttl:     time.Duration(ttlSeconds) * time.Second,
	}
}

func (c *ScheduleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

func (c *ScheduleCache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *ScheduleCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

func (c *ScheduleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]CacheEntry)
}

func (c *ScheduleCache) TTL() time.Duration {
	return c.ttl
}

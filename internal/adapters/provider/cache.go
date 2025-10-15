package provider

import (
	"sync"
	"time"
)

const (
	CleanupInterval = 30 * time.Second
)

type CacheItem[T any] struct {
	Value     T
	ExpiresAt time.Time
}

func (c *CacheItem[T]) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

type Cache[T any] struct {
	items map[string]*CacheItem[T]
	mutex sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	cache := &Cache[T]{
		items: make(map[string]*CacheItem[T]),
	}

	go cache.cleanup()

	return cache
}

func newCacheWithInterval[T any](cleanupInterval time.Duration) *Cache[T] {
	cache := &Cache[T]{
		items: make(map[string]*CacheItem[T]),
	}

	go cache.cleanupWithInterval(cleanupInterval)

	return cache
}

func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &CacheItem[T]{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

//nolint:ireturn
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var zero T
	item, exists := c.items[key]
	if !exists {
		return zero, false
	}

	if item.IsExpired() {
		return zero, false
	}

	return item.Value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *Cache[T]) cleanup() {
	c.cleanupWithInterval(CleanupInterval)
}

func (c *Cache[T]) cleanupWithInterval(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpiredItems()
	}
}

func (c *Cache[T]) cleanupExpiredItems() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, item := range c.items {
		if item.IsExpired() {
			delete(c.items, key)
		}
	}
}

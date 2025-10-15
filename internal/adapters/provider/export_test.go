package provider

import "time"

func (c *Cache[T]) CleanupExpiredItems() {
	c.cleanupExpiredItems()
}

func NewCacheWithInterval[T any](cleanupInterval time.Duration) *Cache[T] {
	return newCacheWithInterval[T](cleanupInterval)
}

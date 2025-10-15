package provider_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[string]()

	cache.Set("key1", "value1", 10*time.Second)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)
}

func TestCache_GetNonExistent(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[string]()

	value, found := cache.Get("nonexistent")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Expiration(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[int]()

	cache.Set("key1", 42, 1*time.Millisecond)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, 42, value)

	time.Sleep(5 * time.Millisecond)

	value, found = cache.Get("key1")
	assert.False(t, found)
	assert.Equal(t, 0, value)
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[string]()

	cache.Set("key1", "value1", 10*time.Second)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	cache.Delete("key1")

	value, found = cache.Get("key1")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_OverwriteValue(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[string]()

	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key1", "value2", 10*time.Second)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value2", value)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[int]()

	done := make(chan bool, 10)

	for i := range 5 {
		go func(val int) {
			for range 10 {
				cache.Set("concurrent", val, 10*time.Second)
			}
			done <- true
		}(i)
	}

	for range 5 {
		go func() {
			for range 10 {
				cache.Get("concurrent")
			}
			done <- true
		}()
	}

	for range 10 {
		<-done
	}

	cache.Set("final", 999, 10*time.Second)
	value, found := cache.Get("final")
	assert.True(t, found)
	assert.Equal(t, 999, value)
}

func TestCacheItem_IsExpired(t *testing.T) {
	t.Parallel()
	expiredItem := &provider.CacheItem[string]{
		Value:     "test",
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}
	assert.True(t, expiredItem.IsExpired())

	validItem := &provider.CacheItem[string]{
		Value:     "test",
		ExpiresAt: time.Now().Add(1 * time.Second),
	}
	assert.False(t, validItem.IsExpired())
}

func TestCache_DifferentTypes(t *testing.T) {
	t.Parallel()
	stringCache := provider.NewCache[string]()
	intCache := provider.NewCache[int]()
	structCache := provider.NewCache[*provider.CachedRateResult]()

	stringCache.Set("str", "hello", 10*time.Second)
	value1, found1 := stringCache.Get("str")
	assert.True(t, found1)
	assert.Equal(t, "hello", value1)

	intCache.Set("num", 42, 10*time.Second)
	value2, found2 := intCache.Get("num")
	assert.True(t, found2)
	assert.Equal(t, 42, value2)

	rate := &provider.CachedRateResult{Rate: 100.5, Change24h: 2.5}
	structCache.Set("rate", rate, 10*time.Second)
	value3, found3 := structCache.Get("rate")
	assert.True(t, found3)
	assert.InEpsilon(t, 100.5, value3.Rate, 0.001)
	assert.InEpsilon(t, 2.5, value3.Change24h, 0.001)
}

//nolint:paralleltest // This test is timing-sensitive and cannot run in parallel
func TestCache_Cleanup(t *testing.T) {
	caches := make([]*provider.Cache[string], 10)
	for i := range caches {
		caches[i] = provider.NewCache[string]()
		caches[i].Set("shortlived", "value", 1*time.Millisecond)
		caches[i].Set("longlived", "value", 10*time.Second)
	}

	time.Sleep(10 * time.Millisecond)

	for i, cache := range caches {
		_, foundShort := cache.Get("shortlived")
		_, foundLong := cache.Get("longlived")
		assert.False(t, foundShort, "Cache %d short-lived item should be expired", i)
		assert.True(t, foundLong, "Cache %d long-lived item should still exist", i)
	}

	testCache := provider.NewCache[int]()
	testCache.Set("test", 42, 5*time.Second)

	value, found := testCache.Get("test")
	assert.True(t, found)
	assert.Equal(t, 42, value)

	for i := range 100 {
		testCache.Set(fmt.Sprintf("stress_%d", i), i, 1*time.Millisecond)
	}

	time.Sleep(5 * time.Millisecond)

	expiredCount := 0
	for i := range 100 {
		_, found := testCache.Get(fmt.Sprintf("stress_%d", i))
		if !found {
			expiredCount++
		}
	}

	assert.Greater(t, expiredCount, 90, "Most stress test items should be expired")
}

func TestCache_CleanupExpiredItems(t *testing.T) {
	t.Parallel()
	cache := provider.NewCache[string]()

	cache.Set("expired1", "value1", -1*time.Second)
	cache.Set("expired2", "value2", -1*time.Second)
	cache.Set("valid1", "value3", 10*time.Second)
	cache.Set("expired3", "value4", -1*time.Second)
	cache.Set("valid2", "value5", 10*time.Second)

	cache.CleanupExpiredItems()

	_, found1 := cache.Get("expired1")
	_, found2 := cache.Get("expired2")
	_, found3 := cache.Get("valid1")
	_, found4 := cache.Get("expired3")
	_, found5 := cache.Get("valid2")

	assert.False(t, found1)
	assert.False(t, found2)
	assert.True(t, found3)
	assert.False(t, found4)
	assert.True(t, found5)

	value3, found3Again := cache.Get("valid1")
	value5, found5Again := cache.Get("valid2")
	assert.True(t, found3Again)
	assert.True(t, found5Again)
	assert.Equal(t, "value3", value3)
	assert.Equal(t, "value5", value5)

	intCache := provider.NewCache[int]()
	intCache.Set("expired_int", 42, -1*time.Second)
	intCache.Set("valid_int", 99, 10*time.Second)

	intCache.CleanupExpiredItems()

	_, expiredFound := intCache.Get("expired_int")
	validValue, validFound := intCache.Get("valid_int")

	assert.False(t, expiredFound)
	assert.True(t, validFound)
	assert.Equal(t, 99, validValue)
}

//nolint:paralleltest // This test is timing-sensitive and cannot run in parallel
func TestCache_CleanupGoroutine(t *testing.T) {
	caches := make([]*provider.Cache[string], 50)
	for i := range caches {
		caches[i] = provider.NewCache[string]()
		caches[i].Set("temp", fmt.Sprintf("value%d", i), 1*time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)

	testCache := provider.NewCache[int]()
	testCache.Set("test", 42, 5*time.Second)

	value, found := testCache.Get("test")
	assert.True(t, found)
	assert.Equal(t, 42, value)

	for i, cache := range caches {
		cache.Set("final", fmt.Sprintf("final%d", i), 1*time.Second)
		val, exists := cache.Get("final")
		assert.True(t, exists)
		assert.Equal(t, fmt.Sprintf("final%d", i), val)
	}
}

//nolint:paralleltest // This test is timing-sensitive and cannot run in parallel
func TestCache_CleanupWithShortInterval(t *testing.T) {
	cache := provider.NewCacheWithInterval[string](10 * time.Millisecond)

	cache.Set("will_expire", "value1", 1*time.Millisecond)
	cache.Set("will_remain", "value2", 1*time.Second)

	time.Sleep(5 * time.Millisecond)
	time.Sleep(15 * time.Millisecond)

	_, expiredFound := cache.Get("will_expire")
	remainValue, remainFound := cache.Get("will_remain")

	assert.False(t, expiredFound)
	assert.True(t, remainFound)
	assert.Equal(t, "value2", remainValue)

	cache.Set("new_item", "new_value", 1*time.Second)
	newValue, newFound := cache.Get("new_item")
	assert.True(t, newFound)
	assert.Equal(t, "new_value", newValue)
}

//nolint:paralleltest // This test is timing-sensitive and cannot run in parallel
func TestCache_NewCacheUsesCleanup(t *testing.T) {
	cache := provider.NewCache[string]()

	cache.Set("test1", "value1", 1*time.Second)
	cache.Set("test2", "value2", 1*time.Second)

	value1, found1 := cache.Get("test1")
	value2, found2 := cache.Get("test2")

	assert.True(t, found1)
	assert.True(t, found2)
	assert.Equal(t, "value1", value1)
	assert.Equal(t, "value2", value2)

	cache.Delete("test1")
	_, found1After := cache.Get("test1")
	value2After, found2After := cache.Get("test2")

	assert.False(t, found1After)
	assert.True(t, found2After)
	assert.Equal(t, "value2", value2After)
}

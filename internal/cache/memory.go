package cache

import (
	"context"
	"strings"
	"sync"
	"time"
)

// cacheItem 缓存项，包含值和过期时间
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// isExpired 检查缓存项是否已过期
func (item *cacheItem) isExpired() bool {
	if item.expiration.IsZero() {
		return false // 如果没有设置过期时间，则永不过期
	}
	return time.Now().After(item.expiration)
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
	stop chan struct{}
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache() (Cache, error) {
	mc := &MemoryCache{
		data: make(map[string]*cacheItem),
		mu:   sync.RWMutex{},
		stop: make(chan struct{}),
	}

	// 启动后台清理 goroutine，定期清理过期项
	go mc.cleanup()

	return mc, nil
}

// cleanup 定期清理过期的缓存项
func (mc *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.mu.Lock()
			for key, item := range mc.data {
				if item.isExpired() {
					delete(mc.data, key)
				}
			}
			mc.mu.Unlock()
		case <-mc.stop:
			return
		}
	}
}

// Set 设置缓存值，支持过期时间
func (mc *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item := &cacheItem{
		value: value,
	}

	// 如果设置了过期时间，则计算过期时间点
	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
	}

	mc.data[key] = item
	return nil
}

// Get 获取缓存值，如果过期则返回 nil
func (mc *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return nil, nil
	}

	// 检查是否过期
	if item.isExpired() {
		// 如果过期，异步删除（避免在读锁中加写锁）
		go func() {
			mc.mu.Lock()
			delete(mc.data, key)
			mc.mu.Unlock()
		}()
		return nil, nil
	}

	return item.value, nil
}

// Get 获取缓存值，如果过期则返回 nil
func (mc *MemoryCache) GetAll(ctx context.Context, pattern string) (interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	result := make(map[string]*cacheItem)
	for key, value := range mc.data {
		if strings.HasPrefix(key, pattern) {
			result[key] = value
		}
	}

	return result, nil
}

// Delete 删除指定的缓存键
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.data, key)
	return nil
}

// Exists 检查键是否存在且未过期
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return false, nil
	}

	// 如果过期，返回 false
	if item.isExpired() {
		// 异步删除过期项
		go func() {
			mc.mu.Lock()
			delete(mc.data, key)
			mc.mu.Unlock()
		}()
		return false, nil
	}

	return true, nil
}

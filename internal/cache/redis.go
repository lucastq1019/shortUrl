package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis 缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建新的 Redis 缓存实例
func NewRedisCache(addr string) (Cache, error) {
	// 解析 Redis 地址，支持 "host:port" 格式
	opts := &redis.Options{
		Addr:     addr,
		Password: "", // 如果需要密码，可以从配置中读取
		DB:       0,  // 默认使用 0 号数据库
	}

	client := redis.NewClient(opts)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Set 设置缓存值，支持过期时间
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 将 value 序列化为 JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 如果设置了过期时间，使用 SetEX；否则使用 Set（永不过期）
	if expiration > 0 {
		return rc.client.Set(ctx, key, data, expiration).Err()
	}
	return rc.client.Set(ctx, key, data, 0).Err()
}

// Get 获取缓存值
func (rc *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// 键不存在，返回 nil
			return nil, nil
		}
		return nil, err
	}

	// 尝试将 JSON 反序列化为 interface{}
	var result any
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// 如果反序列化失败，直接返回原始字符串
		return val, nil
	}

	return result, nil
}

// Delete 删除指定的缓存键
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

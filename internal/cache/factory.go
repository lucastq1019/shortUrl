package cache

import (
	"github.com/username/shorturl/internal/config"
)

func NewCache(cfg *config.Config) (Cache, error) {
	if cfg.RedisAddr != "" {
		return NewRedisCache(cfg.RedisAddr)
	}
	return NewMemoryCache()
}

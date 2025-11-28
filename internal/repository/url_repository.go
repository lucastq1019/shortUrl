package repository

import (
	"context"

	"github.com/username/shorturl/internal/model"
)

type URLRepository interface {
	// Get 从多个数据源并发获取，谁先返回就用谁的
	// 优先级：RedisCache > MemoryCache > MySQLDB > SQLiteDB
	Get(ctx context.Context, shortCode string) (*model.ShortURL, error)

	// Save 保存到多个数据源
	// 缓存：优先写入 Redis，如果失败则写入 Memory
	// 数据库：优先写入 MySQL，如果失败则写入 SQLite
	Save(ctx context.Context, url *model.ShortURL) error

	// 以下方法保持向后兼容
	SaveToCache(ctx context.Context, url *model.ShortURL) error
	GetFromCache(ctx context.Context, shortCode string) (*model.ShortURL, error)

	SaveToDB(ctx context.Context, url *model.ShortURL) error
	GetFromDB(ctx context.Context, shortCode string) (*model.ShortURL, error)

	DeleteFromCache(ctx context.Context, shortCode string) error
}

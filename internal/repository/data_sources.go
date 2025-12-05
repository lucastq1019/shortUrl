package repository

import (
	"log"

	"github.com/username/shorturl/internal/cache"
	"github.com/username/shorturl/internal/config"
	"github.com/username/shorturl/internal/db"
	"gorm.io/gorm"
)

// DataSources 管理所有可用的数据源
type DataSources struct {
	// 缓存：Redis 优先，Memory 作为 fallback
	RedisCache  cache.Cache
	MemoryCache cache.Cache

	// 数据库：MySQL 优先，SQLite 作为 fallback
	MySQLDB  db.Database
	SQLiteDB db.Database
}

// NewDataSources 创建数据源管理器
// 尝试初始化所有数据源，即使某些失败也会继续
func NewDataSources(cfg *config.Config) *DataSources {
	ds := &DataSources{}

	// 初始化缓存
	// 尝试创建 Redis（如果配置了）
	if cfg.RedisAddr != "" {
		if redisCache, err := cache.NewRedisCache(cfg.RedisAddr); err == nil {
			log.Println("初始化redis:", cfg.RedisAddr)
			ds.RedisCache = redisCache
		}
	}
	// MemoryCache 总是可用的
	if memoryCache, err := cache.NewMemoryCache(); err == nil {
		log.Println("初始化memoryCache")
		ds.MemoryCache = memoryCache
	}

	// 初始化数据库
	// 尝试创建 MySQL（如果配置了）
	if cfg.MySQLDSN != "" {
		if mysqlDB, err := db.NewMySQLDB(cfg.MySQLDSN); err == nil {
			log.Println("初始化MysqlDB")
			ds.MySQLDB = mysqlDB
		}
	}
	// 尝试创建 SQLite（总是尝试，作为最后的 fallback）
	if sqliteDB, err := db.NewSQLiteDB(cfg.SQLitePath); err == nil {
		log.Println("初始化SqlLite")
		ds.SQLiteDB = sqliteDB
	}

	return ds
}

var GloablDataSources *DataSources

func GetDataSources() (*DataSources, error) {
	if GloablDataSources == nil {
		GloablDataSources = NewDataSources(config.GetConfig())

	}
	return GloablDataSources, nil
}

// migrateTable is migrate db table.
func migrateTable(db *gorm.DB, tables ...interface{}) {
	err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(tables...)
	if err != nil {

		panic(err)
	}
}

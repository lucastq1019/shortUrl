package db

import (
	"github.com/username/shorturl/internal/config"
)

// NewDatabase 根据配置创建数据库连接
// 如果配置了 MySQLDSN，则使用 MySQL；否则使用 SQLite
func NewDatabase(cfg *config.Config) (Database, error) {
	if cfg.MySQLDSN != "" {
		return NewMySQLDB(cfg.MySQLDSN)
	}
	return NewSQLiteDB(cfg.SQLitePath)
}


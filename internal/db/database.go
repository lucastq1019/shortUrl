package db

import (
	"database/sql"
)

// Database 数据库接口，只负责管理连接
// 具体的查询操作由调用方自行决定
type Database interface {
	// GetDB 获取数据库连接
	GetDB() *sql.DB
	// Close 关闭数据库连接
	Close() error
}

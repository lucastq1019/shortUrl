package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB MySQL 数据库实现
type MySQLDB struct {
	db *sql.DB
}

// NewMySQLDB 创建新的 MySQL 数据库连接
func NewMySQLDB(dsn string) (Database, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return &MySQLDB{db: db}, nil
}

// GetDB 获取数据库连接
func (m *MySQLDB) GetDB() *sql.DB {
	return m.db
}

// Close 关闭数据库连接
func (m *MySQLDB) Close() error {
	return m.db.Close()
}

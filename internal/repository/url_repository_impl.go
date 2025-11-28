package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/username/shorturl/internal/model"
)

// urlRepository 实现 URLRepository 接口
type urlRepository struct {
	sources *DataSources
}

// NewURLRepository 创建新的 URL Repository
func NewURLRepository(sources *DataSources) URLRepository {
	return &urlRepository{
		sources: sources,
	}
}

// Get 从多个数据源并发获取，谁先返回就用谁的
// 优先级：RedisCache > MemoryCache > MySQLDB > SQLiteDB
func (r *urlRepository) Get(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	type result struct {
		url *model.ShortURL
		err error
	}

	resultCh := make(chan result, 1) // 只需要第一个结果
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once

	// 辅助函数：发送结果并取消其他 goroutine
	sendResult := func(url *model.ShortURL, err error) {
		once.Do(func() {
			select {
			case resultCh <- result{url: url, err: err}:
				cancel() // 取消其他 goroutine
			case <-ctx.Done():
			}
		})
	}

	// 从 Redis 获取
	if r.sources.RedisCache != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := r.getFromRedis(ctx, shortCode)
			if err == nil && url != nil {
				sendResult(url, nil)
			}
		}()
	}

	// 从 Memory 获取
	if r.sources.MemoryCache != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := r.getFromMemory(ctx, shortCode)
			if err == nil && url != nil {
				sendResult(url, nil)
			}
		}()
	}

	// 从 MySQL 获取
	if r.sources.MySQLDB != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := r.getFromMySQL(ctx, shortCode)
			if err == nil && url != nil {
				sendResult(url, nil)
			}
		}()
	}

	// 从 SQLite 获取
	if r.sources.SQLiteDB != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := r.getFromSQLite(ctx, shortCode)
			if err == nil && url != nil {
				sendResult(url, nil)
			}
		}()
	}

	// 等待第一个结果或所有 goroutine 完成
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// 等待第一个结果
	select {
	case res := <-resultCh:
		if res.err != nil {
			return nil, res.err
		}
		if res.url != nil {
			// 异步回写缓存（如果从数据库获取的）
			go r.asyncWriteToCache(context.Background(), res.url)
			return res.url, nil
		}
		// 如果收到 nil，继续等待其他结果
		// 这种情况理论上不应该发生，因为只有 url != nil 时才发送
		// 但为了安全，我们继续等待
		select {
		case res := <-resultCh:
			if res.url != nil {
				go r.asyncWriteToCache(context.Background(), res.url)
				return res.url, nil
			}
			// 继续等待所有完成
			<-doneCh
			return nil, fmt.Errorf("short code not found: %s", shortCode)
		case <-doneCh:
			return nil, fmt.Errorf("short code not found: %s", shortCode)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	case <-doneCh:
		// 所有 goroutine 都完成了，但没有找到结果
		return nil, fmt.Errorf("short code not found: %s", shortCode)
	case <-ctx.Done():
		// 外部 context 被取消
		wg.Wait()
		return nil, ctx.Err()
	}
}

// Save 保存到多个数据源
// 缓存：优先写入 Redis，如果失败则写入 Memory
// 数据库：优先写入 MySQL，如果失败则写入 SQLite
func (r *urlRepository) Save(ctx context.Context, url *model.ShortURL) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 4)

	// 保存到缓存
	if r.sources.RedisCache != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.saveToRedis(ctx, url); err != nil {
				// Redis 失败，尝试 Memory
				if r.sources.MemoryCache != nil {
					if err := r.saveToMemory(ctx, url); err != nil {
						errCh <- fmt.Errorf("failed to save to cache: %w", err)
					}
				}
			}
		}()
	} else if r.sources.MemoryCache != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.saveToMemory(ctx, url); err != nil {
				errCh <- fmt.Errorf("failed to save to memory cache: %w", err)
			}
		}()
	}

	// 保存到数据库
	if r.sources.MySQLDB != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.saveToMySQL(ctx, url); err != nil {
				// MySQL 失败，尝试 SQLite
				if r.sources.SQLiteDB != nil {
					if err := r.saveToSQLite(ctx, url); err != nil {
						errCh <- fmt.Errorf("failed to save to database: %w", err)
					}
				} else {
					errCh <- fmt.Errorf("failed to save to MySQL: %w", err)
				}
			}
		}()
	} else if r.sources.SQLiteDB != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.saveToSQLite(ctx, url); err != nil {
				errCh <- fmt.Errorf("failed to save to SQLite: %w", err)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	// 收集错误
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("some saves failed: %v", errs)
	}

	return nil
}

// 从各个数据源获取的辅助方法

func (r *urlRepository) getFromRedis(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	key := "shorturl:" + shortCode
	val, err := r.sources.RedisCache.Get(ctx, key)
	if err != nil || val == nil {
		return nil, err
	}

	var url model.ShortURL
	data, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cache value type")
	}
	if err := json.Unmarshal([]byte(data), &url); err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) getFromMemory(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	key := "shorturl:" + shortCode
	val, err := r.sources.MemoryCache.Get(ctx, key)
	if err != nil || val == nil {
		return nil, err
	}

	var url model.ShortURL
	data, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cache value type")
	}
	if err := json.Unmarshal([]byte(data), &url); err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) getFromMySQL(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	db := r.sources.MySQLDB.GetDB()
	query := `SELECT id, short_code, long_url, created_at, expires_at FROM short_urls WHERE short_code = ?`

	var url model.ShortURL
	var expiresAt sql.NullTime
	err := db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID, &url.ShortCode, &url.LongURL, &url.CreatedAt, &expiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if expiresAt.Valid {
		url.ExpiresAt = &expiresAt.Time
	}

	return &url, nil
}

func (r *urlRepository) getFromSQLite(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	db := r.sources.SQLiteDB.GetDB()
	query := `SELECT id, short_code, long_url, created_at, expires_at FROM short_urls WHERE short_code = ?`

	var url model.ShortURL
	var expiresAt sql.NullTime
	err := db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID, &url.ShortCode, &url.LongURL, &url.CreatedAt, &expiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if expiresAt.Valid {
		url.ExpiresAt = &expiresAt.Time
	}

	return &url, nil
}

// 保存到各个数据源的辅助方法

func (r *urlRepository) saveToRedis(ctx context.Context, url *model.ShortURL) error {
	key := "shorturl:" + url.ShortCode
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	var expiration time.Duration
	if url.ExpiresAt != nil {
		expiration = time.Until(*url.ExpiresAt)
		if expiration < 0 {
			return fmt.Errorf("url already expired")
		}
	}

	return r.sources.RedisCache.Set(ctx, key, string(data), expiration)
}

func (r *urlRepository) saveToMemory(ctx context.Context, url *model.ShortURL) error {
	key := "shorturl:" + url.ShortCode
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	var expiration time.Duration
	if url.ExpiresAt != nil {
		expiration = time.Until(*url.ExpiresAt)
		if expiration < 0 {
			return fmt.Errorf("url already expired")
		}
	}

	return r.sources.MemoryCache.Set(ctx, key, string(data), expiration)
}

func (r *urlRepository) saveToMySQL(ctx context.Context, url *model.ShortURL) error {
	db := r.sources.MySQLDB.GetDB()
	query := `INSERT INTO short_urls (short_code, long_url, created_at, expires_at) 
	          VALUES (?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE long_url = ?, expires_at = ?`

	var expiresAt interface{}
	if url.ExpiresAt != nil {
		expiresAt = *url.ExpiresAt
	}

	_, err := db.ExecContext(ctx, query,
		url.ShortCode, url.LongURL, url.CreatedAt, expiresAt,
		url.LongURL, expiresAt,
	)
	return err
}

func (r *urlRepository) saveToSQLite(ctx context.Context, url *model.ShortURL) error {
	db := r.sources.SQLiteDB.GetDB()
	query := `INSERT INTO short_urls (short_code, long_url, created_at, expires_at) 
	          VALUES (?, ?, ?, ?) 
	          ON CONFLICT(short_code) DO UPDATE SET long_url = ?, expires_at = ?`

	var expiresAt interface{}
	if url.ExpiresAt != nil {
		expiresAt = *url.ExpiresAt
	}

	_, err := db.ExecContext(ctx, query,
		url.ShortCode, url.LongURL, url.CreatedAt, expiresAt,
		url.LongURL, expiresAt,
	)
	return err
}

// 异步回写缓存
func (r *urlRepository) asyncWriteToCache(ctx context.Context, url *model.ShortURL) {
	if r.sources.RedisCache != nil {
		_ = r.saveToRedis(ctx, url)
	} else if r.sources.MemoryCache != nil {
		_ = r.saveToMemory(ctx, url)
	}
}

// 实现 URLRepository 接口的旧方法（保持兼容性）

func (r *urlRepository) SaveToCache(ctx context.Context, url *model.ShortURL) error {
	return r.Save(ctx, url)
}

func (r *urlRepository) GetFromCache(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	return r.Get(ctx, shortCode)
}

func (r *urlRepository) SaveToDB(ctx context.Context, url *model.ShortURL) error {
	return r.Save(ctx, url)
}

func (r *urlRepository) GetFromDB(ctx context.Context, shortCode string) (*model.ShortURL, error) {
	return r.Get(ctx, shortCode)
}

func (r *urlRepository) DeleteFromCache(ctx context.Context, shortCode string) error {
	key := "shorturl:" + shortCode
	var errs []error

	if r.sources.RedisCache != nil {
		if err := r.sources.RedisCache.Delete(ctx, key); err != nil {
			errs = append(errs, err)
		}
	}

	if r.sources.MemoryCache != nil {
		if err := r.sources.MemoryCache.Delete(ctx, key); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete from cache: %v", errs)
	}

	return nil
}

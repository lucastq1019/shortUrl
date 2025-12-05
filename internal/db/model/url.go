package model

import "time"

// ShortURL 短链接模型
type ShortURL struct {
	ID        int64     `json:"id"`
	ShortCode string    `json:"short_code"` // 短码
	LongURL   string    `json:"long_url"`   // 原始长链接
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // 可选：过期时间
}


package utils

import (
	"net/url"
	"strings"
)

// ValidateURL 验证 URL 是否有效
func ValidateURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	// 如果没有协议，添加 http://
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

// IsValidShortCode 验证短码是否有效（只包含字母数字和连字符）
func IsValidShortCode(code string) bool {
	if len(code) < 4 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}

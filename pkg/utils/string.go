package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateShortCode 生成随机短码
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", &InvalidLengthError{Length: length}
	}
	
	// 计算需要的字节数（base64 编码后长度约为原长度的 4/3）
	bytesNeeded := (length*3 + 3) / 4
	bytes := make([]byte, bytesNeeded)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// InvalidLengthError 无效长度错误
type InvalidLengthError struct {
	Length int
}

func (e *InvalidLengthError) Error() string {
	return "invalid length: must be greater than 0"
}


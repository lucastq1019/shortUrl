package utils

import (
	"testing"
)

// ========== 测试函数示例 ==========

// TestValidateURL 测试 URL 验证
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    bool
	}{
		{
			name: "有效的 HTTP URL",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "有效的 HTTPS URL",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "没有协议的 URL（应该自动添加 http://）",
			url:  "example.com",
			want: true,
		},
		{
			name: "空字符串",
			url:  "",
			want: false,
		},
		{
			name: "无效的 URL",
			url:  "not a url",
			want: false,
		},
		{
			name: "带路径的 URL",
			url:  "https://example.com/path/to/page",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateURL(tt.url)
			if got != tt.want {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

// TestIsValidShortCode 测试短码验证
func TestIsValidShortCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "有效的短码（字母数字）",
			code: "abc123",
			want: true,
		},
		{
			name: "有效的短码（包含连字符）",
			code: "abc-123",
			want: true,
		},
		{
			name: "有效的短码（包含下划线）",
			code: "abc_123",
			want: true,
		},
		{
			name: "太短的短码",
			code: "abc",
			want: false,
		},
		{
			name: "太长的短码",
			code: "abcdefghijklmnopqrstuvwxyz1234567890",
			want: false,
		},
		{
			name: "包含无效字符",
			code: "abc@123",
			want: false,
		},
		{
			name: "空字符串",
			code: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidShortCode(tt.code)
			if got != tt.want {
				t.Errorf("IsValidShortCode(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// ========== 基准测试示例 ==========

// BenchmarkValidateURL 基准测试 URL 验证
func BenchmarkValidateURL(b *testing.B) {
	testURL := "https://example.com/path/to/page?query=value"
	
	b.ResetTimer() // 重置计时器，排除准备时间
	for i := 0; i < b.N; i++ {
		_ = ValidateURL(testURL)
	}
}

// BenchmarkValidateURL_NoProtocol 测试没有协议的 URL
func BenchmarkValidateURL_NoProtocol(b *testing.B) {
	testURL := "example.com/path"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateURL(testURL)
	}
}

// BenchmarkIsValidShortCode 基准测试短码验证
func BenchmarkIsValidShortCode(b *testing.B) {
	testCode := "abc123XYZ"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsValidShortCode(testCode)
	}
}

// BenchmarkIsValidShortCode_Invalid 测试无效短码的性能
func BenchmarkIsValidShortCode_Invalid(b *testing.B) {
	testCode := "abc@123" // 包含无效字符
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsValidShortCode(testCode)
	}
}

// ========== 并行测试示例 ==========

// TestValidateURL_Parallel 并行测试
func TestValidateURL_Parallel(t *testing.T) {
	tests := []string{
		"http://example.com",
		"https://example.com",
		"example.com",
	}

	for _, url := range tests {
		url := url // 重要：避免闭包问题
		t.Run(url, func(t *testing.T) {
			t.Parallel() // 标记为并行测试
			_ = ValidateURL(url)
		})
	}
}

// ========== 子测试和清理示例 ==========

// TestWithCleanup 演示测试清理
func TestWithCleanup(t *testing.T) {
	// 设置测试环境
	t.Log("Setting up test environment")
	
	// 注册清理函数
	t.Cleanup(func() {
		t.Log("Cleaning up test environment")
		// 这里可以清理资源，如关闭文件、删除临时文件等
	})
	
	// 执行测试
	t.Log("Running test")
	// 测试代码...
}


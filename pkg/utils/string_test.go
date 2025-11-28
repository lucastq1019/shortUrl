package utils

import (
	"testing"
)

// ========== 测试函数示例 ==========

// TestGenerateShortCode 测试 GenerateShortCode 函数
// 函数名必须以 Test 开头，参数必须是 *testing.T
func TestGenerateShortCode(t *testing.T) {
	// 测试用例1：正常生成短码
	code, err := GenerateShortCode(6)
	if err != nil {
		t.Fatalf("GenerateShortCode failed: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("Expected length 6, got %d", len(code))
	}

	// 测试用例2：不同长度
	code2, err := GenerateShortCode(10)
	if err != nil {
		t.Fatalf("GenerateShortCode failed: %v", err)
	}
	if len(code2) != 10 {
		t.Errorf("Expected length 10, got %d", len(code2))
	}

	// 测试用例3：验证生成的短码不重复（概率上）
	code3, _ := GenerateShortCode(8)
	if code2 == code3 {
		t.Error("Generated codes should be different")
	}
}

// TestGenerateShortCode_InvalidLength 测试无效长度
// 可以用下划线分隔测试场景
func TestGenerateShortCode_InvalidLength(t *testing.T) {
	// 测试负数长度
	_, err := GenerateShortCode(-1)
	if err == nil {
		t.Error("Expected error for negative length, got nil")
	}

	// 测试零长度
	_, err = GenerateShortCode(0)
	if err == nil {
		t.Error("Expected error for zero length, got nil")
	}
}

// TestGenerateShortCode_TableDriven 表驱动测试（推荐方式）
func TestGenerateShortCode_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
		wantErr bool
	}{
		{
			name:    "正常长度 6",
			length:  6,
			wantLen: 6,
			wantErr: false,
		},
		{
			name:    "正常长度 10",
			length:  10,
			wantLen: 10,
			wantErr: false,
		},
		{
			name:    "正常长度 16",
			length:  16,
			wantLen: 16,
			wantErr: false,
		},
		{
			name:    "负数长度",
			length:  -1,
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "零长度",
			length:  0,
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// 使用 t.Run 运行子测试
		t.Run(tt.name, func(t *testing.T) {
			code, err := GenerateShortCode(tt.length)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShortCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && len(code) != tt.wantLen {
				t.Errorf("GenerateShortCode() length = %d, want %d", len(code), tt.wantLen)
			}
		})
	}
}

// ========== 基准测试示例 ==========

// BenchmarkGenerateShortCode 基准测试 GenerateShortCode 函数
// 函数名必须以 Benchmark 开头，参数必须是 *testing.B
func BenchmarkGenerateShortCode(b *testing.B) {
	// b.N 是 Go 自动调整的迭代次数
	for i := 0; i < b.N; i++ {
		_, err := GenerateShortCode(8)
		if err != nil {
			b.Fatalf("GenerateShortCode failed: %v", err)
		}
	}
}

// BenchmarkGenerateShortCode_Length6 测试不同长度的性能
func BenchmarkGenerateShortCode_Length6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortCode(6)
	}
}

// BenchmarkGenerateShortCode_Length10 测试不同长度的性能
func BenchmarkGenerateShortCode_Length10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortCode(10)
	}
}

// BenchmarkGenerateShortCode_Length16 测试不同长度的性能
func BenchmarkGenerateShortCode_Length16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortCode(16)
	}
}

// ========== 辅助测试函数示例 ==========

// ExampleGenerateShortCode 示例函数（会在文档中显示）
// 函数名必须以 Example 开头
// 运行 go doc 时会显示这个示例
func ExampleGenerateShortCode() {
	code, err := GenerateShortCode(8)
	if err != nil {
		return
	}
	_ = code
	// 如果函数有输出，可以用 Output 注释验证
	// fmt.Println(code)
	// Output: abc123XY
	// 注意：如果不需要验证输出，可以删除 Output 注释
}


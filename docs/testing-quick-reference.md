# Go 测试快速参考

## 测试文件命名

- 测试文件：`xxx_test.go`
- 测试函数：`TestXxx(t *testing.T)`
- 基准测试：`BenchmarkXxx(b *testing.B)`

## 常用命令

### 运行测试

```bash
# 运行所有测试
go test

# 运行并显示详细信息
go test -v

# 运行特定测试
go test -run TestValidateURL

# 运行所有包的测试
go test ./...

# 显示测试覆盖率
go test -cover
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -bench=.

# 运行并显示内存分配
go test -bench=. -benchmem

# 运行特定基准测试
go test -bench=BenchmarkGenerateShortCode

# 运行多次取平均值
go test -bench=. -count=5
```

### 测试覆盖率

```bash
# 生成覆盖率文件
go test -coverprofile=coverage.out

# 查看 HTML 报告
go tool cover -html=coverage.out

# 查看文本报告
go tool cover -func=coverage.out
```

## 测试函数模板

### 基本测试

```go
func TestFunction(t *testing.T) {
    result := Function(input)
    if result != expected {
        t.Errorf("Function() = %v, want %v", result, expected)
    }
}
```

### 表驱动测试

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        want bool
    }{
        {"case1", "input1", true},
        {"case2", "input2", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Function(tt.input)
            if got != tt.want {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 基准测试

```go
func BenchmarkFunction(b *testing.B) {
    data := prepareData()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = Function(data)
    }
}
```

## testing.T 常用方法

- `t.Log()` - 记录日志
- `t.Error()` - 标记失败但继续
- `t.Fatal()` - 标记失败并停止
- `t.Run()` - 运行子测试
- `t.Parallel()` - 并行测试
- `t.Cleanup()` - 注册清理函数

## testing.B 常用方法

- `b.ResetTimer()` - 重置计时器
- `b.StopTimer()` - 停止计时器
- `b.StartTimer()` - 开始计时器
- `b.ReportAllocs()` - 报告内存分配


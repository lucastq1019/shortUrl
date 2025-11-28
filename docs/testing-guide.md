# Go 测试指南

## 测试文件命名

- 测试文件必须以 `_test.go` 结尾
- 测试文件必须和要测试的文件在同一个包中
- 例如：`string.go` 的测试文件是 `string_test.go`

```
pkg/utils/
├── string.go          # 源代码
├── string_test.go     # 测试文件
├── validator.go       # 源代码
└── validator_test.go  # 测试文件
```

## 测试函数（Test Functions）

### 基本语法

```go
func TestFunctionName(t *testing.T) {
    // 测试代码
}
```

**规则**：
- 函数名必须以 `Test` 开头
- 参数必须是 `*testing.T`
- 文件必须以 `_test.go` 结尾

### 测试函数示例

```go
func TestGenerateShortCode(t *testing.T) {
    code, err := GenerateShortCode(6)
    if err != nil {
        t.Fatalf("GenerateShortCode failed: %v", err)
    }
    if len(code) != 6 {
        t.Errorf("Expected length 6, got %d", len(code))
    }
}
```

### 常用的 testing.T 方法

| 方法 | 说明 |
|------|------|
| `t.Log()` | 记录日志（只在失败或 `-v` 时显示） |
| `t.Logf()` | 格式化日志 |
| `t.Error()` | 标记测试失败，但继续执行 |
| `t.Errorf()` | 格式化错误消息 |
| `t.Fatal()` | 标记测试失败并立即停止 |
| `t.Fatalf()` | 格式化致命错误 |
| `t.Skip()` | 跳过测试 |
| `t.Run()` | 运行子测试 |
| `t.Parallel()` | 标记为并行测试 |
| `t.Cleanup()` | 注册清理函数 |

### 表驱动测试（推荐）

```go
func TestValidateURL(t *testing.T) {
    tests := []struct {
        name string
        url  string
        want bool
    }{
        {
            name: "有效的 HTTP URL",
            url:  "http://example.com",
            want: true,
        },
        {
            name: "无效的 URL",
            url:  "not a url",
            want: false,
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
```

## 基准测试（Benchmark Functions）

### 基本语法

```go
func BenchmarkFunctionName(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 要测试的代码
    }
}
```

**规则**：
- 函数名必须以 `Benchmark` 开头
- 参数必须是 `*testing.B`
- 必须包含 `for i := 0; i < b.N; i++` 循环
- `b.N` 是 Go 自动调整的迭代次数

### 基准测试示例

```go
func BenchmarkGenerateShortCode(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, err := GenerateShortCode(8)
        if err != nil {
            b.Fatalf("GenerateShortCode failed: %v", err)
        }
    }
}
```

### 常用的 testing.B 方法

| 方法 | 说明 |
|------|------|
| `b.ResetTimer()` | 重置计时器，排除准备时间 |
| `b.StopTimer()` | 停止计时器 |
| `b.StartTimer()` | 开始计时器 |
| `b.ReportAllocs()` | 报告内存分配 |
| `b.RunParallel()` | 并行基准测试 |

### 排除准备时间

```go
func BenchmarkValidateURL(b *testing.B) {
    // 准备数据（不计入基准测试时间）
    testURL := "https://example.com/path"
    
    b.ResetTimer() // 重置计时器
    for i := 0; i < b.N; i++ {
        _ = ValidateURL(testURL)
    }
}
```

## 执行测试

### 运行所有测试

```bash
# 运行当前包的所有测试
go test

# 运行所有包的测试
go test ./...

# 运行并显示详细信息
go test -v

# 运行并显示覆盖率
go test -cover
```

### 运行特定测试

```bash
# 运行匹配的测试函数
go test -run TestGenerateShortCode

# 使用正则表达式
go test -run "TestGenerate.*"

# 运行特定包的测试
go test ./pkg/utils
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -bench=.

# 运行匹配的基准测试
go test -bench=BenchmarkGenerateShortCode

# 运行基准测试并显示详细信息
go test -bench=. -benchmem

# 运行基准测试并显示内存分配
go test -bench=. -benchmem -memprofile=mem.prof

# 运行基准测试并生成 CPU profile
go test -bench=. -cpuprofile=cpu.prof
```

### 常用测试标志

| 标志 | 说明 |
|------|------|
| `-v` | 显示详细输出 |
| `-run` | 运行匹配的测试函数（支持正则） |
| `-bench` | 运行匹配的基准测试（支持正则） |
| `-benchmem` | 显示内存分配统计 |
| `-cover` | 显示测试覆盖率 |
| `-coverprofile` | 生成覆盖率文件 |
| `-count` | 运行测试的次数 |
| `-timeout` | 设置超时时间 |
| `-parallel` | 设置并行测试数 |

### 测试覆盖率

```bash
# 显示覆盖率
go test -cover

# 生成覆盖率文件
go test -coverprofile=coverage.out

# 查看覆盖率详情（HTML）
go tool cover -html=coverage.out

# 查看覆盖率详情（文本）
go tool cover -func=coverage.out
```

### 并行测试

```bash
# 设置并行测试数
go test -parallel 4
```

## 示例命令

### 1. 运行所有测试（详细输出）

```bash
go test -v ./...
```

### 2. 运行特定包的测试

```bash
go test -v ./pkg/utils
```

### 3. 运行特定测试函数

```bash
go test -v -run TestValidateURL
```

### 4. 运行所有基准测试

```bash
go test -bench=. -benchmem
```

### 5. 运行特定基准测试

```bash
go test -bench=BenchmarkGenerateShortCode -benchmem
```

### 6. 比较不同实现的性能

```bash
# 运行基准测试并比较
go test -bench=. -benchmem > old.txt
# 修改代码后
go test -bench=. -benchmem > new.txt
# 使用 benchstat 比较（需要安装：go install golang.org/x/perf/cmd/benchstat@latest）
benchstat old.txt new.txt
```

### 7. 生成测试覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 测试最佳实践

### 1. 使用表驱动测试

```go
// ✅ 推荐
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        want bool
    }{
        // 测试用例
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试代码
        })
    }
}
```

### 2. 测试命名清晰

```go
// ✅ 推荐
func TestValidateURL_InvalidURL(t *testing.T) {}
func TestGenerateShortCode_EmptyString(t *testing.T) {}

// ❌ 不推荐
func Test1(t *testing.T) {}
func TestURL(t *testing.T) {}
```

### 3. 使用 t.Run 组织子测试

```go
func TestFunction(t *testing.T) {
    t.Run("case1", func(t *testing.T) {
        // 子测试1
    })
    t.Run("case2", func(t *testing.T) {
        // 子测试2
    })
}
```

### 4. 基准测试要准确

```go
func BenchmarkFunction(b *testing.B) {
    // 准备数据
    data := prepareData()
    
    b.ResetTimer() // 重置计时器
    for i := 0; i < b.N; i++ {
        // 只测试要测量的部分
        _ = Function(data)
    }
}
```

### 5. 使用 t.Cleanup 清理资源

```go
func TestWithFile(t *testing.T) {
    file := createTempFile(t)
    t.Cleanup(func() {
        os.Remove(file.Name())
    })
    // 测试代码
}
```

## 示例输出

### 测试输出示例

```
=== RUN   TestValidateURL
=== RUN   TestValidateURL/有效的_HTTP_URL
=== RUN   TestValidateURL/无效的_URL
--- PASS: TestValidateURL (0.00s)
    --- PASS: TestValidateURL/有效的_HTTP_URL (0.00s)
    --- PASS: TestValidateURL/无效的_URL (0.00s)
PASS
ok      github.com/username/myproject/pkg/utils    0.123s
```

### 基准测试输出示例

```
goos: darwin
goarch: amd64
pkg: github.com/username/myproject/pkg/utils
BenchmarkGenerateShortCode-8         1000000    1234 ns/op    256 B/op    2 allocs/op
PASS
ok      github.com/username/myproject/pkg/utils    1.234s
```

输出说明：
- `1000000` - 迭代次数
- `1234 ns/op` - 每次操作耗时（纳秒）
- `256 B/op` - 每次操作内存分配（字节）
- `2 allocs/op` - 每次操作内存分配次数

## 参考资源

- [Go Testing Package](https://pkg.go.dev/testing)
- [Effective Go - Testing](https://go.dev/doc/effective_go#testing)
- [Go Blog - The cover story](https://go.dev/blog/cover)


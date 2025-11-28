# Go 语言命名约定规范

## 核心原则

### 1. 可见性规则（最重要！）
- **首字母大写** = 公开（Public/Exported），可被其他包导入使用
- **首字母小写** = 私有（Private/Unexported），只能在当前包内使用

```go
// 公开的，其他包可以导入使用
func PublicFunction() {}

// 私有的，只能在当前包内使用
func privateFunction() {}
```

## 命名规范详解

### 1. 包名（Package Name）
- **全小写**，简短，通常是一个单词
- 不要使用下划线或混合大小写
- 不要使用复数形式

```go
// ✅ 正确
package handler
package service
package utils

// ❌ 错误
package handlers  // 不要用复数
package myHandler // 不要混合大小写
package my_handler // 不要用下划线
```

### 2. 函数名（Function Name）
- **驼峰命名**（CamelCase）
- 公开函数：首字母大写
- 私有函数：首字母小写
- 函数名应该是动词或动词短语

```go
// ✅ 正确
func GetUser() {}        // 公开函数
func getUser() {}        // 私有函数
func CreateShortURL() {} // 公开函数
func validateURL() {}    // 私有函数

// ❌ 错误
func get_user() {}       // 不要用下划线
func Get_User() {}       // 不要用下划线
```

### 3. 方法名（Method Name）
- 与函数名规则相同
- 接收者类型通常用 1-2 个字母的缩写

```go
// ✅ 正确
func (s *Service) CreateURL() {}  // 公开方法
func (r *Repository) save() {}     // 私有方法
func (h *Handler) HandleRequest() {}

// 接收者命名约定
func (u *User) GetName() {}        // u = user
func (s *Service) Process() {}     // s = service
func (db *DB) Query() {}           // db = database
```

### 4. 变量名（Variable Name）
- **驼峰命名**
- 公开变量：首字母大写（很少使用，通常用函数代替）
- 私有变量：首字母小写
- 简短但有意义
- 布尔变量通常以 `is`, `has`, `can`, `should` 开头

```go
// ✅ 正确
var maxRetries = 3              // 私有变量
var MaxRetries = 3              // 公开变量（不推荐，用函数代替）
var userCount int               // 私有变量
var isActive bool               // 布尔变量
var hasPermission bool          // 布尔变量
var canEdit bool                // 布尔变量

// ❌ 错误
var max_retries = 3             // 不要用下划线
var Max_Retries = 3             // 不要用下划线
```

### 5. 常量（Constants）
- **两种风格都可以**：
  - 全大写 + 下划线（类似 C 风格）
  - 驼峰命名（Go 推荐，更常见）

```go
// ✅ 风格1：全大写+下划线（用于全局常量）
const (
    MAX_RETRIES = 3
    DEFAULT_PORT = 8080
    API_VERSION = "v1"
)

// ✅ 风格2：驼峰命名（Go 推荐，更常见）
const (
    MaxRetries = 3
    DefaultPort = 8080
    APIVersion = "v1"
)

// ✅ 混合使用（根据用途）
const (
    // 公开常量用驼峰
    DefaultTimeout = 30 * time.Second
    
    // 私有常量用小写
    maxConnections = 100
)
```

### 6. 类型名（Type Name）
- **驼峰命名**，首字母大写（通常是公开的）
- 接口名通常以 `-er` 结尾

```go
// ✅ 正确
type User struct {}
type ShortURL struct {}
type URLRepository interface {}  // 接口
type URLService interface {}     // 接口

// 接口命名（推荐以 -er 结尾）
type Reader interface {}
type Writer interface {}
type URLRepository interface {}  // Repository 本身就是 -er 结尾

// ❌ 错误
type user struct {}              // 如果要在包外使用，应该大写
type short_url struct {}         // 不要用下划线
```

### 7. 接口名（Interface Name）
- 通常以 `-er` 结尾
- 或者使用描述性的名词

```go
// ✅ 正确
type Reader interface {}
type Writer interface {}
type Closer interface {}
type URLRepository interface {}  // Repository 本身就是 -er
type URLService interface {}     // Service 也可以

// 或者描述性命名
type HTTPClient interface {}
type Database interface {}
```

### 8. 错误变量（Error Variables）
- 通常以 `Err` 开头
- 使用 `errors.New()` 或 `fmt.Errorf()`

```go
// ✅ 正确
var ErrNotFound = errors.New("not found")
var ErrInvalidURL = errors.New("invalid URL")
var ErrDatabaseConnection = errors.New("database connection failed")

// 使用
if err != nil {
    return ErrNotFound
}
```

### 9. 缩写词
- 缩写词要么全大写，要么全小写，不要混合

```go
// ✅ 正确
var userID int          // ID 全大写
var apiKey string       // API 全大写
var httpClient *Client  // HTTP 全大写
var urlPath string      // URL 全大写

// ❌ 错误
var userId int          // 应该是 userID
var api_key string      // 应该是 apiKey
var http_client *Client // 应该是 httpClient
```

### 10. 局部变量
- 简短命名，上下文清晰时可使用单字母
- 常见单字母变量：`i`, `j`, `k`（循环），`r`（reader），`w`（writer），`err`（error）

```go
// ✅ 正确
for i := 0; i < 10; i++ {}
for _, item := range items {}
if err != nil { return err }

// 上下文清晰时可以用短名
func (u *User) GetName() string {
    return u.name  // u 是接收者，上下文清晰
}
```

## 命名示例对比

### 好的命名 vs 坏的命名

```go
// ✅ 好的命名
func CreateShortURL(longURL string) (string, error) {}
func (s *Service) GetUserByID(id int64) (*User, error) {}
var maxRetries = 3
const DefaultTimeout = 30 * time.Second
type URLRepository interface {}

// ❌ 坏的命名
func create_short_url(long_url string) (string, error) {}  // 下划线
func (s *Service) GetUserById(id int64) (*User, error) {}   // ID 应该是大写
var Max_Retries = 3                                          // 下划线
const DEFAULT_TIMEOUT = 30 * time.Second                     // 风格不一致
type url_repository interface {}                             // 下划线，小写
```

## 特殊命名约定

### 1. Getter/Setter（不推荐使用 Get/Set 前缀）
```go
// ❌ 不推荐（Java 风格）
func (u *User) GetName() string { return u.name }
func (u *User) SetName(name string) { u.name = name }

// ✅ 推荐（Go 风格）
func (u *User) Name() string { return u.name }
func (u *User) SetName(name string) { u.name = name }
```

### 2. 构造函数
```go
// ✅ 推荐使用 New 前缀
func NewUser(name string) *User {}
func NewService(repo Repository) *Service {}
func NewRepository(db *sql.DB) *Repository {}
```

### 3. 测试函数
```go
// 测试函数必须以 Test 开头
func TestCreateURL(t *testing.T) {}
func TestGetUser_NotFound(t *testing.T) {}  // 可以用下划线分隔场景
```

### 4. 基准测试
```go
// 基准测试必须以 Benchmark 开头
func BenchmarkCreateURL(b *testing.B) {}
```

## 总结

1. **可见性最重要**：大写=公开，小写=私有
2. **使用驼峰命名**：不要用下划线（除了常量可选）
3. **简短但有意义**：`i` 可以，但 `index` 更清晰
4. **一致性**：整个项目保持一致的命名风格
5. **遵循 Go 习惯**：多看标准库的命名方式

## 参考

- [Effective Go - Names](https://go.dev/doc/effective_go#names)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)


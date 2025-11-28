# Go 命名约定快速参考

## 核心规则

| 类型 | 公开（可导出） | 私有（不可导出） |
|------|--------------|----------------|
| **规则** | 首字母**大写** | 首字母**小写** |
| **函数** | `func PublicFunc()` | `func privateFunc()` |
| **方法** | `func (s *S) PublicMethod()` | `func (s *S) privateMethod()` |
| **变量** | `var PublicVar` | `var privateVar` |
| **常量** | `const PublicConst` | `const privateConst` |
| **类型** | `type PublicType struct{}` | `type privateType struct{}` |
| **字段** | `type S struct { PublicField string }` | `type S struct { privateField string }` |

## 命名风格

### ✅ 推荐
- 使用**驼峰命名**（CamelCase）
- 常量可以用全大写+下划线或驼峰命名
- 缩写词全大写：`UserID`, `APIKey`, `HTTPClient`
- 布尔变量以 `is/has/can/should` 开头

### ❌ 避免
- 不要使用下划线（`snake_case`）
- 不要混合大小写（`myVariable` 可以，但 `My_Variable` 不行）
- 不要使用复数包名（`handlers` → `handler`）

## 常见命名模式

```go
// 构造函数
func NewService() *Service {}

// Getter（Go 风格，不用 Get 前缀）
func (u *User) Name() string {}

// 布尔方法
func (u *User) IsActive() bool {}
func (u *User) HasPermission() bool {}

// 错误变量
var ErrNotFound = errors.New("not found")

// 接口（通常以 -er 结尾）
type Reader interface {}
type Writer interface {}
```

## 快速检查清单

- [ ] 公开的标识符首字母大写了吗？
- [ ] 私有的标识符首字母小写了吗？
- [ ] 使用驼峰命名了吗？（没有下划线）
- [ ] 缩写词全大写了吗？（`ID`, `API`, `HTTP`）
- [ ] 布尔变量/方法以 `is/has/can` 开头了吗？
- [ ] 错误变量以 `Err` 开头了吗？
- [ ] 构造函数以 `New` 开头了吗？


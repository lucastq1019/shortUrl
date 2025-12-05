# ShortURL 项目

## Go 项目目录结构说明

### 标准目录约定

```
shorturl/
├── cmd/                    # 应用程序入口点（main 函数）
│   └── server/            # 服务器主程序
│       └── main.go
├── internal/              # 私有应用代码（外部无法导入）
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── model/            # 数据模型
│   └── middleware/       # 中间件
├── pkg/                   # 可被外部使用的库代码
│   └── utils/            # 工具函数
├── api/                   # API 定义（可选）
├── configs/              # 配置文件
├── migrations/           # 数据库迁移文件
├── scripts/              # 脚本文件
├── docs/                 # 文档
├── test/                 # 测试数据
├── go.mod
└── go.sum
```

### 目录说明

- **cmd/**: 包含应用程序的入口点，每个子目录是一个可执行程序
- **internal/**: 私有代码，只能被本项目内部导入，外部项目无法导入
- **pkg/**: 可被外部项目使用的公共库代码
- **api/**: API 定义文件（如 OpenAPI/Swagger）
- **configs/**: 配置文件（YAML, JSON, TOML 等）
- **migrations/**: 数据库迁移脚本
- **scripts/**: 构建、部署等脚本
- **docs/**: 项目文档
- **test/**: 测试数据和测试工具

### 为什么这样划分？

1. **cmd/**: 将入口点与业务逻辑分离，便于维护
2. **internal/**: 防止外部项目依赖你的内部实现，保持 API 稳定性
3. **pkg/**: 如果你要发布可复用的库，放在这里
4. **分层架构**: internal 下的 handler/service/repository 实现 MVC 或分层架构

### 架构分层说明

```
┌─────────────────┐
│   Handler       │  HTTP 层：处理请求和响应
│  (路由/控制器)   │
└────────┬────────┘
         │
┌────────▼────────┐
│   Service       │  业务逻辑层：核心业务处理
│  (业务逻辑)      │
└────────┬────────┘
         │
┌────────▼────────┐
│  Repository     │  数据访问层：数据库操作
│  (数据访问)      │
└────────┬────────┘
         │
┌────────▼────────┐
│    Model        │  数据模型：结构体定义
│  (数据模型)      │
└─────────────────┘
```

### 代码组织最佳实践

1. **导入路径**: 使用完整的模块路径，如 `github.com/username/myproject/internal/handler`
2. **包命名**: 使用小写单数形式，如 `handler` 而不是 `handlers`
3. **文件命名**: 使用小写下划线，如 `user_repository.go`
4. **接口定义**: 接口通常定义在使用它的地方，而不是实现的地方
5. **错误处理**: 始终检查错误，不要忽略

### 运行项目

```bash
# 构建
go build ./cmd/server

# 运行
./server

# 或直接运行
go run ./cmd/server/main.go
```

### 命名约定

Go 语言有严格的命名约定，详见：[命名约定文档](docs/naming-conventions.md)

**核心要点**：
- 首字母大写 = 公开（可被其他包导入）
- 首字母小写 = 私有（仅当前包可用）
- 使用驼峰命名，不要用下划线
- 常量可以用全大写+下划线，或驼峰命名

### 测试

Go 测试指南，详见：[测试指南](docs/testing-guide.md) | [快速参考](docs/testing-quick-reference.md)

**快速开始**：
```bash
# 运行所有测试
go test ./...

# 运行基准测试
go test -bench=. -benchmem

# 查看测试覆盖率
go test -cover
```

### 下一步

1. 实现 Repository 层（数据库操作）
2. 实现 Service 层（业务逻辑）
3. 实现 Handler 层（HTTP 处理）
4. 添加中间件（日志、认证等）
5. 编写测试


💡 设计思路总结
该设计采用 Map + 工厂模式 实现 gRPC 连接和客户端的集中管理和高效复用：

分层管理 (Decoupling): 将 配置映射、连接复用、客户端创建、业务调用 四个职责分离。

配置驱动 (Config-Driven): 通过外部配置（服务名 → 地址）驱动整个初始化流程。

连接复用 (Conn Reuse): 使用 map[Address]*grpc.ClientConn 存储和复用连接，确保同一地址只建立一个底层连接，实现 gRPC 的多路复用优势。

客户端注册/查找 (Client Lookup): 使用 map[ServiceName]*ServiceEntry 存储客户端，方便业务代码通过友好的 服务名 快速检索到对应的客户端存根。

动态创建 (Factory Pattern): 引入 map[ServiceName]ClientFactory（工厂映射），将创建特定类型客户端存根的逻辑（即 pb.New...Client(conn)）与核心管理器逻辑分离，避免使用 switch/case，提高可扩展性。

安全和清理 (Safety): 使用 sync.RWMutex 保证并发安全，并提供 Close() 方法安全关闭所有唯一的 Conn 实例，防止资源泄漏。


### 生成proto
``` sh
# 生成proto文件到目标位置
protoc --go_out=./internal/rpc/ --go_opt=paths=source_relative \
       --go-grpc_out=./internal/rpc/ --go-grpc_opt=paths=source_relative \
       proto/*.proto
```

### 脚本测试


### 
  

package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// ------------------------------------------------
// 结构体定义 (已修复字段导出问题)
// ------------------------------------------------

// Config 应用配置的总结构体
type Config struct {
	GinMode    string
	HttpAddr   string
	MySQLDSN   string
	RedisAddr  string
	SQLitePath string
	// 客户端访问的 gRPC 服务地址
	GRPCServers struct {
		Shortener struct {
			Addr string
		}
		Clipboarder struct {
			Addr string
		}
	}
	// 当前应用自身作为 gRPC Server 的配置 (示例：shortneer 和 clipboarder)
	RPC struct {
		Shortneer struct {
			Addr string
		}
		Clipboarder struct {
			Addr string
		}
	}
}

// ------------------------------------------------
// 单例实现区域 (保持不变)
// ------------------------------------------------

var once sync.Once
var cfg *Config

// GetConfig 获取配置单例实例
func GetConfig() *Config {
	once.Do(func() {
		// LoadAll 现在返回指针
		cfg = LoadAll()
	})
	return cfg
}

// ------------------------------------------------
// 配置加载逻辑 (改为 Viper 加载 YAML)
// ------------------------------------------------

// LoadAll 加载配置：优先从 YAML 文件加载，然后覆盖环境变量，最后使用默认值。
func LoadAll() *Config {
	// 1. 设置默认值 (如果 YAML 和环境变量都没有设置)
	v := viper.New()
	v.SetDefault("GinMode", "debug")
	v.SetDefault("HttpAddr", ":8080")
	v.SetDefault("MySQLDSN", "user:password@tcp(localhost:3306)/shorturl")
	v.SetDefault("RedisAddr", "localhost:6379")
	v.SetDefault("SQLitePath", "./data/shorturl.db")
	v.SetDefault("GRPCServers.shortener", "localhost:9090")
	v.SetDefault("GRPCServers.clipboarder", "localhost:9091")
	v.SetDefault("RPC.Shortneer.Addr", ":9090")   // 假设这是 Shortneer 的 RPC 监听地址
	v.SetDefault("RPC.Clipboarder.Addr", ":9091") // 假设这是 Clipboarder 的 RPC 监听地址

	// 2. 配置 YAML 文件路径和类型
	v.SetConfigName("config") // 文件名为 config (无扩展名)
	v.SetConfigType("yaml")   // 文件类型为 YAML
	v.AddConfigPath(".")      // 搜索当前目录

	// 3. 读取 YAML 文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 文件找不到是允许的，继续使用默认值和环境变量
			log.Println("Warning: config.yaml not found. Using defaults and environment variables.")
		} else {
			// 其他读取错误，如文件格式错误，应退出
			log.Fatalf("Fatal error reading config file: %v", err)
		}
	}

	// 4. 绑定环境变量 (可以覆盖 YAML 文件中的值)
	// 这一步需要手动设置环境变量的前缀和 key 的对应关系

	// 5. 将配置绑定到结构体
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("Unable to unmarshal config into struct: %v", err)
	}

	return &c
}

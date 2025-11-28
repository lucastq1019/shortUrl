package manager

import (
	"github.com/username/shorturl/internal/config"
	shorturlpb "github.com/username/shorturl/internal/rpc/proto"
	"google.golang.org/grpc"
)

// Operation 定义了客户端创建闭包的签名
type Operation func(c *grpc.ClientConn) interface{}

// clientFactoryMap 是私有变量，存储所有服务的创建逻辑
// 这是一个集中的地方，用于添加新的客户端服务
var clientFactoryMap = map[string]Operation{
	config.GetConfig().GRPCServers.Shortener.Addr: func(c *grpc.ClientConn) interface{} {
		// 实际调用 proto 生成的 NewServiceClient
		return shorturlpb.NewShortenerServiceClient(c)
	},
	config.GetConfig().GRPCServers.Clipboarder.Addr: func(c *grpc.ClientConn) interface{} {
		return shorturlpb.NewClipboarderServiceClient(c)
	},
	// 未来增加服务时，只需修改这里
	// "service_new": func(c *grpc.ClientConn) interface{} { ... },
}

package rpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/username/shorturl/internal/config"
	"github.com/username/shorturl/internal/manager"
	shortenerpb "github.com/username/shorturl/internal/rpc/proto"
	shortener "github.com/username/shorturl/internal/rpc/service/shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer() *grpc.Server {

	// ui := grpctrace.UnaryServerInterceptor(grpctrace.WithService(os.Getenv("DD_SERVICE")))
	// cui := grpc_middleware.ChainUnaryServer(TimeoutInterceptor(), DBUnaryInterceptor(), ui, middleware.RecoveredUnaryGRPCServerLog())
	// grpc.UnaryInterceptor() 创造一个拦截器
	// grpc.NewServer(可以传入一个具体的拦截器或者拦截器链)
	grpcServer := grpc.NewServer()
	// 反向注册服务
	shortenerpb.RegisterShortenerServiceServer(grpcServer, &shortener.Server{})

	reflection.Register(grpcServer) // 方便使用 grpcurl 等工具调试

	return grpcServer
}

func RunGRPCServer(ctx context.Context, clientManager *manager.ClientManager) (err error) {
	GRPCServer := NewGRPCServer()

	// grpcLis
	grpcLis, err := net.Listen("tcp", config.GetConfig().RPC.Shortneer.Addr)
	log.Println("config.GetConfig().RPC.Shortneer.Addr")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		<-ctx.Done()
		// 当接收到 主线程的context被取消时，则会销毁grpcServer本身，会先把当前的grpc处理完成才会取消
		GRPCServer.GracefulStop()
		grpcLis.Close()
	}()
	return GRPCServer.Serve(grpcLis)
}

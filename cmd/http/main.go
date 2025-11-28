package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/username/shorturl/internal/config"
	internal "github.com/username/shorturl/internal/handler"
	"github.com/username/shorturl/internal/manager"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 加载配置
	config.LoadAll()
	log.Println("Configuration loaded")
	// 阻塞启动块

	// 异步启动块
	ctx, cancel := context.WithCancel(context.Background())

	g, gCtx := errgroup.WithContext(ctx)

	cliManager := manager.NewClientManer()

	grpcServers := config.GetConfig().GRPCServers
	log.Println(grpcServers.Shortener, grpcServers.Clipboarder)
	var grpcAddrs []string = []string{grpcServers.Shortener.Addr, grpcServers.Clipboarder.Addr}

	// 1. 启动客户端管理器初始化
	g.Go(func() error {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		log.Println("Initializing services...")
		err := cliManager.InitServices(gCtx, grpcAddrs, opts...)
		return err
	})
	// 启动 HTTP
	g.Go(func() error {
		return internal.RunHTTPServer(gCtx, cliManager)
	})

	// 监听关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 等待信号或者服务错误
	select {
	case <-sigChan:
		log.Println("Shutdown signal received")
	case <-gCtx.Done():
		log.Println("Service error, shutting down")
	}

	log.Println("开始执行cancel")
	cancel()
	log.Println("执行cancel结束")
	if err := g.Wait(); err != nil {
		log.Printf("Service error: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/username/shorturl/internal/config"
	"github.com/username/shorturl/internal/manager"
	"github.com/username/shorturl/internal/rpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.LoadAll()
	log.Println("Configuration loaded")

	ctx, cancel := context.WithCancel(context.Background())
	// 延迟取消操作保留，作为主函数退出的二次保障，但不是退出流程的触发器。

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
	// 2. 启动gRPC Server
	g.Go(func() error {
		// 确保 rpc.RunGRPCServer(ctx, ...) 内部监听 ctx.Done()，并在 context 取消时优雅停止
		return rpc.RunGRPCServer(ctx, cliManager)
	})

	// 3. 监听关闭信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 等待信号或者服务错误
	select {
	case <-sigChan:
		log.Println("Shutdown signal received, initiating graceful shutdown...")
	case <-gCtx.Done():
		log.Println("Service error detected, shutting down...")
	}

	// 修正 1：显式调用 cancel()，触发所有 Goroutine 退出
	cancel()

	// 关闭客户端管理器
	if err := cliManager.Close(); err != nil {
		log.Printf("Warning: Client Manager Close failed: %v", err)
	}

	// 等待所有 goroutine 退出
	if err := g.Wait(); err != nil && err != context.Canceled {
		log.Printf("Service shut down with error: %v", err)
	} else {
		log.Println("All services shut down gracefully.")
	}
}

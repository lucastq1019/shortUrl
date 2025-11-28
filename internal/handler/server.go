package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/username/shorturl/internal/config"
	"github.com/username/shorturl/internal/manager"
	shortenerpb "github.com/username/shorturl/internal/rpc/proto"
)

func RunHTTPServer(ctx context.Context, clientManager *manager.ClientManager) (err error) {
	errCh := make(chan error, 1)

	httpServer := &http.Server{
		Addr:    config.GetConfig().HttpAddr,
		Handler: NewRouter(clientManager),
	}
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("http server start close")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}

// NewRouter 创建路由
func NewRouter(clientManager *manager.ClientManager) *gin.Engine {
	gin.SetMode(config.GetConfig().GinMode)

	router := gin.New()

	// 默认中间件 日志，恢复
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "hello shortURL service")
	})

	router.GET("/testRpc", func(ctx *gin.Context) {
		// 1. 创建请求结构体
		req := &shortenerpb.CreateShortLinkRequest{
			LongUrl: "https://www.google.com",
		}

		// 2. 发起 gRPC 调用
		// 替换 `client.MyMethod` 为你实际的 gRPC 客户端方法

		client, err := clientManager.GetClient(config.GetConfig().GRPCServers.Shortener.Addr)

		// conn, err := grpc.NewClient(
		// 	"localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()), // 👈 加上这行
		// )
		if err != nil {
			log.Printf("%v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "RPC create conn failed",
				"details": err.Error(),
			})
			return
		}
		tempre, ok := client.(shortenerpb.ShortenerServiceClient)
		log.Println(tempre, ok)
		if !ok {
			log.Panicln("错误转换")
		}
		resp, err := tempre.CreateShortLink(ctx, req)

		// 3. 处理响应
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "RPC call failed",
				"details": err.Error(),
			})
			return
		}

		// 4. 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"message": "RPC call successful",
			"data":    resp, // 假设 resp 可以被 JSON 序列化
		})
	})

	router.GET("/t", func(ctx *gin.Context) {

		// 1. 创建请求结构体
		req := &shortenerpb.GetLongURLRequest{
			ShortKey: "https://www.google.com",
		}

		client, err := clientManager.GetClient(config.GetConfig().GRPCServers.Shortener.Addr)

		if err != nil {
			log.Printf("%s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "RPC create conn failed",
				"details": err.Error(),
			})
			return
		}
		tempre, ok := client.(shortenerpb.ShortenerServiceClient)
		if !ok {
			log.Panicln("错误转换")
		}
		resp, err := tempre.GetLongURL(ctx, req)

		// 3. 处理响应
		if err != nil {
			log.Printf("%s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "RPC call failed",
				"details": err.Error(),
			})
			return
		}

		// 4. 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"message": "RPC call successful",
			"data":    resp, // 假设 resp 可以被 JSON 序列化
		})
	})

	return router
}

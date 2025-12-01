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

// 只负责 协调，而每个功能的 HTTP 适配细节 都隔离在自己的文件中，极大地提高了代码的可读性和可维护性。
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

// RouterHandlers 聚合所有资源（客户端依赖）
// 我们在它上面定义方法作为 HTTP Handlers
type RouterHandlers struct {
	Shortener shortenerpb.ShortenerServiceClient
	Clipboard shortenerpb.ClipboarderServiceClient
}

// NewRouter 创建路由，并注册所有 Handler
func NewRouter(clientManager *manager.ClientManager) *gin.Engine {
	log.Println(clientManager)
	ShortenerServiceClient, err := clientManager.GetClient(config.GetConfig().GRPCServers.Shortener.Addr)
	if err != nil {
		log.Fatalf("获取链接失败,%v", err)
	}
	clientShort, ok := ShortenerServiceClient.(shortenerpb.ShortenerServiceClient)
	if !ok {
		log.Fatalf("无法创建ShortenerServiceClient,%v", err)
	}

	// 1. 初始化 RouterHandlers（注入依赖，只做一次）
	rh := &RouterHandlers{
		Shortener: clientShort,
		// Clipboard: clientManager,
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "API Gateway is running")
	})

	// 2. 按功能资源创建路由分组
	// --- Shortener 路由 ---
	shortenerGroup := router.Group("/shortener/v1")
	// 调用外部文件中的注册函数
	rh.RegisterShortenerRoutes(shortenerGroup)

	// --- Clipboard 路由 ---
	// clipboardGroup := router.Group("/clipboard/v1")
	// 调用外部文件中的注册函数
	// rh.RegisterClipboardRoutes(clipboardGroup)

	return router
}

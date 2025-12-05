package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	shortenerpb "github.com/username/shorturl/internal/rpc/proto"
)

func (rh *RouterHandlers) RegisterShortenerRoutes(group *gin.RouterGroup) {
	group.POST("/c", rh.HandleCreateShortLink)
	group.GET("/:key", rh.HandleGetLongURL)
	group.GET("/all", rh.HandleGetAllShortLink)
}

// HandleCreateShortLink 是 Shortener 资源的 HTTP Handler
func (rh *RouterHandlers) HandleCreateShortLink(ctx *gin.Context) {
	var reqBody struct {
		LongURL string `json:"long_url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 调用 gRPC 客户端封装层（核心：转发请求）
	resp, err := rh.Shortener.CreateShortLink(ctx, &shortenerpb.CreateShortLinkRequest{
		LongUrl: reqBody.LongURL,
	})

	if err != nil {
		log.Printf("Shortener RPC failed: %v", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Backend service unavailable"})
		return
	}

	// 格式化并返回 HTTP 响应
	ctx.JSON(http.StatusOK, gin.H{"short_url": resp.ShortKey})
}

func (rh *RouterHandlers) HandleGetLongURL(ctx *gin.Context) {
	var reqBody struct {
		ShortKey string `json:"short_key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 调用 gRPC 客户端封装层
	resp, err := rh.Shortener.GetLongURL(ctx, &shortenerpb.GetLongURLRequest{
		ShortKey: reqBody.ShortKey,
	})

	if err != nil {
		// ... 错误处理
		return
	}

	// ... 返回响应
	ctx.JSON(http.StatusOK, gin.H{"long_url": resp.GetLongUrl()})
}

func (rh *RouterHandlers) HandleGetAllShortLink(ctx *gin.Context) {

	// 调用 gRPC 客户端封装层
	resp, err := rh.Shortener.GetAllShortLink(ctx, &shortenerpb.GetAllShortLinkRequest{})

	if err != nil {
		log.Printf("调用获取全部链接错误,%v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"short_links": resp.ShortLinks})
}

package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	"goweb/services/file-api/internal/handler"
)

// NewRouter 创建路由
func NewRouter(cfg *config.Config, log logger.Logger, fileHandler *handler.FileHandler, chunkHandler *handler.ChunkHandler) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLogger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "file-api"})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		files := api.Group("/files")
		{
			// 上传文件
			files.POST("/upload", fileHandler.UploadFile)

			// 文件列表
			files.GET("", fileHandler.ListFiles)

			// 文件详情
			files.GET("/:id", fileHandler.GetFile)

			// 下载文件
			files.GET("/:id/download", fileHandler.DownloadFile)

			// 删除文件
			files.DELETE("/:id", fileHandler.DeleteFile)

			// 根据哈希获取文件
			files.GET("/hash/:hash", fileHandler.GetFileByHash)

			// 统计信息
			files.GET("/statistics", fileHandler.GetStatistics)

			// 分片上传相关接口
			chunks := files.Group("/chunks")
			{
				// 初始化分片上传
				chunks.POST("/init", chunkHandler.InitiateChunkUpload)

				// 上传分片
				chunks.POST("/upload", chunkHandler.UploadChunk)

				// 合并分片
				chunks.POST("/merge", chunkHandler.MergeChunks)

				// 获取分片上传状态
				chunks.GET("/status/:upload_id", chunkHandler.GetChunkUploadStatus)

				// 列出分片
				chunks.GET("/list/:upload_id", chunkHandler.ListChunks)

				// 中止分片上传
				chunks.DELETE("/abort/:upload_id", chunkHandler.AbortChunkUpload)
			}
		}
	}

	// 静态文件服务（可选，用于直接访问上传的文件）
	// 注意：生产环境建议使用专门的静态文件服务器或对象存储
	storageBasePath := cfg.GetString("storage.local.base_path")
	if storageBasePath != "" {
		r.Static("/uploads", storageBasePath)
	}

	return r
}

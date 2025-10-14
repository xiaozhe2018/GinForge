package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"goweb/pkg/base"
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/file-api/internal/service"
)

// ChunkHandler 分片上传处理器
type ChunkHandler struct {
	*base.BaseHandler
	fileService *service.FileService
	logger      logger.Logger
	config      *config.Config
}

// NewChunkHandler 创建分片上传处理器
func NewChunkHandler(fileService *service.FileService, cfg *config.Config, log logger.Logger) *ChunkHandler {
	return &ChunkHandler{
		BaseHandler: base.NewBaseHandler(log),
		fileService: fileService,
		logger:      log,
		config:      cfg,
	}
}

// InitiateChunkUpload 初始化分片上传
// @Summary 初始化分片上传
// @Description 初始化分片上传，获取上传ID
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file_name formData string true "文件名"
// @Param total_chunks formData int true "总分片数"
// @Param total_size formData int true "文件总大小"
// @Param chunk_size formData int true "分片大小"
// @Param file_hash formData string false "文件哈希"
// @Param sub_path formData string false "子路径"
// @Param user_id formData int false "用户ID"
// @Param user_type formData string false "用户类型"
// @Success 200 {object} response.Response{data=string}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/init [post]
func (h *ChunkHandler) InitiateChunkUpload(c *gin.Context) {
	// 1. 解析参数
	fileName := c.PostForm("file_name")
	if fileName == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	totalChunksStr := c.PostForm("total_chunks")
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil || totalChunks <= 0 {
		response.BadRequest(c, "无效的分片总数")
		return
	}

	totalSizeStr := c.PostForm("total_size")
	totalSize, err := strconv.ParseInt(totalSizeStr, 10, 64)
	if err != nil || totalSize <= 0 {
		response.BadRequest(c, "无效的文件总大小")
		return
	}

	chunkSizeStr := c.PostForm("chunk_size")
	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil || chunkSize <= 0 {
		response.BadRequest(c, "无效的分片大小")
		return
	}

	fileHash := c.PostForm("file_hash")
	subPath := c.PostForm("sub_path")
	userIDStr := c.PostForm("user_id")
	userType := c.PostForm("user_type")

	var userID uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的用户ID")
			return
		}
		userID = uint(id)
	}

	// 2. 构建请求
	req := &service.ChunkUploadRequest{
		FileName:    fileName,
		TotalChunks: totalChunks,
		TotalSize:   totalSize,
		ChunkSize:   chunkSize,
		FileHash:    fileHash,
		SubPath:     subPath,
		UserID:      userID,
		UserType:    userType,
		IP:          c.ClientIP(),
	}

	// 3. 调用服务初始化分片上传
	uploadID, err := h.fileService.InitiateChunkUpload(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to initiate chunk upload", "error", err)
		response.InternalError(c, fmt.Sprintf("初始化分片上传失败: %v", err))
		return
	}

	// 4. 返回上传ID
	response.Success(c, uploadID)
}

// UploadChunk 上传分片
// @Summary 上传分片
// @Description 上传单个分片
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param upload_id formData string true "上传ID"
// @Param chunk_index formData int true "分片索引"
// @Param total_chunks formData int true "总分片数"
// @Param file formData file true "分片文件"
// @Param file_name formData string true "文件名"
// @Param total_size formData int true "文件总大小"
// @Param chunk_size formData int true "分片大小"
// @Param file_hash formData string false "文件哈希"
// @Param sub_path formData string false "子路径"
// @Param user_id formData int false "用户ID"
// @Param user_type formData string false "用户类型"
// @Success 200 {object} response.Response{data=service.ChunkUploadResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/upload [post]
func (h *ChunkHandler) UploadChunk(c *gin.Context) {
	// 1. 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file", "error", err)
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 2. 解析参数
	uploadID := c.PostForm("upload_id")
	if uploadID == "" {
		response.BadRequest(c, "上传ID不能为空")
		return
	}

	chunkIndexStr := c.PostForm("chunk_index")
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || chunkIndex < 0 {
		response.BadRequest(c, "无效的分片索引")
		return
	}

	totalChunksStr := c.PostForm("total_chunks")
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil || totalChunks <= 0 {
		response.BadRequest(c, "无效的分片总数")
		return
	}

	fileName := c.PostForm("file_name")
	if fileName == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	totalSizeStr := c.PostForm("total_size")
	totalSize, err := strconv.ParseInt(totalSizeStr, 10, 64)
	if err != nil || totalSize <= 0 {
		response.BadRequest(c, "无效的文件总大小")
		return
	}

	chunkSizeStr := c.PostForm("chunk_size")
	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil || chunkSize <= 0 {
		response.BadRequest(c, "无效的分片大小")
		return
	}

	fileHash := c.PostForm("file_hash")
	subPath := c.PostForm("sub_path")
	userIDStr := c.PostForm("user_id")
	userType := c.PostForm("user_type")

	var userID uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的用户ID")
			return
		}
		userID = uint(id)
	}

	// 3. 构建请求
	req := &service.ChunkUploadRequest{
		UploadID:    uploadID,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
		File:        file,
		FileName:    fileName,
		TotalSize:   totalSize,
		ChunkSize:   chunkSize,
		FileHash:    fileHash,
		SubPath:     subPath,
		UserID:      userID,
		UserType:    userType,
		IP:          c.ClientIP(),
	}

	// 4. 调用服务上传分片
	result, err := h.fileService.UploadChunk(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to upload chunk", "error", err, "upload_id", uploadID, "chunk", chunkIndex)
		response.InternalError(c, fmt.Sprintf("上传分片失败: %v", err))
		return
	}

	// 5. 返回结果
	response.Success(c, result)
}

// MergeChunks 合并分片
// @Summary 合并分片
// @Description 合并已上传的所有分片
// @Tags 文件管理
// @Accept application/json
// @Produce json
// @Param request body service.ChunkMergeRequest true "合并请求"
// @Success 200 {object} response.Response{data=service.UploadResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/merge [post]
func (h *ChunkHandler) MergeChunks(c *gin.Context) {
	// 1. 解析参数
	var req service.ChunkMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, fmt.Sprintf("无效的请求参数: %v", err))
		return
	}

	// 2. 设置IP
	req.IP = c.ClientIP()

	// 3. 调用服务合并分片
	result, err := h.fileService.MergeChunks(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to merge chunks", "error", err, "upload_id", req.UploadID)
		response.InternalError(c, fmt.Sprintf("合并分片失败: %v", err))
		return
	}

	// 4. 返回结果
	response.Success(c, result)
}

// GetChunkUploadStatus 获取分片上传状态
// @Summary 获取分片上传状态
// @Description 获取分片上传的当前状态
// @Tags 文件管理
// @Produce json
// @Param upload_id path string true "上传ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/status/{upload_id} [get]
func (h *ChunkHandler) GetChunkUploadStatus(c *gin.Context) {
	// 1. 获取上传ID
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		response.BadRequest(c, "上传ID不能为空")
		return
	}

	// 2. 调用服务获取状态
	status, err := h.fileService.GetChunkUploadStatus(c.Request.Context(), uploadID)
	if err != nil {
		h.logger.Error("Failed to get chunk upload status", "error", err, "upload_id", uploadID)
		response.NotFound(c, "未找到上传记录或已过期")
		return
	}

	// 3. 返回结果
	response.Success(c, status)
}

// ListChunks 列出分片
// @Summary 列出已上传的分片
// @Description 获取已上传分片的列表
// @Tags 文件管理
// @Produce json
// @Param upload_id path string true "上传ID"
// @Success 200 {object} response.Response{data=[]int}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/list/{upload_id} [get]
func (h *ChunkHandler) ListChunks(c *gin.Context) {
	// 1. 获取上传ID
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		response.BadRequest(c, "上传ID不能为空")
		return
	}

	// 2. 调用服务获取分片列表
	chunks, err := h.fileService.ListChunks(c.Request.Context(), uploadID)
	if err != nil {
		h.logger.Error("Failed to list chunks", "error", err, "upload_id", uploadID)
		response.NotFound(c, "未找到上传记录或已过期")
		return
	}

	// 3. 返回结果
	response.Success(c, chunks)
}

// AbortChunkUpload 中止分片上传
// @Summary 中止分片上传
// @Description 中止分片上传并清理资源
// @Tags 文件管理
// @Produce json
// @Param upload_id path string true "上传ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/chunks/abort/{upload_id} [delete]
func (h *ChunkHandler) AbortChunkUpload(c *gin.Context) {
	// 1. 获取上传ID
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		response.BadRequest(c, "上传ID不能为空")
		return
	}

	// 2. 调用服务中止上传
	if err := h.fileService.AbortChunkUpload(c.Request.Context(), uploadID); err != nil {
		h.logger.Error("Failed to abort chunk upload", "error", err, "upload_id", uploadID)
		response.NotFound(c, "未找到上传记录或已过期")
		return
	}

	// 3. 返回结果
	response.Success(c, nil)
}
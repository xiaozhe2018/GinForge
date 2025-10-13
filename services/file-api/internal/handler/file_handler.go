package handler

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"goweb/pkg/base"
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/file-api/internal/service"
)

// FileHandler 文件处理器
type FileHandler struct {
	*base.BaseHandler
	fileService *service.FileService
	logger      logger.Logger
	config      *config.Config
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileService *service.FileService, cfg *config.Config, log logger.Logger) *FileHandler {
	return &FileHandler{
		BaseHandler: base.NewBaseHandler(log),
		fileService: fileService,
		logger:      log,
		config:      cfg,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param description formData string false "文件描述"
// @Param tags formData string false "文件标签(逗号分隔)"
// @Param sub_path formData string false "子路径"
// @Param user_id formData int false "用户ID"
// @Param user_type formData string false "用户类型(admin/user/merchant)"
// @Success 200 {object} response.Response{data=service.UploadResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 1. 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file", err)
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 2. 获取其他参数
	description := c.PostForm("description")
	tags := c.PostForm("tags")
	subPath := c.PostForm("sub_path")
	userIDStr := c.PostForm("user_id")
	userType := c.PostForm("user_type")

	var userID uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			h.logger.Error("failed to parse user_id", err)
			response.BadRequest(c, "无效的用户ID")
			return
		}
		userID = uint(id)
	}

	// 3. 构建上传请求
	req := &service.UploadRequest{
		File:        file,
		Description: description,
		Tags:        tags,
		SubPath:     subPath,
		UserID:      userID,
		UserType:    userType,
		IP:          c.ClientIP(),
	}

	// 4. 调用服务上传文件
	result, err := h.fileService.UploadFile(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to upload file", err)
		response.InternalError(c, fmt.Sprintf("上传文件失败: %v", err))
		return
	}

	// 5. 返回结果
	response.Success(c, result)
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 通过文件ID下载文件
// @Tags 文件管理
// @Produce octet-stream
// @Param id path int true "文件ID"
// @Param user_id query int false "用户ID"
// @Success 200 {file} file "文件内容"
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/{id}/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	// 1. 解析请求参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("failed to parse file id", err)
		response.BadRequest(c, "无效的文件ID")
		return
	}

	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		uid, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			h.logger.Error("failed to parse user_id", err)
			response.BadRequest(c, "无效的用户ID")
			return
		}
		userID = uint(uid)
	}

	// 2. 构建下载请求
	req := &service.DownloadRequest{
		FileID: uint(id),
		UserID: userID,
		IP:     c.ClientIP(),
	}

	// 3. 调用服务获取文件
	fileRecord, err := h.fileService.DownloadFile(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to get file for download", err)
		response.NotFound(c, "文件不存在或已被删除")
		return
	}

	// 4. 构建完整文件路径
	basePath := h.config.GetString("storage.local.base_path")
	fullPath := filepath.Join(basePath, fileRecord.RelativePath)

	// 5. 发送文件
	c.FileAttachment(fullPath, fileRecord.OriginalName)
}

// ListFiles 列出文件
// @Summary 列出文件
// @Description 分页获取文件列表
// @Tags 文件管理
// @Produce json
// @Param user_id query int false "用户ID"
// @Param file_type query string false "文件类型(image/video/document/other)"
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} response.Response{data=service.ListFilesResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	// 1. 解析请求参数
	userIDStr := c.Query("user_id")
	fileType := c.Query("file_type")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest(c, "无效的页码")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		response.BadRequest(c, "无效的每页数量")
		return
	}

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
	req := &service.ListFilesRequest{
		UserID:   userID,
		FileType: fileType,
		Page:     page,
		PageSize: pageSize,
	}

	// 3. 调用服务获取文件列表
	result, err := h.fileService.ListFiles(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to list files", err)
		response.InternalError(c, "获取文件列表失败")
		return
	}

	// 4. 返回结果
	response.Success(c, result)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 通过文件ID删除文件
// @Tags 文件管理
// @Produce json
// @Param id path int true "文件ID"
// @Param user_id query int false "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	// 1. 解析请求参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文件ID")
		return
	}

	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		uid, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的用户ID")
			return
		}
		userID = uint(uid)
	}

	// 2. 构建请求
	req := &service.DeleteFileRequest{
		FileID: uint(id),
		UserID: userID,
	}

	// 3. 调用服务删除文件
	if err := h.fileService.DeleteFile(c.Request.Context(), req); err != nil {
		h.logger.Error("failed to delete file", err)
		response.InternalError(c, "删除文件失败")
		return
	}

	// 4. 返回结果
	response.Success(c, nil)
}

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Description 通过文件ID获取文件详细信息
// @Tags 文件管理
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response{data=model.FileRecord}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	// 1. 解析请求参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文件ID")
		return
	}

	// 2. 构建请求
	req := &service.GetFileRequest{
		FileID: uint(id),
	}

	// 3. 调用服务获取文件信息
	fileRecord, err := h.fileService.GetFile(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to get file", err)
		response.NotFound(c, "文件不存在或已被删除")
		return
	}

	// 4. 返回结果
	response.Success(c, fileRecord)
}

// GetFileByHash 根据哈希获取文件
// @Summary 根据哈希获取文件
// @Description 通过文件哈希获取文件信息
// @Tags 文件管理
// @Produce json
// @Param hash path string true "文件哈希"
// @Success 200 {object} response.Response{data=model.FileRecord}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/files/hash/{hash} [get]
func (h *FileHandler) GetFileByHash(c *gin.Context) {
	// 1. 获取哈希参数
	hash := c.Param("hash")
	if hash == "" {
		response.BadRequest(c, "文件哈希不能为空")
		return
	}

	// 2. 构建请求
	req := &service.GetFileByHashRequest{
		Hash: hash,
	}

	// 3. 调用服务获取文件信息
	fileRecord, err := h.fileService.GetFileByHash(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to get file by hash", err)
		response.NotFound(c, "文件不存在或已被删除")
		return
	}

	// 4. 返回结果
	response.Success(c, fileRecord)
}

// GetStatistics 获取统计信息
// @Summary 获取文件统计信息
// @Description 获取文件总数、总大小、类型分布等统计信息
// @Tags 文件管理
// @Produce json
// @Success 200 {object} response.Response{data=service.GetStatisticsResponse}
// @Failure 500 {object} response.Response
// @Router /api/v1/files/statistics [get]
func (h *FileHandler) GetStatistics(c *gin.Context) {
	// 调用服务获取统计信息
	statistics, err := h.fileService.GetStatistics(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get statistics", err)
		response.InternalError(c, "获取统计信息失败")
		return
	}

	// 返回结果
	response.Success(c, statistics)
}

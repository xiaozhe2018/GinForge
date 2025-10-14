// Package chunked 提供分片上传实现
package chunked

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"goweb/pkg/logger"
	"goweb/pkg/storage"
	"goweb/pkg/upload"
)

// 错误定义
var (
	ErrChunkNotFound      = errors.New("chunk not found")
	ErrChunkUploadExpired = errors.New("chunk upload expired")
	ErrChunkInvalidIndex  = errors.New("invalid chunk index")
	ErrChunkInvalidTotal  = errors.New("invalid total chunks")
	ErrChunkAlreadyExists = errors.New("chunk already exists")
	ErrChunkMergeInProgress = errors.New("chunk merge in progress")
)

// ChunkUploadStatus 分片上传状态
type ChunkUploadStatus struct {
	UploadID    string
	FileName    string
	TotalChunks int
	TotalSize   int64
	ChunkSize   int64
	SubPath     string
	UserID      uint
	UserType    string
	IP          string
	StartTime   time.Time
	LastUpdate  time.Time
	Chunks      map[int]bool  // 已上传的分片索引
	MergeInProgress bool      // 是否正在合并
	Metadata    map[string]interface{} // 元数据
}

// Config 分片上传配置
type Config struct {
	Enabled      bool   // 是否启用
	ChunkSize    int64  // 分片大小
	ChunkTimeout int    // 分片超时时间（秒）
	ChunkDir     string // 分片存储目录
	TempDir      string // 临时目录
}

// ChunkUploader 分片上传实现
type ChunkUploader struct {
	*upload.BaseUploader
	config     Config
	logger     logger.Logger
	storage    storage.Storage
	chunkLock  sync.Mutex
	chunkUploads map[string]*ChunkUploadStatus
}

// New 创建分片上传实现
func New(config Config, storage storage.Storage, log logger.Logger) *ChunkUploader {
	// 使用默认值
	if config.ChunkSize <= 0 {
		config.ChunkSize = 5 * 1024 * 1024 // 5MB
	}
	if config.ChunkTimeout <= 0 {
		config.ChunkTimeout = 3600 // 1小时
	}
	if config.ChunkDir == "" {
		config.ChunkDir = "./uploads/chunks"
	}
	if config.TempDir == "" {
		config.TempDir = "./uploads/temp"
	}

	// 确保目录存在
	os.MkdirAll(config.ChunkDir, 0755)
	os.MkdirAll(config.TempDir, 0755)

	return &ChunkUploader{
		BaseUploader: upload.NewBaseUploader(upload.UploadTypeChunked, storage, log),
		config:       config,
		logger:       log,
		storage:      storage,
		chunkUploads: make(map[string]*ChunkUploadStatus),
	}
}

// Upload 上传文件
func (u *ChunkUploader) Upload(ctx context.Context, req *upload.UploadRequest) (*upload.UploadResponse, error) {
	// 分片上传不支持直接上传，应该使用分片上传接口
	return nil, errors.New("chunked uploader does not support direct upload, use chunked upload API instead")
}

// CanHandle 是否可以处理
func (u *ChunkUploader) CanHandle(req *upload.UploadRequest) bool {
	// 分片上传不支持直接上传
	return false
}

// InitiateUpload 初始化分片上传
func (u *ChunkUploader) InitiateUpload(ctx context.Context, req *upload.ChunkUploadRequest) (string, error) {
	// 检查分片上传是否启用
	if !u.config.Enabled {
		return "", errors.New("chunked upload is disabled")
	}

	// 生成上传ID
	uploadID := uuid.New().String()

	// 创建分片上传状态
	u.chunkLock.Lock()
	defer u.chunkLock.Unlock()

	u.chunkUploads[uploadID] = &ChunkUploadStatus{
		UploadID:    uploadID,
		FileName:    req.FileName,
		TotalChunks: req.TotalChunks,
		TotalSize:   req.TotalSize,
		ChunkSize:   req.ChunkSize,
		SubPath:     req.SubPath,
		UserID:      req.UserID,
		UserType:    req.UserType,
		IP:          req.IP,
		StartTime:   time.Now(),
		LastUpdate:  time.Now(),
		Chunks:      make(map[int]bool),
		Metadata:    req.Metadata,
	}

	// 确保分片目录存在
	chunkDir := u.getChunkDir(uploadID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		u.logger.Error("Failed to create chunk directory", "error", err, "dir", chunkDir)
		delete(u.chunkUploads, uploadID)
		return "", fmt.Errorf("failed to create chunk directory: %w", err)
	}

	u.logger.Info("Initiated chunked upload", "upload_id", uploadID, "file", req.FileName, "chunks", req.TotalChunks)
	return uploadID, nil
}

// UploadChunk 上传分片
func (u *ChunkUploader) UploadChunk(ctx context.Context, req *upload.ChunkUploadRequest) (*upload.ChunkUploadResponse, error) {
	// 检查分片上传是否启用
	if !u.config.Enabled {
		return nil, errors.New("chunked upload is disabled")
	}

	// 如果没有上传ID，则初始化上传
	if req.UploadID == "" {
		uploadID, err := u.InitiateUpload(ctx, req)
		if err != nil {
			return nil, err
		}
		req.UploadID = uploadID
	}

	// 获取上传状态
	u.chunkLock.Lock()
	uploadStatus, exists := u.chunkUploads[req.UploadID]
	if !exists {
		u.chunkLock.Unlock()
		return nil, ErrChunkNotFound
	}

	// 检查是否已过期
	if time.Since(uploadStatus.LastUpdate) > time.Duration(u.config.ChunkTimeout)*time.Second {
		u.cleanupExpiredUpload(req.UploadID)
		u.chunkLock.Unlock()
		return nil, ErrChunkUploadExpired
	}

	// 检查分片索引是否有效
	if req.ChunkIndex < 0 || req.ChunkIndex >= req.TotalChunks {
		u.chunkLock.Unlock()
		return nil, ErrChunkInvalidIndex
	}

	// 检查分片是否已上传
	if uploadStatus.Chunks[req.ChunkIndex] {
		u.chunkLock.Unlock()
		return nil, ErrChunkAlreadyExists
	}

	// 检查是否正在合并
	if uploadStatus.MergeInProgress {
		u.chunkLock.Unlock()
		return nil, ErrChunkMergeInProgress
	}

	// 更新上传状态
	uploadStatus.LastUpdate = time.Now()
	u.chunkLock.Unlock()

	// 打开源文件
	src, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 保存分片
	chunkPath := u.getChunkPath(req.UploadID, req.ChunkIndex)
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		u.logger.Error("Failed to save chunk", "error", err, "upload_id", req.UploadID, "chunk", req.ChunkIndex)
		return nil, fmt.Errorf("failed to save chunk: %w", err)
	}

	// 标记分片已上传
	u.chunkLock.Lock()
	uploadStatus.Chunks[req.ChunkIndex] = true
	allUploaded := len(uploadStatus.Chunks) == uploadStatus.TotalChunks
	u.chunkLock.Unlock()

	u.logger.Info("Uploaded chunk", "upload_id", req.UploadID, "chunk", req.ChunkIndex, "total", req.TotalChunks, "completed", allUploaded)

	return &upload.ChunkUploadResponse{
		UploadID:   req.UploadID,
		ChunkIndex: req.ChunkIndex,
		TotalChunks: req.TotalChunks,
		Completed:  allUploaded,
	}, nil
}

// MergeChunks 合并分片
func (u *ChunkUploader) MergeChunks(ctx context.Context, req *upload.ChunkMergeRequest) (*upload.UploadResponse, error) {
	// 检查分片上传是否启用
	if !u.config.Enabled {
		return nil, errors.New("chunked upload is disabled")
	}

	// 获取上传状态
	u.chunkLock.Lock()
	uploadStatus, exists := u.chunkUploads[req.UploadID]
	if !exists {
		u.chunkLock.Unlock()
		return nil, ErrChunkNotFound
	}

	// 检查是否已过期
	if time.Since(uploadStatus.LastUpdate) > time.Duration(u.config.ChunkTimeout)*time.Second {
		u.cleanupExpiredUpload(req.UploadID)
		u.chunkLock.Unlock()
		return nil, ErrChunkUploadExpired
	}

	// 检查是否所有分片都已上传
	if len(uploadStatus.Chunks) != uploadStatus.TotalChunks {
		u.chunkLock.Unlock()
		return nil, fmt.Errorf("not all chunks uploaded: %d/%d", len(uploadStatus.Chunks), uploadStatus.TotalChunks)
	}

	// 检查是否正在合并
	if uploadStatus.MergeInProgress {
		u.chunkLock.Unlock()
		return nil, ErrChunkMergeInProgress
	}

	// 标记为正在合并
	uploadStatus.MergeInProgress = true
	u.chunkLock.Unlock()

	// 创建临时文件
	tempFilePath := u.getTempFilePath(req.UploadID)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		u.logger.Error("Failed to create temp file", "error", err, "path", tempFilePath)
		u.cleanupMergeStatus(req.UploadID)
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// 获取所有分片路径
	chunkPaths := make([]string, uploadStatus.TotalChunks)
	for i := 0; i < uploadStatus.TotalChunks; i++ {
		chunkPaths[i] = u.getChunkPath(req.UploadID, i)
	}

	// 按顺序合并分片
	hash := md5.New()
	totalSize := int64(0)
	for i, chunkPath := range chunkPaths {
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			u.logger.Error("Failed to open chunk file", "error", err, "chunk", i, "path", chunkPath)
			u.cleanupMergeStatus(req.UploadID)
			return nil, fmt.Errorf("failed to open chunk file: %w", err)
		}

		// 复制分片内容到临时文件
		written, err := io.Copy(tempFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			u.logger.Error("Failed to copy chunk data", "error", err, "chunk", i)
			u.cleanupMergeStatus(req.UploadID)
			return nil, fmt.Errorf("failed to copy chunk data: %w", err)
		}

		// 更新哈希和大小
		chunkFile, err = os.Open(chunkPath)
		if err == nil {
			io.Copy(hash, chunkFile)
			chunkFile.Close()
		}

		totalSize += written
	}

	// 关闭临时文件
	tempFile.Close()

	// 计算文件哈希
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// 检查文件大小是否匹配
	if totalSize != uploadStatus.TotalSize {
		u.logger.Warn("File size mismatch", "expected", uploadStatus.TotalSize, "actual", totalSize)
		// 继续处理，但记录警告
	}

	// 检查文件哈希是否匹配（如果提供）
	if req.FileHash != "" && req.FileHash != fileHash {
		u.logger.Warn("File hash mismatch", "expected", req.FileHash, "actual", fileHash)
		// 继续处理，但记录警告
	}

	// 读取临时文件内容
	fileData, err := os.ReadFile(tempFilePath)
	if err != nil {
		u.logger.Error("Failed to read temp file", "error", err)
		u.cleanupMergeStatus(req.UploadID)
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	// 使用存储服务保存文件
	fileInfo, err := u.storage.SaveFileWithContext(ctx, fileData, req.FileName, req.SubPath)
	if err != nil {
		u.logger.Error("Failed to save merged file", "error", err)
		u.cleanupMergeStatus(req.UploadID)
		return nil, fmt.Errorf("failed to save merged file: %w", err)
	}

	// 添加元数据
	if uploadStatus.Metadata != nil {
		for k, v := range uploadStatus.Metadata {
			fileInfo.Metadata[k] = v
		}
	}
	fileInfo.Metadata["user_id"] = uploadStatus.UserID
	fileInfo.Metadata["user_type"] = uploadStatus.UserType
	fileInfo.Metadata["ip"] = uploadStatus.IP
	fileInfo.Metadata["upload_method"] = "chunked"

	// 处理文件
	fileInfo, err = u.ProcessFile(ctx, fileInfo)
	if err != nil {
		u.logger.Error("Failed to process file", "error", err)
		u.cleanupMergeStatus(req.UploadID)
		return nil, fmt.Errorf("failed to process file: %w", err)
	}

	// 清理分片文件
	u.cleanupUpload(req.UploadID)

	// 构建响应
	response := upload.ConvertToResponse(fileInfo, 0) // ID将由存储层填充

	u.logger.Info("Merged chunks successfully", "upload_id", req.UploadID, "file", req.FileName, "size", totalSize)
	return response, nil
}

// GetUploadStatus 获取上传状态
func (u *ChunkUploader) GetUploadStatus(ctx context.Context, uploadID string) (interface{}, error) {
	u.chunkLock.Lock()
	defer u.chunkLock.Unlock()

	status, exists := u.chunkUploads[uploadID]
	if !exists {
		return nil, ErrChunkNotFound
	}

	// 创建副本以避免并发修改
	result := &ChunkUploadStatus{
		UploadID:    status.UploadID,
		FileName:    status.FileName,
		TotalChunks: status.TotalChunks,
		TotalSize:   status.TotalSize,
		ChunkSize:   status.ChunkSize,
		SubPath:     status.SubPath,
		UserID:      status.UserID,
		UserType:    status.UserType,
		IP:          status.IP,
		StartTime:   status.StartTime,
		LastUpdate:  status.LastUpdate,
		MergeInProgress: status.MergeInProgress,
	}

	// 复制已上传的分片
	result.Chunks = make(map[int]bool, len(status.Chunks))
	for chunk, uploaded := range status.Chunks {
		result.Chunks[chunk] = uploaded
	}

	return result, nil
}

// ListChunks 列出分片
func (u *ChunkUploader) ListChunks(ctx context.Context, uploadID string) ([]int, error) {
	u.chunkLock.Lock()
	status, exists := u.chunkUploads[uploadID]
	u.chunkLock.Unlock()

	if !exists {
		return nil, ErrChunkNotFound
	}

	chunks := make([]int, 0, len(status.Chunks))
	for chunk := range status.Chunks {
		chunks = append(chunks, chunk)
	}

	sort.Ints(chunks)
	return chunks, nil
}

// AbortUpload 中止上传
func (u *ChunkUploader) AbortUpload(ctx context.Context, uploadID string) error {
	u.chunkLock.Lock()
	_, exists := u.chunkUploads[uploadID]
	u.chunkLock.Unlock()

	if !exists {
		return ErrChunkNotFound
	}

	u.cleanupUpload(uploadID)
	return nil
}

// cleanupUpload 清理上传
func (u *ChunkUploader) cleanupUpload(uploadID string) {
	u.chunkLock.Lock()
	delete(u.chunkUploads, uploadID)
	u.chunkLock.Unlock()

	// 删除分片目录
	chunkDir := u.getChunkDir(uploadID)
	os.RemoveAll(chunkDir)

	// 删除临时文件
	tempFile := u.getTempFilePath(uploadID)
	os.Remove(tempFile)

	u.logger.Info("Cleaned up upload", "upload_id", uploadID)
}

// cleanupExpiredUpload 清理过期上传
func (u *ChunkUploader) cleanupExpiredUpload(uploadID string) {
	u.chunkLock.Lock()
	delete(u.chunkUploads, uploadID)
	u.chunkLock.Unlock()

	// 删除分片目录
	chunkDir := u.getChunkDir(uploadID)
	os.RemoveAll(chunkDir)

	// 删除临时文件
	tempFile := u.getTempFilePath(uploadID)
	os.Remove(tempFile)

	u.logger.Info("Cleaned up expired upload", "upload_id", uploadID)
}

// cleanupMergeStatus 清理合并状态
func (u *ChunkUploader) cleanupMergeStatus(uploadID string) {
	u.chunkLock.Lock()
	if status, exists := u.chunkUploads[uploadID]; exists {
		status.MergeInProgress = false
	}
	u.chunkLock.Unlock()
}

// getChunkDir 获取分片目录
func (u *ChunkUploader) getChunkDir(uploadID string) string {
	return filepath.Join(u.config.ChunkDir, uploadID)
}

// getChunkPath 获取分片路径
func (u *ChunkUploader) getChunkPath(uploadID string, chunkIndex int) string {
	chunkDir := u.getChunkDir(uploadID)
	return filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", chunkIndex))
}

// getTempFilePath 获取临时文件路径
func (u *ChunkUploader) getTempFilePath(uploadID string) string {
	return filepath.Join(u.config.TempDir, fmt.Sprintf("merged_%s", uploadID))
}

// CleanupExpiredUploads 清理过期的上传
func (u *ChunkUploader) CleanupExpiredUploads() {
	u.chunkLock.Lock()
	defer u.chunkLock.Unlock()

	timeout := time.Duration(u.config.ChunkTimeout) * time.Second
	now := time.Now()
	expiredIDs := make([]string, 0)

	// 查找过期的上传
	for id, status := range u.chunkUploads {
		if now.Sub(status.LastUpdate) > timeout {
			expiredIDs = append(expiredIDs, id)
		}
	}

	// 释放锁后清理
	u.chunkLock.Unlock()
	for _, id := range expiredIDs {
		u.cleanupExpiredUpload(id)
	}
	u.chunkLock.Lock()

	u.logger.Info("Cleaned up expired uploads", "count", len(expiredIDs))
}

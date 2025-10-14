// Package upload 提供文件上传接口和实现
package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"goweb/pkg/logger"
	"goweb/pkg/storage"
)

// UploadType 上传类型
type UploadType string

// 上传类型常量
const (
	UploadTypeSimple  UploadType = "simple"  // 简单上传
	UploadTypeChunked UploadType = "chunked" // 分片上传
)

// UploadRequest 上传请求
type UploadRequest struct {
	File        *multipart.FileHeader // 文件
	Description string                // 描述
	Tags        string                // 标签
	SubPath     string                // 子路径
	UserID      uint                  // 用户ID
	UserType    string                // 用户类型
	IP          string                // IP地址
	Metadata    map[string]interface{} // 元数据
}

// UploadResponse 上传响应
type UploadResponse struct {
	ID           uint      // 文件ID
	OriginalName string    // 原始文件名
	FileName     string    // 存储文件名
	FileSize     int64     // 文件大小
	FileType     string    // 文件类型
	MimeType     string    // MIME类型
	URL          string    // 访问URL
	Hash         string    // 文件哈希
	UploadTime   time.Time // 上传时间
	Metadata     map[string]interface{} // 元数据
}

// ChunkUploadRequest 分片上传请求
type ChunkUploadRequest struct {
	UploadID   string                // 上传ID
	File       *multipart.FileHeader // 文件分片
	FileName   string                // 文件名
	ChunkIndex int                   // 分片索引
	TotalChunks int                  // 总分片数
	ChunkSize  int64                 // 分片大小
	TotalSize  int64                 // 总大小
	FileHash   string                // 文件哈希
	SubPath    string                // 子路径
	UserID     uint                  // 用户ID
	UserType   string                // 用户类型
	IP         string                // IP地址
	Metadata   map[string]interface{} // 元数据
}

// ChunkUploadResponse 分片上传响应
type ChunkUploadResponse struct {
	UploadID   string // 上传ID
	ChunkIndex int    // 分片索引
	TotalChunks int   // 总分片数
	Completed  bool   // 是否完成
}

// ChunkMergeRequest 分片合并请求
type ChunkMergeRequest struct {
	UploadID   string // 上传ID
	FileName   string // 文件名
	TotalChunks int   // 总分片数
	FileHash   string // 文件哈希
	SubPath    string // 子路径
	UserID     uint   // 用户ID
	UserType   string // 用户类型
	IP         string // IP地址
	Metadata   map[string]interface{} // 元数据
}

// Uploader 上传接口
type Uploader interface {
	// Type 获取上传类型
	Type() UploadType
	
	// Upload 上传文件
	Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
	
	// CanHandle 是否可以处理
	CanHandle(req *UploadRequest) bool
}

// ChunkUploader 分片上传接口
type ChunkUploader interface {
	Uploader
	
	// InitiateUpload 初始化分片上传
	InitiateUpload(ctx context.Context, req *ChunkUploadRequest) (string, error)
	
	// UploadChunk 上传分片
	UploadChunk(ctx context.Context, req *ChunkUploadRequest) (*ChunkUploadResponse, error)
	
	// MergeChunks 合并分片
	MergeChunks(ctx context.Context, req *ChunkMergeRequest) (*UploadResponse, error)
	
	// GetUploadStatus 获取上传状态
	GetUploadStatus(ctx context.Context, uploadID string) (interface{}, error)
	
	// ListChunks 列出分片
	ListChunks(ctx context.Context, uploadID string) ([]int, error)
	
	// AbortUpload 中止上传
	AbortUpload(ctx context.Context, uploadID string) error
}

// FileProcessor 文件处理器接口
type FileProcessor interface {
	// Process 处理文件
	Process(ctx context.Context, fileInfo *storage.FileInfo) (*storage.FileInfo, error)
	
	// CanProcess 是否可以处理
	CanProcess(fileInfo *storage.FileInfo) bool
}

// BaseUploader 基础上传实现
type BaseUploader struct {
	uploadType UploadType
	storage    storage.Storage
	processors []FileProcessor
	logger     logger.Logger
}

// NewBaseUploader 创建基础上传实现
func NewBaseUploader(uploadType UploadType, storage storage.Storage, log logger.Logger) *BaseUploader {
	return &BaseUploader{
		uploadType: uploadType,
		storage:    storage,
		processors: make([]FileProcessor, 0),
		logger:     log,
	}
}

// Type 获取上传类型
func (u *BaseUploader) Type() UploadType {
	return u.uploadType
}

// AddProcessor 添加处理器
func (u *BaseUploader) AddProcessor(processor FileProcessor) {
	u.processors = append(u.processors, processor)
}

// ProcessFile 处理文件
func (u *BaseUploader) ProcessFile(ctx context.Context, fileInfo *storage.FileInfo) (*storage.FileInfo, error) {
	// 处理文件
	for _, processor := range u.processors {
		if processor.CanProcess(fileInfo) {
			var err error
			fileInfo, err = processor.Process(ctx, fileInfo)
			if err != nil {
				return nil, err
			}
		}
	}
	
	return fileInfo, nil
}

// SimpleUploader 简单上传实现
type SimpleUploader struct {
	*BaseUploader
}

// NewSimpleUploader 创建简单上传实现
func NewSimpleUploader(storage storage.Storage, log logger.Logger) *SimpleUploader {
	return &SimpleUploader{
		BaseUploader: NewBaseUploader(UploadTypeSimple, storage, log),
	}
}

// Upload 上传文件
func (u *SimpleUploader) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	// 上传文件
	fileInfo, err := u.storage.UploadFileWithContext(ctx, req.File, req.SubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 添加元数据
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			fileInfo.Metadata[k] = v
		}
	}
	fileInfo.Metadata["user_id"] = req.UserID
	fileInfo.Metadata["user_type"] = req.UserType
	fileInfo.Metadata["ip"] = req.IP
	fileInfo.Metadata["description"] = req.Description
	fileInfo.Metadata["tags"] = req.Tags

	// 处理文件
	fileInfo, err = u.ProcessFile(ctx, fileInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to process file: %w", err)
	}

	// 构建响应
	response := ConvertToResponse(fileInfo, 0) // ID将由存储层填充

	return response, nil
}

// CanHandle 是否可以处理
func (u *SimpleUploader) CanHandle(req *UploadRequest) bool {
	// 简单上传可以处理所有请求
	return true
}

// ConvertToResponse 转换为响应
func ConvertToResponse(fileInfo *storage.FileInfo, id uint) *UploadResponse {
	return &UploadResponse{
		ID:           id,
		OriginalName: fileInfo.OriginalName,
		FileName:     fileInfo.FileName,
		FileSize:     fileInfo.Size,
		FileType:     string(storage.GetFileType(fileInfo.MimeType)),
		MimeType:     fileInfo.MimeType,
		URL:          fileInfo.URL,
		Hash:         fileInfo.Hash,
		UploadTime:   fileInfo.UploadTime,
		Metadata:     fileInfo.Metadata,
	}
}

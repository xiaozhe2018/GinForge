package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"goweb/pkg/upload"
	"goweb/services/file-api/internal/model"
)

// ChunkUploadRequest 分片上传请求
type ChunkUploadRequest struct {
	UploadID   string `form:"upload_id"`
	ChunkIndex int    `form:"chunk_index" binding:"required,min=0"`
	TotalChunks int   `form:"total_chunks" binding:"required,min=1"`
	File       *multipart.FileHeader `form:"file"`
	FileName   string `form:"file_name" binding:"required"`
	TotalSize  int64  `form:"total_size" binding:"required,min=1"`
	ChunkSize  int64  `form:"chunk_size" binding:"required,min=1"`
	FileHash   string `form:"file_hash"`
	SubPath    string `form:"sub_path"`
	UserID     uint   `form:"user_id"`
	UserType   string `form:"user_type"`
	IP         string `form:"-"` // 由处理器设置
}

// ChunkUploadResponse 分片上传响应
type ChunkUploadResponse struct {
	UploadID    string `json:"upload_id"`
	ChunkIndex  int    `json:"chunk_index"`
	TotalChunks int    `json:"total_chunks"`
	Completed   bool   `json:"completed"`
}

// ChunkMergeRequest 分片合并请求
type ChunkMergeRequest struct {
	UploadID  string `form:"upload_id" binding:"required"`
	FileName  string `form:"file_name" binding:"required"`
	FileHash  string `form:"file_hash"`
	SubPath   string `form:"sub_path"`
	UserID    uint   `form:"user_id"`
	UserType  string `form:"user_type"`
	IP        string `form:"-"` // 由处理器设置
}

// InitiateChunkUpload 初始化分片上传
func (s *FileService) InitiateChunkUpload(ctx context.Context, req *ChunkUploadRequest) (string, error) {
	// 检查文件安全性
	if s.security != nil {
		if err := s.security.ValidateFile(ctx, req.FileName, "", req.TotalSize); err != nil {
			return "", fmt.Errorf("file validation failed: %w", err)
		}

		if err := s.security.CheckUploadLimit(ctx, req.UserID, req.TotalSize); err != nil {
			return "", fmt.Errorf("upload limit exceeded: %w", err)
		}
	}

	// 转换为上传包的请求
	chunkReq := &upload.ChunkUploadRequest{
		FileName:    req.FileName,
		TotalChunks: req.TotalChunks,
		TotalSize:   req.TotalSize,
		ChunkSize:   req.ChunkSize,
		FileHash:    req.FileHash,
		SubPath:     req.SubPath,
		UserID:      req.UserID,
		UserType:    req.UserType,
		IP:          req.IP,
		Metadata:    make(map[string]interface{}),
	}

	// 初始化分片上传
	return s.chunkUploader.InitiateUpload(ctx, chunkReq)
}

// UploadChunk 上传分片
func (s *FileService) UploadChunk(ctx context.Context, req *ChunkUploadRequest) (*ChunkUploadResponse, error) {
	// 转换为上传包的请求
	chunkReq := &upload.ChunkUploadRequest{
		UploadID:    req.UploadID,
		ChunkIndex:  req.ChunkIndex,
		TotalChunks: req.TotalChunks,
		File:        req.File,
		FileName:    req.FileName,
		TotalSize:   req.TotalSize,
		ChunkSize:   req.ChunkSize,
		FileHash:    req.FileHash,
		SubPath:     req.SubPath,
		UserID:      req.UserID,
		UserType:    req.UserType,
		IP:          req.IP,
		Metadata:    make(map[string]interface{}),
	}

	// 上传分片
	result, err := s.chunkUploader.UploadChunk(ctx, chunkReq)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunk: %w", err)
	}

	return &ChunkUploadResponse{
		UploadID:    result.UploadID,
		ChunkIndex:  result.ChunkIndex,
		TotalChunks: result.TotalChunks,
		Completed:   result.Completed,
	}, nil
}

// MergeChunks 合并分片
func (s *FileService) MergeChunks(ctx context.Context, req *ChunkMergeRequest) (*UploadResponse, error) {
	startTime := time.Now()

	// 转换为上传包的请求
	mergeReq := &upload.ChunkMergeRequest{
		UploadID:  req.UploadID,
		FileName:  req.FileName,
		FileHash:  req.FileHash,
		SubPath:   req.SubPath,
		UserID:    req.UserID,
		UserType:  req.UserType,
		IP:        req.IP,
		Metadata:  make(map[string]interface{}),
	}

	// 合并分片
	result, err := s.chunkUploader.MergeChunks(ctx, mergeReq)
	if err != nil {
		s.logger.Error("Failed to merge chunks", "error", err, "upload_id", req.UploadID)
		return nil, fmt.Errorf("failed to merge chunks: %w", err)
	}

	// 保存文件记录到数据库
	fileRecord := &model.FileRecord{
		OriginalName: result.OriginalName,
		FileName:     result.FileName,
		RelativePath: result.FileName,
		FileSize:     result.FileSize,
		MimeType:     result.MimeType,
		FileType:     result.FileType,
		Hash:         result.Hash,
		StorageType:  string(s.storage.Type()),
		URL:          result.URL,
		UploadedBy:   req.UserID,
		UserType:     req.UserType,
		UploadIP:     req.IP,
		UploadTime:   result.UploadTime,
		Description:  "Chunked upload",
		Tags:         "chunked",
		Status:       1, // 正常状态
	}

	if err := s.repo.Create(ctx, fileRecord); err != nil {
		s.logger.Error("Failed to save file record", "error", err)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// 记录上传
	if s.security != nil {
		if err := s.security.RecordUpload(ctx, req.UserID, result.FileSize); err != nil {
			s.logger.Warn("Failed to record upload", "error", err)
		}
	}

	// 创建上传日志
	s.createChunkUploadLog(ctx, req, result, nil, startTime)

	// 更新响应中的ID
	result.ID = fileRecord.ID

	return &UploadResponse{
		ID:           fileRecord.ID,
		OriginalName: result.OriginalName,
		FileName:     result.FileName,
		FileSize:     result.FileSize,
		FileType:     result.FileType,
		MimeType:     result.MimeType,
		URL:          result.URL,
		Hash:         result.Hash,
		UploadTime:   result.UploadTime,
	}, nil
}

// GetChunkUploadStatus 获取分片上传状态
func (s *FileService) GetChunkUploadStatus(ctx context.Context, uploadID string) (interface{}, error) {
	return s.chunkUploader.GetUploadStatus(ctx, uploadID)
}

// ListChunks 列出分片
func (s *FileService) ListChunks(ctx context.Context, uploadID string) ([]int, error) {
	return s.chunkUploader.ListChunks(ctx, uploadID)
}

// AbortChunkUpload 中止分片上传
func (s *FileService) AbortChunkUpload(ctx context.Context, uploadID string) error {
	return s.chunkUploader.AbortUpload(ctx, uploadID)
}

// 创建分片上传日志
func (s *FileService) createChunkUploadLog(ctx context.Context, req *ChunkMergeRequest, result *upload.UploadResponse, err error, startTime time.Time) {
	duration := time.Since(startTime).Milliseconds()

	log := &model.FileUploadLog{
		FileName:   req.FileName,
		FileSize:   result.FileSize,
		UploadedBy: req.UserID,
		UploadIP:   req.IP,
		UploadTime: time.Now(),
		Success:    err == nil,
		Duration:   duration,
		UploadType: "chunked",
	}

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	if err := s.repo.CreateUploadLog(ctx, log); err != nil {
		s.logger.Error("failed to create upload log", err)
	}
}
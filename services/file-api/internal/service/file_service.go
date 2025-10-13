package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"gorm.io/gorm"

	"goweb/pkg/logger"
	"goweb/pkg/storage"
	"goweb/services/file-api/internal/model"
	"goweb/services/file-api/internal/repository"
)

// FileService 文件服务
type FileService struct {
	repo           repository.FileRepository
	storageService *StorageService
	logger         logger.Logger
	db             *gorm.DB
}

// NewFileService 创建文件服务
func NewFileService(repo repository.FileRepository, storageService *StorageService, db *gorm.DB, log logger.Logger) *FileService {
	return &FileService{
		repo:           repo,
		storageService: storageService,
		logger:         log,
		db:             db,
	}
}

// UploadRequest 上传请求
type UploadRequest struct {
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Description string                `form:"description"`
	Tags        string                `form:"tags"`
	SubPath     string                `form:"sub_path"`
	UserID      uint                  `form:"user_id"`
	UserType    string                `form:"user_type"`
	IP          string                `form:"-"` // 由处理器设置
}

// UploadResponse 上传响应
type UploadResponse struct {
	ID           uint      `json:"id"`
	OriginalName string    `json:"original_name"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	MimeType     string    `json:"mime_type"`
	URL          string    `json:"url"`
	Hash         string    `json:"hash"`
	UploadTime   time.Time `json:"upload_time"`
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	startTime := time.Now()

	// 1. 上传文件到存储
	fileInfo, err := s.storageService.UploadFile(ctx, req.File, req.SubPath)
	if err != nil {
		// 记录上传失败日志
		s.createUploadLog(ctx, req, nil, err, startTime)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 2. 获取文件类型
	fileType := s.storageService.GetFileTypeByMime(fileInfo.MimeType)
	if fileType == "" {
		fileType = s.storageService.GetFileTypeByExt(fileInfo.OriginalName)
	}

	// 3. 保存文件记录到数据库
	fileRecord := &model.FileRecord{
		OriginalName: fileInfo.OriginalName,
		FileName:     fileInfo.FileName,
		RelativePath: fileInfo.RelativePath,
		FileSize:     fileInfo.Size,
		MimeType:     fileInfo.MimeType,
		FileType:     fileType,
		Hash:         fileInfo.Hash,
		StorageType:  s.storageService.GetStorageType(),
		URL:          fileInfo.URL,
		UploadedBy:   req.UserID,
		UserType:     req.UserType,
		UploadIP:     req.IP,
		UploadTime:   fileInfo.UploadTime,
		Description:  req.Description,
		Tags:         req.Tags,
		Status:       1, // 正常状态
	}

	if err := s.repo.Create(ctx, fileRecord); err != nil {
		// 记录上传失败日志
		s.createUploadLog(ctx, req, fileInfo, err, startTime)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// 4. 记录上传成功日志
	s.createUploadLog(ctx, req, fileInfo, nil, startTime)

	// 5. 构建响应
	response := &UploadResponse{
		ID:           fileRecord.ID,
		OriginalName: fileRecord.OriginalName,
		FileName:     fileRecord.FileName,
		FileSize:     fileRecord.FileSize,
		FileType:     fileRecord.FileType,
		MimeType:     fileRecord.MimeType,
		URL:          fileRecord.URL,
		Hash:         fileRecord.Hash,
		UploadTime:   fileRecord.UploadTime,
	}

	return response, nil
}

// 创建上传日志
func (s *FileService) createUploadLog(ctx context.Context, req *UploadRequest, fileInfo *storage.FileInfo, err error, startTime time.Time) {
	duration := time.Since(startTime).Milliseconds()

	log := &model.FileUploadLog{
		FileName:   req.File.Filename,
		FileSize:   req.File.Size,
		UploadedBy: req.UserID,
		UploadIP:   req.IP,
		UploadTime: time.Now(),
		Success:    err == nil,
		Duration:   duration,
	}

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	if err := s.repo.CreateUploadLog(ctx, log); err != nil {
		s.logger.Error("failed to create upload log", err)
	}
}

// DownloadRequest 下载请求
type DownloadRequest struct {
	FileID uint   `uri:"id" binding:"required"`
	UserID uint   `form:"user_id"`
	IP     string `form:"-"` // 由处理器设置
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(ctx context.Context, req *DownloadRequest) (*model.FileRecord, error) {
	// 1. 查询文件记录
	fileRecord, err := s.repo.Get(ctx, req.FileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 2. 检查文件是否存在
	if !s.storageService.FileExists(ctx, fileRecord.RelativePath) {
		return nil, fmt.Errorf("file not found in storage")
	}

	// 3. 增加下载次数
	if err := s.repo.IncrementDownloadCount(ctx, req.FileID); err != nil {
		s.logger.Warn("failed to increment download count", "error", err)
	}

	// 4. 记录下载日志
	s.createDownloadLog(ctx, fileRecord, req.UserID, req.IP, nil)

	return fileRecord, nil
}

// 创建下载日志
func (s *FileService) createDownloadLog(ctx context.Context, file *model.FileRecord, userID uint, ip string, err error) {
	log := &model.FileDownloadLog{
		FileID:       file.ID,
		FileName:     file.FileName,
		DownloadedBy: userID,
		DownloadIP:   ip,
		DownloadTime: time.Now(),
		Success:      err == nil,
	}

	if err := s.repo.CreateDownloadLog(ctx, log); err != nil {
		s.logger.Error("failed to create download log", err)
	}
}

// ListFilesRequest 列出文件请求
type ListFilesRequest struct {
	UserID   uint   `form:"user_id"`
	FileType string `form:"file_type"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
}

// ListFilesResponse 列出文件响应
type ListFilesResponse struct {
	Total int64               `json:"total"`
	List  []*model.FileRecord `json:"list"`
}

// ListFiles 列出文件
func (s *FileService) ListFiles(ctx context.Context, req *ListFilesRequest) (*ListFilesResponse, error) {
	var files []*model.FileRecord
	var total int64
	var err error

	// 根据条件查询
	if req.UserID > 0 {
		files, total, err = s.repo.FindByUserID(ctx, req.UserID, req.Page, req.PageSize)
	} else if req.FileType != "" {
		files, total, err = s.repo.FindByType(ctx, req.FileType, req.Page, req.PageSize)
	} else {
		// 查询所有文件
		files, total, err = s.repo.List(ctx, req.Page, req.PageSize)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return &ListFilesResponse{
		Total: total,
		List:  files,
	}, nil
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	FileID uint `uri:"id" binding:"required"`
	UserID uint `form:"user_id"`
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, req *DeleteFileRequest) error {
	// 1. 查询文件记录
	fileRecord, err := s.repo.Get(ctx, req.FileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// 2. 检查权限（如果需要）
	if req.UserID > 0 && fileRecord.UploadedBy != req.UserID {
		return fmt.Errorf("permission denied")
	}

	// 3. 软删除文件记录
	fileRecord.Status = 2 // 已删除状态
	if err := s.repo.Update(ctx, fileRecord); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	// 注意：这里不实际删除存储中的文件，可以通过定期清理任务处理

	return nil
}

// GetFileRequest 获取文件请求
type GetFileRequest struct {
	FileID uint `uri:"id" binding:"required"`
}

// GetFile 获取文件信息
func (s *FileService) GetFile(ctx context.Context, req *GetFileRequest) (*model.FileRecord, error) {
	return s.repo.Get(ctx, req.FileID)
}

// GetFileByHashRequest 根据哈希获取文件请求
type GetFileByHashRequest struct {
	Hash string `uri:"hash" binding:"required"`
}

// GetFileByHash 根据哈希获取文件
func (s *FileService) GetFileByHash(ctx context.Context, req *GetFileByHashRequest) (*model.FileRecord, error) {
	return s.repo.FindByHash(ctx, req.Hash)
}

// GetStatisticsResponse 获取统计信息响应
type GetStatisticsResponse struct {
	TotalCount     int64            `json:"total_count"`
	TotalSize      int64            `json:"total_size"`
	TypeStatistics map[string]int64 `json:"type_statistics"`
}

// GetStatistics 获取统计信息
func (s *FileService) GetStatistics(ctx context.Context) (*GetStatisticsResponse, error) {
	// 获取总数量
	totalCount, err := s.repo.GetTotalCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// 获取总大小
	totalSize, err := s.repo.GetTotalSize(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total size: %w", err)
	}

	// 获取类型统计
	typeStatistics, err := s.repo.GetTypeStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get type statistics: %w", err)
	}

	return &GetStatisticsResponse{
		TotalCount:     totalCount,
		TotalSize:      totalSize,
		TypeStatistics: typeStatistics,
	}, nil
}

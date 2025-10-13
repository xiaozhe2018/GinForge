package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goweb/pkg/logger"
	"goweb/pkg/storage"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 阿里云OSS
	StorageTypeS3    StorageType = "s3"    // AWS S3
	StorageTypeMinio StorageType = "minio" // MinIO
)

// StorageConfig 存储配置
type StorageConfig struct {
	Type          StorageType `json:"type"`            // 存储类型
	LocalBasePath string      `json:"local_base_path"` // 本地存储基础路径
	OSSConfig     OSSConfig   `json:"oss_config"`      // OSS配置
	S3Config      S3Config    `json:"s3_config"`       // S3配置
	MinioConfig   MinioConfig `json:"minio_config"`    // MinIO配置
	URLPrefix     string      `json:"url_prefix"`      // URL前缀
	MaxFileSize   int64       `json:"max_file_size"`   // 最大文件大小(字节)
}

// OSSConfig OSS配置
type OSSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
}

// S3Config S3配置
type S3Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
}

// MinioConfig MinIO配置
type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	UseSSL          bool   `json:"use_ssl"`
}

// StorageService 存储服务
type StorageService struct {
	config  *StorageConfig
	storage storage.Storage
	logger  logger.Logger
}

// GetStorageType 获取存储类型
func (s *StorageService) GetStorageType() string {
	return string(s.config.Type)
}

// NewStorageService 创建存储服务
func NewStorageService(config *StorageConfig, log logger.Logger) (*StorageService, error) {
	var store storage.Storage
	var err error

	switch config.Type {
	case StorageTypeLocal:
		store = storage.NewLocalStorage(config.LocalBasePath, log)
	case StorageTypeOSS:
		// TODO: 实现OSS存储
		return nil, errors.New("OSS storage not implemented yet")
	case StorageTypeS3:
		// TODO: 实现S3存储
		return nil, errors.New("S3 storage not implemented yet")
	case StorageTypeMinio:
		// TODO: 实现MinIO存储
		return nil, errors.New("MinIO storage not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}

	if err != nil {
		return nil, err
	}

	return &StorageService{
		config:  config,
		storage: store,
		logger:  log,
	}, nil
}

// UploadFile 上传文件
func (s *StorageService) UploadFile(ctx context.Context, file *multipart.FileHeader, subPath string) (*storage.FileInfo, error) {
	// 检查文件大小
	if s.config.MaxFileSize > 0 && file.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d > %d", file.Size, s.config.MaxFileSize)
	}

	// 上传文件
	fileInfo, err := s.storage.UploadFile(file, subPath)
	if err != nil {
		return nil, err
	}

	// 添加URL前缀
	if s.config.URLPrefix != "" {
		fileInfo.URL = s.config.URLPrefix + "/" + fileInfo.RelativePath
	}

	return fileInfo, nil
}

// SaveFile 保存文件
func (s *StorageService) SaveFile(ctx context.Context, data []byte, fileName, subPath string) (*storage.FileInfo, error) {
	// 检查文件大小
	if s.config.MaxFileSize > 0 && int64(len(data)) > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d > %d", len(data), s.config.MaxFileSize)
	}

	// 保存文件
	fileInfo, err := s.storage.SaveFile(data, fileName, subPath)
	if err != nil {
		return nil, err
	}

	// 添加URL前缀
	if s.config.URLPrefix != "" {
		fileInfo.URL = s.config.URLPrefix + "/" + fileInfo.RelativePath
	}

	return fileInfo, nil
}

// GetFile 获取文件
func (s *StorageService) GetFile(ctx context.Context, relativePath string) (*storage.FileInfo, error) {
	return s.storage.GetFile(relativePath)
}

// DeleteFile 删除文件
func (s *StorageService) DeleteFile(ctx context.Context, relativePath string) error {
	return s.storage.DeleteFile(relativePath)
}

// ListFiles 列出文件
func (s *StorageService) ListFiles(ctx context.Context, subPath string) ([]*storage.FileInfo, error) {
	return s.storage.ListFiles(subPath)
}

// FileExists 检查文件是否存在
func (s *StorageService) FileExists(ctx context.Context, relativePath string) bool {
	return s.storage.FileExists(relativePath)
}

// GetFileSize 获取文件大小
func (s *StorageService) GetFileSize(ctx context.Context, relativePath string) (int64, error) {
	return s.storage.GetFileSize(relativePath)
}

// GetContentType 获取文件内容类型
func (s *StorageService) GetContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 只读取前512字节用于判断内容类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// 使用net/http包检测内容类型
	contentType := detectContentType(buffer)
	return contentType, nil
}

// detectContentType 检测内容类型
func detectContentType(data []byte) string {
	// 简单的MIME类型检测
	if len(data) > 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if len(data) > 4 && data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
		return "image/png"
	}
	if len(data) > 4 && data[0] == 'G' && data[1] == 'I' && data[2] == 'F' && data[3] == '8' {
		return "image/gif"
	}
	if len(data) > 4 && data[0] == '%' && data[1] == 'P' && data[2] == 'D' && data[3] == 'F' {
		return "application/pdf"
	}
	// 默认返回二进制流
	return "application/octet-stream"
}

// GetFileTypeByMime 根据MIME类型获取文件类型
func (s *StorageService) GetFileTypeByMime(mimeType string) string {
	fileType := storage.GetFileType(mimeType)
	return string(fileType)
}

// GetFileTypeByExt 根据扩展名获取文件类型
func (s *StorageService) GetFileTypeByExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
		return string(storage.FileTypeImage)
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return string(storage.FileTypeVideo)
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return string(storage.FileTypeAudio)
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md":
		return string(storage.FileTypeDoc)
	default:
		return string(storage.FileTypeOther)
	}
}

// Cleanup 清理过期文件
func (s *StorageService) Cleanup(ctx context.Context, subPath string, maxAge time.Duration) error {
	return s.storage.Cleanup(subPath, maxAge)
}

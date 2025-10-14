// Package storage 提供文件存储接口和实现
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

// StorageType 存储类型
type StorageType string

// 存储类型常量
const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 阿里云OSS
	StorageTypeS3    StorageType = "s3"    // AWS S3
	StorageTypeCOS   StorageType = "cos"   // 腾讯云COS
	StorageTypeQiniu StorageType = "qiniu" // 七牛云
	StorageTypeMinio StorageType = "minio" // MinIO
)

// FileType 文件类型
type FileType string

// 文件类型常量
const (
	FileTypeImage   FileType = "image"   // 图片
	FileTypeVideo   FileType = "video"   // 视频
	FileTypeAudio   FileType = "audio"   // 音频
	FileTypeDoc     FileType = "document" // 文档
	FileTypeArchive FileType = "archive" // 压缩文件
	FileTypeOther   FileType = "other"   // 其他
)

// FileInfo 文件信息
type FileInfo struct {
	OriginalName string    `json:"original_name"` // 原始文件名
	FileName     string    `json:"file_name"`     // 存储文件名
	RelativePath string    `json:"relative_path"` // 相对路径
	Size         int64     `json:"size"`          // 文件大小
	MimeType     string    `json:"mime_type"`     // MIME类型
	Hash         string    `json:"hash"`          // 文件哈希
	URL          string    `json:"url"`           // 访问URL
	UploadTime   time.Time `json:"upload_time"`   // 上传时间
	Metadata     map[string]interface{} `json:"metadata"` // 元数据
}

// Storage 存储接口
type Storage interface {
	// Name 获取存储名称
	Name() string
	
	// Type 获取存储类型
	Type() StorageType
	
	// UploadFile 上传文件
	UploadFile(file *multipart.FileHeader, subPath string) (*FileInfo, error)
	
	// UploadFileWithContext 带上下文的上传文件
	UploadFileWithContext(ctx context.Context, file *multipart.FileHeader, subPath string) (*FileInfo, error)
	
	// SaveFile 保存文件
	SaveFile(data []byte, fileName, subPath string) (*FileInfo, error)
	
	// SaveFileWithContext 带上下文的保存文件
	SaveFileWithContext(ctx context.Context, data []byte, fileName, subPath string) (*FileInfo, error)
	
	// SaveFileFromReader 从读取器保存文件
	SaveFileFromReader(reader io.Reader, size int64, fileName, subPath string) (*FileInfo, error)
	
	// SaveFileFromReaderWithContext 带上下文从读取器保存文件
	SaveFileFromReaderWithContext(ctx context.Context, reader io.Reader, size int64, fileName, subPath string) (*FileInfo, error)
	
	// GetFile 获取文件
	GetFile(relativePath string) (*FileInfo, error)
	
	// GetFileWithContext 带上下文获取文件
	GetFileWithContext(ctx context.Context, relativePath string) (*FileInfo, error)
	
	// DeleteFile 删除文件
	DeleteFile(relativePath string) error
	
	// DeleteFileWithContext 带上下文删除文件
	DeleteFileWithContext(ctx context.Context, relativePath string) error
	
	// ListFiles 列出文件
	ListFiles(subPath string) ([]*FileInfo, error)
	
	// ListFilesWithContext 带上下文列出文件
	ListFilesWithContext(ctx context.Context, subPath string) ([]*FileInfo, error)
	
	// FileExists 检查文件是否存在
	FileExists(relativePath string) bool
	
	// FileExistsWithContext 带上下文检查文件是否存在
	FileExistsWithContext(ctx context.Context, relativePath string) bool
	
	// GetFileSize 获取文件大小
	GetFileSize(relativePath string) (int64, error)
	
	// GetFileSizeWithContext 带上下文获取文件大小
	GetFileSizeWithContext(ctx context.Context, relativePath string) (int64, error)
	
	// GetFileContent 获取文件内容
	GetFileContent(relativePath string) ([]byte, error)
	
	// GetFileContentWithContext 带上下文获取文件内容
	GetFileContentWithContext(ctx context.Context, relativePath string) ([]byte, error)
	
	// GetFileReader 获取文件读取器
	GetFileReader(relativePath string) (io.ReadCloser, error)
	
	// GetFileReaderWithContext 带上下文获取文件读取器
	GetFileReaderWithContext(ctx context.Context, relativePath string) (io.ReadCloser, error)
	
	// GetFileURL 获取文件URL
	GetFileURL(relativePath string) string
	
	// GetFileURLWithContext 带上下文获取文件URL
	GetFileURLWithContext(ctx context.Context, relativePath string) string
	
	// GetSignedURL 获取签名URL
	GetSignedURL(relativePath string, expireSeconds int) (string, error)
	
	// GetSignedURLWithContext 带上下文获取签名URL
	GetSignedURLWithContext(ctx context.Context, relativePath string, expireSeconds int) (string, error)
	
	// Cleanup 清理过期文件
	Cleanup(subPath string, maxAge time.Duration) error
	
	// CleanupWithContext 带上下文清理过期文件
	CleanupWithContext(ctx context.Context, subPath string, maxAge time.Duration) error
}

// StorageConfig 存储配置接口
type StorageConfig interface {
	// GetType 获取存储类型
	GetType() StorageType
	
	// GetConfig 获取配置
	GetConfig() map[string]interface{}
	
	// GetString 获取字符串配置
	GetString(key string) string
	
	// GetInt 获取整数配置
	GetInt(key string) int
	
	// GetInt64 获取64位整数配置
	GetInt64(key string) int64
	
	// GetBool 获取布尔配置
	GetBool(key string) bool
	
	// GetDuration 获取时间配置
	GetDuration(key string) time.Duration
}

// StorageOption 存储选项函数
type StorageOption func(storage Storage)

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) StorageOption {
	return func(storage Storage) {
		if configurable, ok := storage.(interface{ SetBaseURL(string) }); ok {
			configurable.SetBaseURL(baseURL)
		}
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) StorageOption {
	return func(storage Storage) {
		if configurable, ok := storage.(interface{ SetTimeout(time.Duration) }); ok {
			configurable.SetTimeout(timeout)
		}
	}
}

// WithMaxFileSize 设置最大文件大小
func WithMaxFileSize(maxFileSize int64) StorageOption {
	return func(storage Storage) {
		if configurable, ok := storage.(interface{ SetMaxFileSize(int64) }); ok {
			configurable.SetMaxFileSize(maxFileSize)
		}
	}
}

// GetFileType 根据MIME类型获取文件类型
func GetFileType(mimeType string) FileType {
	switch {
	case mimeType == "":
		return FileTypeOther
	case mimeType == "application/octet-stream":
		return FileTypeOther
	case mimeType == "text/plain":
		return FileTypeDoc
	case mimeType == "application/pdf":
		return FileTypeDoc
	case mimeType == "application/msword", mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return FileTypeDoc
	case mimeType == "application/vnd.ms-excel", mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return FileTypeDoc
	case mimeType == "application/vnd.ms-powerpoint", mimeType == "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return FileTypeDoc
	case mimeType == "application/zip", mimeType == "application/x-rar-compressed", mimeType == "application/x-7z-compressed":
		return FileTypeArchive
	case mimeType == "audio/mpeg", mimeType == "audio/wav", mimeType == "audio/ogg", mimeType == "audio/webm":
		return FileTypeAudio
	case mimeType == "video/mp4", mimeType == "video/webm", mimeType == "video/ogg", mimeType == "video/quicktime":
		return FileTypeVideo
	case mimeType == "image/jpeg", mimeType == "image/png", mimeType == "image/gif", mimeType == "image/webp", mimeType == "image/svg+xml":
		return FileTypeImage
	default:
		if len(mimeType) >= 5 {
			prefix := mimeType[:5]
			switch prefix {
			case "image":
				return FileTypeImage
			case "video":
				return FileTypeVideo
			case "audio":
				return FileTypeAudio
			}
		}
		return FileTypeOther
	}
}

// BaseStorage 基础存储实现
type BaseStorage struct {
	storageType StorageType
	name        string
	baseURL     string
	timeout     time.Duration
	maxFileSize int64
}

// NewBaseStorage 创建基础存储
func NewBaseStorage(storageType StorageType, name string) *BaseStorage {
	return &BaseStorage{
		storageType: storageType,
		name:        name,
		timeout:     30 * time.Second,
	}
}

// Name 获取存储名称
func (s *BaseStorage) Name() string {
	return s.name
}

// Type 获取存储类型
func (s *BaseStorage) Type() StorageType {
	return s.storageType
}

// SetBaseURL 设置基础URL
func (s *BaseStorage) SetBaseURL(baseURL string) {
	s.baseURL = baseURL
}

// SetTimeout 设置超时时间
func (s *BaseStorage) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// SetMaxFileSize 设置最大文件大小
func (s *BaseStorage) SetMaxFileSize(maxFileSize int64) {
	s.maxFileSize = maxFileSize
}

// GetBaseURL 获取基础URL
func (s *BaseStorage) GetBaseURL() string {
	return s.baseURL
}

// GetTimeout 获取超时时间
func (s *BaseStorage) GetTimeout() time.Duration {
	return s.timeout
}

// GetMaxFileSize 获取最大文件大小
func (s *BaseStorage) GetMaxFileSize() int64 {
	return s.maxFileSize
}

// 以下方法需要由实现类提供

// UploadFile 上传文件
func (s *BaseStorage) UploadFile(file *multipart.FileHeader, subPath string) (*FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// SaveFile 保存文件
func (s *BaseStorage) SaveFile(data []byte, fileName, subPath string) (*FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetFile 获取文件
func (s *BaseStorage) GetFile(relativePath string) (*FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// DeleteFile 删除文件
func (s *BaseStorage) DeleteFile(relativePath string) error {
	return fmt.Errorf("not implemented")
}

// ListFiles 列出文件
func (s *BaseStorage) ListFiles(subPath string) ([]*FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// FileExists 检查文件是否存在
func (s *BaseStorage) FileExists(relativePath string) bool {
	return false
}

// GetFileSize 获取文件大小
func (s *BaseStorage) GetFileSize(relativePath string) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

// GetFileContent 获取文件内容
func (s *BaseStorage) GetFileContent(relativePath string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetFileReader 获取文件读取器
func (s *BaseStorage) GetFileReader(relativePath string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetFileURL 获取文件URL
func (s *BaseStorage) GetFileURL(relativePath string) string {
	return ""
}

// GetSignedURL 获取签名URL
func (s *BaseStorage) GetSignedURL(relativePath string, expireSeconds int) (string, error) {
	return "", fmt.Errorf("not implemented")
}

// Cleanup 清理过期文件
func (s *BaseStorage) Cleanup(subPath string, maxAge time.Duration) error {
	return fmt.Errorf("not implemented")
}

// UploadFileWithContext 带上下文的上传文件
func (s *BaseStorage) UploadFileWithContext(ctx context.Context, file *multipart.FileHeader, subPath string) (*FileInfo, error) {
	// 默认实现调用不带上下文的方法
	return s.UploadFile(file, subPath)
}

// SaveFileWithContext 带上下文的保存文件
func (s *BaseStorage) SaveFileWithContext(ctx context.Context, data []byte, fileName, subPath string) (*FileInfo, error) {
	// 默认实现调用不带上下文的方法
	return s.SaveFile(data, fileName, subPath)
}

// GetFileWithContext 带上下文获取文件
func (s *BaseStorage) GetFileWithContext(ctx context.Context, relativePath string) (*FileInfo, error) {
	// 默认实现调用不带上下文的方法
	return s.GetFile(relativePath)
}

// DeleteFileWithContext 带上下文删除文件
func (s *BaseStorage) DeleteFileWithContext(ctx context.Context, relativePath string) error {
	// 默认实现调用不带上下文的方法
	return s.DeleteFile(relativePath)
}

// ListFilesWithContext 带上下文列出文件
func (s *BaseStorage) ListFilesWithContext(ctx context.Context, subPath string) ([]*FileInfo, error) {
	// 默认实现调用不带上下文的方法
	return s.ListFiles(subPath)
}

// FileExistsWithContext 带上下文检查文件是否存在
func (s *BaseStorage) FileExistsWithContext(ctx context.Context, relativePath string) bool {
	// 默认实现调用不带上下文的方法
	return s.FileExists(relativePath)
}

// GetFileSizeWithContext 带上下文获取文件大小
func (s *BaseStorage) GetFileSizeWithContext(ctx context.Context, relativePath string) (int64, error) {
	// 默认实现调用不带上下文的方法
	return s.GetFileSize(relativePath)
}

// GetFileContentWithContext 带上下文获取文件内容
func (s *BaseStorage) GetFileContentWithContext(ctx context.Context, relativePath string) ([]byte, error) {
	// 默认实现调用不带上下文的方法
	return s.GetFileContent(relativePath)
}

// GetFileReaderWithContext 带上下文获取文件读取器
func (s *BaseStorage) GetFileReaderWithContext(ctx context.Context, relativePath string) (io.ReadCloser, error) {
	// 默认实现调用不带上下文的方法
	return s.GetFileReader(relativePath)
}

// GetFileURLWithContext 带上下文获取文件URL
func (s *BaseStorage) GetFileURLWithContext(ctx context.Context, relativePath string) string {
	// 默认实现调用不带上下文的方法
	return s.GetFileURL(relativePath)
}

// GetSignedURLWithContext 带上下文获取签名URL
func (s *BaseStorage) GetSignedURLWithContext(ctx context.Context, relativePath string, expireSeconds int) (string, error) {
	// 默认实现调用不带上下文的方法
	return s.GetSignedURL(relativePath, expireSeconds)
}

// CleanupWithContext 带上下文清理过期文件
func (s *BaseStorage) CleanupWithContext(ctx context.Context, subPath string, maxAge time.Duration) error {
	// 默认实现调用不带上下文的方法
	return s.Cleanup(subPath, maxAge)
}

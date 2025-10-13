package storage

import (
	"mime/multipart"
	"strings"
	"time"
)

// FileInfo 文件信息
type FileInfo struct {
	OriginalName string    `json:"original_name"` // 原始文件名
	FileName     string    `json:"file_name"`     // 存储文件名
	FilePath     string    `json:"file_path"`     // 完整文件路径
	RelativePath string    `json:"relative_path"` // 相对路径
	Size         int64     `json:"size"`          // 文件大小
	MimeType     string    `json:"mime_type"`     // MIME 类型
	Hash         string    `json:"hash"`          // 文件哈希
	UploadTime   time.Time `json:"upload_time"`   // 上传时间
	URL          string    `json:"url,omitempty"` // 访问 URL
}

// Storage 存储接口
type Storage interface {
	// 文件操作
	UploadFile(file *multipart.FileHeader, subPath string) (*FileInfo, error)
	SaveFile(data []byte, fileName, subPath string) (*FileInfo, error)
	GetFile(relativePath string) (*FileInfo, error)
	DeleteFile(relativePath string) error
	ListFiles(subPath string) ([]*FileInfo, error)
	FileExists(relativePath string) bool
	GetFileSize(relativePath string) (int64, error)

	// 管理操作
	Cleanup(subPath string, maxAge time.Duration) error
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int64    `json:"max_size"`      // 最大文件大小
	AllowedTypes []string `json:"allowed_types"` // 允许的文件类型
	AllowedExts  []string `json:"allowed_exts"`  // 允许的扩展名
	SubPath      string   `json:"sub_path"`      // 子路径
	GenerateName bool     `json:"generate_name"` // 是否生成新文件名
}

// UploadResult 上传结果
type UploadResult struct {
	FileInfo *FileInfo `json:"file_info"`
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// FileType 文件类型
type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeVideo FileType = "video"
	FileTypeAudio FileType = "audio"
	FileTypeDoc   FileType = "document"
	FileTypeOther FileType = "other"
)

// GetFileType 根据 MIME 类型获取文件类型
func GetFileType(mimeType string) FileType {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return FileTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return FileTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return FileTypeAudio
	case strings.HasPrefix(mimeType, "application/pdf") ||
		strings.HasPrefix(mimeType, "application/msword") ||
		strings.HasPrefix(mimeType, "application/vnd.openxmlformats-officedocument"):
		return FileTypeDoc
	default:
		return FileTypeOther
	}
}

// IsImage 是否为图片
func (ft FileType) IsImage() bool {
	return ft == FileTypeImage
}

// IsVideo 是否为视频
func (ft FileType) IsVideo() bool {
	return ft == FileTypeVideo
}

// IsAudio 是否为音频
func (ft FileType) IsAudio() bool {
	return ft == FileTypeAudio
}

// IsDocument 是否为文档
func (ft FileType) IsDocument() bool {
	return ft == FileTypeDoc
}

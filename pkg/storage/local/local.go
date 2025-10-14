// Package local 提供本地文件存储实现
package local

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

// LocalStorage 本地存储实现
type LocalStorage struct {
	*storage.BaseStorage
	basePath      string
	tempDir       string
	autoCreateDir bool
	fileMode      os.FileMode
	dirMode       os.FileMode
	logger        logger.Logger
}

// Config 本地存储配置
type Config struct {
	BasePath      string      // 基础路径
	TempDir       string      // 临时目录
	AutoCreateDir bool        // 自动创建目录
	FileMode      os.FileMode // 文件权限
	DirMode       os.FileMode // 目录权限
	BaseURL       string      // 基础URL
}

// New 创建本地存储
func New(config Config, log logger.Logger, options ...storage.StorageOption) (*LocalStorage, error) {
	// 使用默认值
	if config.BasePath == "" {
		config.BasePath = "./uploads"
	}
	if config.TempDir == "" {
		config.TempDir = filepath.Join(config.BasePath, "temp")
	}
	if config.FileMode == 0 {
		config.FileMode = 0644
	}
	if config.DirMode == 0 {
		config.DirMode = 0755
	}

	// 创建基础存储
	baseStorage := storage.NewBaseStorage(storage.StorageTypeLocal, "local")
	if config.BaseURL != "" {
		baseStorage.SetBaseURL(config.BaseURL)
	}

	// 创建本地存储
	s := &LocalStorage{
		BaseStorage:   baseStorage,
		basePath:      config.BasePath,
		tempDir:       config.TempDir,
		autoCreateDir: config.AutoCreateDir,
		fileMode:      config.FileMode,
		dirMode:       config.DirMode,
		logger:        log,
	}

	// 应用选项
	for _, option := range options {
		option(s)
	}

	// 确保目录存在
	if s.autoCreateDir {
		if err := s.ensureDir(s.basePath); err != nil {
			return nil, fmt.Errorf("failed to create base path: %w", err)
		}
		if err := s.ensureDir(s.tempDir); err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	return s, nil
}

// UploadFile 上传文件
func (s *LocalStorage) UploadFile(file *multipart.FileHeader, subPath string) (*storage.FileInfo, error) {
	// 检查文件大小
	if s.GetMaxFileSize() > 0 && file.Size > s.GetMaxFileSize() {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d > %d", file.Size, s.GetMaxFileSize())
	}

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 生成唯一文件名
	originalName := filepath.Base(file.Filename)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	uniqueName := fmt.Sprintf("%s_%d_%s%s", baseName, time.Now().UnixNano(), s.generateRandomString(8), ext)

	// 确保子目录存在
	uploadDir := s.basePath
	if subPath != "" {
		uploadDir = filepath.Join(s.basePath, subPath)
		if s.autoCreateDir {
			if err := s.ensureDir(uploadDir); err != nil {
				return nil, fmt.Errorf("failed to create upload directory: %w", err)
			}
		}
	}

	// 构建文件路径
	filePath := filepath.Join(uploadDir, uniqueName)
	relativePath := uniqueName
	if subPath != "" {
		relativePath = filepath.Join(subPath, uniqueName)
	}

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 计算MD5哈希
	hash := md5.New()
	reader := io.TeeReader(src, hash)

	// 复制文件内容
	size, err := io.Copy(dst, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 设置文件权限
	if err := os.Chmod(filePath, s.fileMode); err != nil {
		s.logger.Warn("Failed to set file permission", "error", err)
	}

	// 获取MIME类型
	mimeType := s.detectMimeType(filePath)

	// 构建文件信息
	fileInfo := &storage.FileInfo{
		OriginalName: originalName,
		FileName:     uniqueName,
		RelativePath: relativePath,
		Size:         size,
		MimeType:     mimeType,
		Hash:         hex.EncodeToString(hash.Sum(nil)),
		URL:          s.GetFileURL(relativePath),
		UploadTime:   time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	return fileInfo, nil
}

// SaveFile 保存文件
func (s *LocalStorage) SaveFile(data []byte, fileName, subPath string) (*storage.FileInfo, error) {
	// 检查文件大小
	if s.GetMaxFileSize() > 0 && int64(len(data)) > s.GetMaxFileSize() {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d > %d", len(data), s.GetMaxFileSize())
	}

	// 生成唯一文件名
	originalName := filepath.Base(fileName)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	uniqueName := fmt.Sprintf("%s_%d_%s%s", baseName, time.Now().UnixNano(), s.generateRandomString(8), ext)

	// 确保子目录存在
	uploadDir := s.basePath
	if subPath != "" {
		uploadDir = filepath.Join(s.basePath, subPath)
		if s.autoCreateDir {
			if err := s.ensureDir(uploadDir); err != nil {
				return nil, fmt.Errorf("failed to create upload directory: %w", err)
			}
		}
	}

	// 构建文件路径
	filePath := filepath.Join(uploadDir, uniqueName)
	relativePath := uniqueName
	if subPath != "" {
		relativePath = filepath.Join(subPath, uniqueName)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, s.fileMode); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 计算MD5哈希
	hash := md5.Sum(data)

	// 获取MIME类型
	mimeType := s.detectMimeType(filePath)

	// 构建文件信息
	fileInfo := &storage.FileInfo{
		OriginalName: originalName,
		FileName:     uniqueName,
		RelativePath: relativePath,
		Size:         int64(len(data)),
		MimeType:     mimeType,
		Hash:         hex.EncodeToString(hash[:]),
		URL:          s.GetFileURL(relativePath),
		UploadTime:   time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	return fileInfo, nil
}

// SaveFileFromReader 从读取器保存文件
func (s *LocalStorage) SaveFileFromReader(reader io.Reader, size int64, fileName, subPath string) (*storage.FileInfo, error) {
	// 检查文件大小
	if s.GetMaxFileSize() > 0 && size > s.GetMaxFileSize() {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d > %d", size, s.GetMaxFileSize())
	}

	// 生成唯一文件名
	originalName := filepath.Base(fileName)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	uniqueName := fmt.Sprintf("%s_%d_%s%s", baseName, time.Now().UnixNano(), s.generateRandomString(8), ext)

	// 确保子目录存在
	uploadDir := s.basePath
	if subPath != "" {
		uploadDir = filepath.Join(s.basePath, subPath)
		if s.autoCreateDir {
			if err := s.ensureDir(uploadDir); err != nil {
				return nil, fmt.Errorf("failed to create upload directory: %w", err)
			}
		}
	}

	// 构建文件路径
	filePath := filepath.Join(uploadDir, uniqueName)
	relativePath := uniqueName
	if subPath != "" {
		relativePath = filepath.Join(subPath, uniqueName)
	}

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 计算MD5哈希
	hash := md5.New()
	teeReader := io.TeeReader(reader, hash)

	// 复制文件内容
	written, err := io.Copy(dst, teeReader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 设置文件权限
	if err := os.Chmod(filePath, s.fileMode); err != nil {
		s.logger.Warn("Failed to set file permission", "error", err)
	}

	// 获取MIME类型
	mimeType := s.detectMimeType(filePath)

	// 构建文件信息
	fileInfo := &storage.FileInfo{
		OriginalName: originalName,
		FileName:     uniqueName,
		RelativePath: relativePath,
		Size:         written,
		MimeType:     mimeType,
		Hash:         hex.EncodeToString(hash.Sum(nil)),
		URL:          s.GetFileURL(relativePath),
		UploadTime:   time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	return fileInfo, nil
}

// GetFile 获取文件
func (s *LocalStorage) GetFile(relativePath string) (*storage.FileInfo, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", relativePath)
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// 获取文件名
	fileName := filepath.Base(relativePath)

	// 获取MIME类型
	mimeType := s.detectMimeType(filePath)

	// 构建文件信息
	fileInfo := &storage.FileInfo{
		OriginalName: fileName,
		FileName:     fileName,
		RelativePath: relativePath,
		Size:         info.Size(),
		MimeType:     mimeType,
		URL:          s.GetFileURL(relativePath),
		UploadTime:   info.ModTime(),
		Metadata:     make(map[string]interface{}),
	}

	return fileInfo, nil
}

// DeleteFile 删除文件
func (s *LocalStorage) DeleteFile(relativePath string) error {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 删除文件
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", relativePath)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListFiles 列出文件
func (s *LocalStorage) ListFiles(subPath string) ([]*storage.FileInfo, error) {
	// 构建目录路径
	dirPath := s.basePath
	if subPath != "" {
		dirPath = filepath.Join(s.basePath, subPath)
	}

	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory not found: %s", subPath)
		}
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}
	defer dir.Close()

	// 读取目录项
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// 构建文件信息
	result := make([]*storage.FileInfo, 0, len(fileInfos))
	for _, info := range fileInfos {
		// 跳过目录
		if info.IsDir() {
			continue
		}

		// 构建相对路径
		fileName := info.Name()
		relativePath := fileName
		if subPath != "" {
			relativePath = filepath.Join(subPath, fileName)
		}

		// 获取MIME类型
		mimeType := s.detectMimeType(filepath.Join(dirPath, fileName))

		// 构建文件信息
		fileInfo := &storage.FileInfo{
			OriginalName: fileName,
			FileName:     fileName,
			RelativePath: relativePath,
			Size:         info.Size(),
			MimeType:     mimeType,
			URL:          s.GetFileURL(relativePath),
			UploadTime:   info.ModTime(),
			Metadata:     make(map[string]interface{}),
		}

		result = append(result, fileInfo)
	}

	return result, nil
}

// FileExists 检查文件是否存在
func (s *LocalStorage) FileExists(relativePath string) bool {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	return err == nil
}

// GetFileSize 获取文件大小
func (s *LocalStorage) GetFileSize(relativePath string) (int64, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 获取文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("file not found: %s", relativePath)
		}
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}

// GetFileContent 获取文件内容
func (s *LocalStorage) GetFileContent(relativePath string) ([]byte, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", relativePath)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// GetFileReader 获取文件读取器
func (s *LocalStorage) GetFileReader(relativePath string) (io.ReadCloser, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, relativePath)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", relativePath)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// GetFileURL 获取文件URL
func (s *LocalStorage) GetFileURL(relativePath string) string {
	// 如果设置了基础URL，则使用基础URL
	if s.GetBaseURL() != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.GetBaseURL(), "/"), relativePath)
	}

	// 否则使用相对路径
	return fmt.Sprintf("/uploads/%s", relativePath)
}

// GetSignedURL 获取签名URL
func (s *LocalStorage) GetSignedURL(relativePath string, expireSeconds int) (string, error) {
	// 本地存储不支持签名URL，直接返回普通URL
	return s.GetFileURL(relativePath), nil
}

// Cleanup 清理过期文件
func (s *LocalStorage) Cleanup(subPath string, maxAge time.Duration) error {
	// 构建目录路径
	dirPath := s.basePath
	if subPath != "" {
		dirPath = filepath.Join(s.basePath, subPath)
	}

	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open directory: %w", err)
	}
	defer dir.Close()

	// 读取目录项
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// 当前时间
	now := time.Now()

	// 删除过期文件
	for _, info := range fileInfos {
		// 跳过目录
		if info.IsDir() {
			continue
		}

		// 检查文件是否过期
		if now.Sub(info.ModTime()) > maxAge {
			// 构建文件路径
			filePath := filepath.Join(dirPath, info.Name())

			// 删除文件
			if err := os.Remove(filePath); err != nil {
				s.logger.Warn("Failed to delete expired file", "error", err, "file", filePath)
			} else {
				s.logger.Info("Deleted expired file", "file", filePath)
			}
		}
	}

	return nil
}

// SaveFileFromReaderWithContext 带上下文从读取器保存文件
func (s *LocalStorage) SaveFileFromReaderWithContext(ctx context.Context, reader io.Reader, size int64, fileName, subPath string) (*storage.FileInfo, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// 继续执行
	}

	return s.SaveFileFromReader(reader, size, fileName, subPath)
}

// ensureDir 确保目录存在
func (s *LocalStorage) ensureDir(dirPath string) error {
	// 检查目录是否存在
	info, err := os.Stat(dirPath)
	if err == nil {
		// 目录已存在，检查是否为目录
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", dirPath)
		}
		return nil
	}

	// 目录不存在，创建目录
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, s.dirMode)
	}

	return err
}

// generateRandomString 生成随机字符串
func (s *LocalStorage) generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(1 * time.Nanosecond) // 确保每次生成的字符不同
	}
	return string(result)
}

// detectMimeType 检测MIME类型
func (s *LocalStorage) detectMimeType(filePath string) string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	// 读取前512字节
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "application/octet-stream"
	}

	// 根据扩展名检测MIME类型
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".zip":
		return "application/zip"
	case ".txt":
		return "text/plain"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	default:
		// 使用简单的检测
		if len(buffer) > 3 && buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
			return "image/jpeg"
		}
		if len(buffer) > 4 && buffer[0] == 0x89 && buffer[1] == 'P' && buffer[2] == 'N' && buffer[3] == 'G' {
			return "image/png"
		}
		if len(buffer) > 4 && buffer[0] == 'G' && buffer[1] == 'I' && buffer[2] == 'F' && buffer[3] == '8' {
			return "image/gif"
		}
		if len(buffer) > 4 && buffer[0] == '%' && buffer[1] == 'P' && buffer[2] == 'D' && buffer[3] == 'F' {
			return "application/pdf"
		}
		return "application/octet-stream"
	}
}

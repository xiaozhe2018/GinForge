package storage

import (
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
)

// LocalStorage 本地文件存储
type LocalStorage struct {
	basePath string
	logger   logger.Logger
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(basePath string, log logger.Logger) *LocalStorage {
	// 确保基础目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Error("failed to create storage directory", err, "path", basePath)
	}

	return &LocalStorage{
		basePath: basePath,
		logger:   log,
	}
}

// UploadFile 上传文件
func (ls *LocalStorage) UploadFile(file *multipart.FileHeader, subPath string) (*FileInfo, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 生成文件名
	fileName := ls.generateFileName(file.Filename)
	fullPath := filepath.Join(ls.basePath, subPath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	written, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 计算文件哈希
	hash, err := ls.calculateFileHash(fullPath)
	if err != nil {
		ls.logger.Warn("failed to calculate file hash", "file", fullPath, "error", err)
	}

	fileInfo := &FileInfo{
		OriginalName: file.Filename,
		FileName:     fileName,
		FilePath:     fullPath,
		RelativePath: filepath.Join(subPath, fileName),
		Size:         written,
		MimeType:     file.Header.Get("Content-Type"),
		Hash:         hash,
		UploadTime:   time.Now(),
	}

	ls.logger.Info("file uploaded", "original_name", file.Filename, "file_name", fileName, "size", written)
	return fileInfo, nil
}

// SaveFile 保存文件
func (ls *LocalStorage) SaveFile(data []byte, fileName, subPath string) (*FileInfo, error) {
	// 生成文件名
	generatedName := ls.generateFileName(fileName)
	fullPath := filepath.Join(ls.basePath, subPath, generatedName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 计算文件哈希
	hash, err := ls.calculateFileHash(fullPath)
	if err != nil {
		ls.logger.Warn("failed to calculate file hash", "file", fullPath, "error", err)
	}

	fileInfo := &FileInfo{
		OriginalName: fileName,
		FileName:     generatedName,
		FilePath:     fullPath,
		RelativePath: filepath.Join(subPath, generatedName),
		Size:         int64(len(data)),
		Hash:         hash,
		UploadTime:   time.Now(),
	}

	ls.logger.Info("file saved", "original_name", fileName, "file_name", generatedName, "size", len(data))
	return fileInfo, nil
}

// GetFile 获取文件
func (ls *LocalStorage) GetFile(relativePath string) (*FileInfo, error) {
	fullPath := filepath.Join(ls.basePath, relativePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// 计算文件哈希
	hash, err := ls.calculateFileHash(fullPath)
	if err != nil {
		ls.logger.Warn("failed to calculate file hash", "file", fullPath, "error", err)
	}

	fileInfo := &FileInfo{
		FileName:     stat.Name(),
		FilePath:     fullPath,
		RelativePath: relativePath,
		Size:         stat.Size(),
		Hash:         hash,
		UploadTime:   stat.ModTime(),
	}

	return fileInfo, nil
}

// DeleteFile 删除文件
func (ls *LocalStorage) DeleteFile(relativePath string) error {
	fullPath := filepath.Join(ls.basePath, relativePath)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	ls.logger.Info("file deleted", "path", relativePath)
	return nil
}

// ListFiles 列出文件
func (ls *LocalStorage) ListFiles(subPath string) ([]*FileInfo, error) {
	dirPath := filepath.Join(ls.basePath, subPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []*FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		relativePath := filepath.Join(subPath, entry.Name())

		stat, err := entry.Info()
		if err != nil {
			ls.logger.Warn("failed to get file info", "file", fullPath, "error", err)
			continue
		}

		// 计算文件哈希
		hash, err := ls.calculateFileHash(fullPath)
		if err != nil {
			ls.logger.Warn("failed to calculate file hash", "file", fullPath, "error", err)
		}

		fileInfo := &FileInfo{
			FileName:     entry.Name(),
			FilePath:     fullPath,
			RelativePath: relativePath,
			Size:         stat.Size(),
			Hash:         hash,
			UploadTime:   stat.ModTime(),
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// FileExists 检查文件是否存在
func (ls *LocalStorage) FileExists(relativePath string) bool {
	fullPath := filepath.Join(ls.basePath, relativePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFileSize 获取文件大小
func (ls *LocalStorage) GetFileSize(relativePath string) (int64, error) {
	fullPath := filepath.Join(ls.basePath, relativePath)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return 0, fmt.Errorf("file not found: %w", err)
	}
	return stat.Size(), nil
}

// generateFileName 生成文件名
func (ls *LocalStorage) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	// 使用时间戳和随机数生成唯一文件名
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d", originalName, timestamp)))
	hashStr := hex.EncodeToString(hash[:])[:8]

	return fmt.Sprintf("%s_%d_%s%s", name, timestamp, hashStr, ext)
}

// calculateFileHash 计算文件哈希
func (ls *LocalStorage) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetBasePath 获取基础路径
func (ls *LocalStorage) GetBasePath() string {
	return ls.basePath
}

// Cleanup 清理过期文件
func (ls *LocalStorage) Cleanup(subPath string, maxAge time.Duration) error {
	dirPath := filepath.Join(ls.basePath, subPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	var deletedCount int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			fullPath := filepath.Join(dirPath, entry.Name())
			if err := os.Remove(fullPath); err != nil {
				ls.logger.Warn("failed to delete expired file", "file", fullPath, "error", err)
				continue
			}
			deletedCount++
		}
	}

	ls.logger.Info("cleanup completed", "deleted_count", deletedCount, "max_age", maxAge)
	return nil
}

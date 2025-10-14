// Package factory 提供存储工厂
package factory

import (
	"fmt"
	"os"
	"time"

	"goweb/pkg/logger"
	"goweb/pkg/storage"
	"goweb/pkg/storage/local"
)

// StorageFactory 存储工厂
type StorageFactory struct {
	logger logger.Logger
}

// New 创建存储工厂
func New(log logger.Logger) *StorageFactory {
	return &StorageFactory{
		logger: log,
	}
}

// CreateStorage 创建存储
func (f *StorageFactory) CreateStorage(storageType storage.StorageType, config map[string]interface{}) (storage.Storage, error) {
	// 根据存储类型创建存储
	switch storageType {
	case storage.StorageTypeLocal:
		return f.createLocalStorage(config)
	case storage.StorageTypeOSS:
		return f.createOSSStorage(config)
	case storage.StorageTypeS3:
		return f.createS3Storage(config)
	case storage.StorageTypeCOS:
		return f.createCOSStorage(config)
	case storage.StorageTypeQiniu:
		return f.createQiniuStorage(config)
	case storage.StorageTypeMinio:
		return f.createMinioStorage(config)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// createLocalStorage 创建本地存储
func (f *StorageFactory) createLocalStorage(config map[string]interface{}) (storage.Storage, error) {
	// 获取本地存储配置
	basePath, _ := config["local_base_path"].(string)
	tempDir, _ := config["local_temp_dir"].(string)
	autoCreateDir, _ := config["local_auto_create_dir"].(bool)
	fileMode, _ := config["local_file_mode"].(int)
	dirMode, _ := config["local_dir_mode"].(int)
	baseURL, _ := config["url_prefix"].(string)

	// 创建本地存储配置
	localConfig := local.Config{
		BasePath:      basePath,
		TempDir:       tempDir,
		AutoCreateDir: autoCreateDir,
		FileMode:      os.FileMode(fileMode),
		DirMode:       os.FileMode(dirMode),
		BaseURL:       baseURL,
	}

	// 创建本地存储
	localStorage, err := local.New(localConfig, f.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create local storage: %w", err)
	}

	// 设置选项
	if maxFileSize, ok := config["max_file_size"].(int64); ok && maxFileSize > 0 {
		localStorage.SetMaxFileSize(maxFileSize)
	}
	if timeout, ok := config["timeout"].(int); ok && timeout > 0 {
		localStorage.SetTimeout(time.Duration(timeout) * time.Second)
	}

	return localStorage, nil
}

// createOSSStorage 创建OSS存储
func (f *StorageFactory) createOSSStorage(config map[string]interface{}) (storage.Storage, error) {
	// TODO: 实现OSS存储
	return nil, fmt.Errorf("OSS storage not implemented yet")
}

// createS3Storage 创建S3存储
func (f *StorageFactory) createS3Storage(config map[string]interface{}) (storage.Storage, error) {
	// TODO: 实现S3存储
	return nil, fmt.Errorf("S3 storage not implemented yet")
}

// createCOSStorage 创建COS存储
func (f *StorageFactory) createCOSStorage(config map[string]interface{}) (storage.Storage, error) {
	// TODO: 实现COS存储
	return nil, fmt.Errorf("COS storage not implemented yet")
}

// createQiniuStorage 创建七牛云存储
func (f *StorageFactory) createQiniuStorage(config map[string]interface{}) (storage.Storage, error) {
	// TODO: 实现七牛云存储
	return nil, fmt.Errorf("Qiniu storage not implemented yet")
}

// createMinioStorage 创建MinIO存储
func (f *StorageFactory) createMinioStorage(config map[string]interface{}) (storage.Storage, error) {
	// TODO: 实现MinIO存储
	return nil, fmt.Errorf("MinIO storage not implemented yet")
}

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"goweb/pkg/logger"
)

// 新增存储类型常量
const (
	StorageTypeCOS   StorageType = "cos"   // 腾讯云COS
	StorageTypeQiniu StorageType = "qiniu" // 七牛云
)

// EnhancedStorageConfig 增强的存储配置
type EnhancedStorageConfig struct {
	Type          StorageType        `json:"type"`                    // 存储类型
	Local         LocalConfig        `json:"local"`                   // 本地存储配置
	OSS           OSSConfig          `json:"oss"`                     // OSS配置
	S3            S3Config           `json:"s3"`                      // S3配置
	COS           COSConfig          `json:"cos"`                     // COS配置
	Qiniu         QiniuConfig        `json:"qiniu"`                   // 七牛云配置
	MinIO         MinioConfig        `json:"minio"`                   // MinIO配置
	URLPrefix     string             `json:"url_prefix"`              // URL前缀
	Domains       DomainsConfig      `json:"domains"`                 // 域名配置
	UploadLimits  UploadLimitsConfig `json:"upload_limits"`           // 上传限制
	FileTypes     FileTypesConfig    `json:"file_types"`              // 文件类型配置
	Processing    ProcessingConfig   `json:"processing"`              // 文件处理配置
	Security      SecurityConfig     `json:"security"`                // 安全配置
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath      string `json:"base_path"`       // 存储基础路径
	TempDir       string `json:"temp_dir"`        // 临时文件目录
	AutoCreateDir bool   `json:"auto_create_dir"` // 是否自动创建目录
	FileMode      uint32 `json:"file_mode"`       // 文件权限
	DirMode       uint32 `json:"dir_mode"`        // 目录权限
}

// COSConfig 腾讯云COS配置
type COSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
	CustomDomain    string `json:"custom_domain"`
	Prefix          string `json:"prefix"`
	SSL             bool   `json:"ssl"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	BucketName string `json:"bucket_name"`
	Domain     string `json:"domain"`
	Zone       string `json:"zone"`
	UseHTTPS   bool   `json:"use_https"`
	Prefix     string `json:"prefix"`
}

// DomainsConfig 域名配置
type DomainsConfig struct {
	UseCustomDomain bool   `json:"use_custom_domain"` // 是否使用自定义域名
	StaticDomain    string `json:"static_domain"`     // 静态资源域名
	UploadDomain    string `json:"upload_domain"`     // 上传服务域名
	CDNDomain       string `json:"cdn_domain"`        // CDN域名
	ImageDomain     string `json:"image_domain"`      // 图片处理域名
}

// UploadLimitsConfig 上传限制配置
type UploadLimitsConfig struct {
	MaxFileSize               int64          `json:"max_file_size"`                 // 最大文件大小
	MaxFilesPerRequest        int            `json:"max_files_per_request"`         // 单次请求最大文件数
	MaxUploadsPerMinute       int            `json:"max_uploads_per_minute"`        // 每分钟最大上传数量
	MaxDailyUploadSizePerUser int64          `json:"max_daily_upload_size_per_user"` // 每天每用户最大上传容量
	MaxDailyUploadsPerUser    int            `json:"max_daily_uploads_per_user"`    // 每天每用户最大上传数量
	MaxConcurrentUploads      int            `json:"max_concurrent_uploads"`        // 并发上传数量限制
	ChunkedUpload            ChunkedUploadConfig `json:"chunked_upload"`           // 分片上传配置
}

// ChunkedUploadConfig 分片上传配置
type ChunkedUploadConfig struct {
	Enabled      bool   `json:"enabled"`       // 是否启用分片上传
	ChunkSize    int64  `json:"chunk_size"`    // 分片大小
	ChunkTimeout int    `json:"chunk_timeout"` // 分片超时时间
	ChunkDir     string `json:"chunk_dir"`     // 分片存储目录
}

// FileTypesConfig 文件类型配置
type FileTypesConfig struct {
	AllowedMimeTypes     map[string][]string `json:"allowed_mime_types"`      // 允许的MIME类型
	AllowedExtensions    map[string][]string `json:"allowed_extensions"`     // 允许的文件扩展名
	DeniedExtensions     []string            `json:"denied_extensions"`      // 禁止的文件扩展名
	UserTypeRestrictions map[string][]string `json:"user_type_restrictions"` // 按用户类型限制文件类型
}

// ProcessingConfig 文件处理配置
type ProcessingConfig struct {
	ImageProcessing bool                 `json:"image_processing"` // 是否启用图片处理
	ImageService    string               `json:"image_service"`    // 图片处理服务类型
	ImageParams     map[string]interface{} `json:"image_params"`     // 图片处理参数
	Thumbnails      ThumbnailsConfig     `json:"thumbnails"`       // 缩略图配置
}

// ThumbnailsConfig 缩略图配置
type ThumbnailsConfig struct {
	AutoGenerate bool             `json:"auto_generate"` // 是否自动生成缩略图
	Sizes        []ThumbnailSize  `json:"sizes"`         // 缩略图尺寸
}

// ThumbnailSize 缩略图尺寸
type ThumbnailSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Suffix string `json:"suffix"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	RefererCheck     bool     `json:"referer_check"`      // 是否启用防盗链
	AllowedReferers  []string `json:"allowed_referers"`   // 允许的来源域名
	AllowEmptyReferer bool     `json:"allow_empty_referer"` // 是否允许空Referer
	URLExpire        int      `json:"url_expire"`         // 签名URL过期时间
	SignedURL        bool     `json:"signed_url"`         // 是否启用签名URL
}

// LoadEnhancedConfig 从Viper加载增强的存储配置
func LoadEnhancedConfig(cfg *viper.Viper, log logger.Logger) (*EnhancedStorageConfig, error) {
	config := &EnhancedStorageConfig{}
	
	// 基本配置
	config.Type = StorageType(cfg.GetString("storage.type"))
	config.URLPrefix = cfg.GetString("storage.url_prefix")
	
	// 本地存储配置
	config.Local = LocalConfig{
		BasePath:      cfg.GetString("storage.local.base_path"),
		TempDir:       cfg.GetString("storage.local.temp_dir"),
		AutoCreateDir: cfg.GetBool("storage.local.auto_create_dir"),
		FileMode:      uint32(cfg.GetInt("storage.local.file_mode")),
		DirMode:       uint32(cfg.GetInt("storage.local.dir_mode")),
	}
	
	// 如果本地存储目录不存在且启用了自动创建目录，则创建目录
	if config.Type == StorageTypeLocal && config.Local.AutoCreateDir {
		if err := ensureDir(config.Local.BasePath, os.FileMode(config.Local.DirMode)); err != nil {
			log.Warn("Failed to create storage base path", "error", err)
		}
		
		if config.Local.TempDir != "" {
			if err := ensureDir(config.Local.TempDir, os.FileMode(config.Local.DirMode)); err != nil {
				log.Warn("Failed to create temp directory", "error", err)
			}
		}
	}
	
	// 云存储配置
	config.OSS = OSSConfig{
		Endpoint:        cfg.GetString("storage.oss.endpoint"),
		AccessKeyID:     cfg.GetString("storage.oss.access_key_id"),
		AccessKeySecret: cfg.GetString("storage.oss.access_key_secret"),
		BucketName:      cfg.GetString("storage.oss.bucket_name"),
		Region:          cfg.GetString("storage.oss.region"),
	}
	
	config.S3 = S3Config{
		Endpoint:        cfg.GetString("storage.s3.endpoint"),
		AccessKeyID:     cfg.GetString("storage.s3.access_key_id"),
		AccessKeySecret: cfg.GetString("storage.s3.access_key_secret"),
		BucketName:      cfg.GetString("storage.s3.bucket_name"),
		Region:          cfg.GetString("storage.s3.region"),
	}
	
	config.MinIO = MinioConfig{
		Endpoint:        cfg.GetString("storage.minio.endpoint"),
		AccessKeyID:     cfg.GetString("storage.minio.access_key_id"),
		AccessKeySecret: cfg.GetString("storage.minio.access_key_secret"),
		BucketName:      cfg.GetString("storage.minio.bucket_name"),
		UseSSL:          cfg.GetBool("storage.minio.use_ssl"),
	}
	
	// 域名配置
	config.Domains = DomainsConfig{
		UseCustomDomain: cfg.GetBool("storage.domains.use_custom_domain"),
		StaticDomain:    cfg.GetString("storage.domains.static_domain"),
		UploadDomain:    cfg.GetString("storage.domains.upload_domain"),
		CDNDomain:       cfg.GetString("storage.domains.cdn_domain"),
		ImageDomain:     cfg.GetString("storage.domains.image_domain"),
	}
	
	// 上传限制配置
	config.UploadLimits = UploadLimitsConfig{
		MaxFileSize:               cfg.GetInt64("storage.upload_limits.max_file_size"),
		MaxFilesPerRequest:        cfg.GetInt("storage.upload_limits.max_files_per_request"),
		MaxUploadsPerMinute:       cfg.GetInt("storage.upload_limits.max_uploads_per_minute"),
		MaxDailyUploadSizePerUser: cfg.GetInt64("storage.upload_limits.max_daily_upload_size_per_user"),
		MaxDailyUploadsPerUser:    cfg.GetInt("storage.upload_limits.max_daily_uploads_per_user"),
		MaxConcurrentUploads:      cfg.GetInt("storage.upload_limits.max_concurrent_uploads"),
		ChunkedUpload: ChunkedUploadConfig{
			Enabled:      cfg.GetBool("storage.upload_limits.chunked_upload.enabled"),
			ChunkSize:    cfg.GetInt64("storage.upload_limits.chunked_upload.chunk_size"),
			ChunkTimeout: cfg.GetInt("storage.upload_limits.chunked_upload.chunk_timeout"),
			ChunkDir:     cfg.GetString("storage.upload_limits.chunked_upload.chunk_dir"),
		},
	}
	
	// 安全配置
	config.Security = SecurityConfig{
		RefererCheck:      cfg.GetBool("storage.security.referer_check"),
		AllowEmptyReferer: cfg.GetBool("storage.security.allow_empty_referer"),
		URLExpire:         cfg.GetInt("storage.security.url_expire"),
		SignedURL:         cfg.GetBool("storage.security.signed_url"),
	}
	
	// 获取允许的来源域名
	if config.Security.RefererCheck {
		config.Security.AllowedReferers = cfg.GetStringSlice("storage.security.allowed_referers")
	}
	
	// 获取环境变量覆盖
	if os.Getenv("USE_CUSTOM_DOMAIN") == "true" {
		config.Domains.UseCustomDomain = true
	}
	
	if staticDomain := os.Getenv("STATIC_DOMAIN"); staticDomain != "" {
		config.Domains.StaticDomain = staticDomain
	}
	
	if uploadDomain := os.Getenv("UPLOAD_DOMAIN"); uploadDomain != "" {
		config.Domains.UploadDomain = uploadDomain
	}
	
	// 转换为旧配置格式，以便兼容现有代码
	_ = &StorageConfig{
		Type:          config.Type,
		LocalBasePath: config.Local.BasePath,
		OSSConfig:     config.OSS,
		S3Config:      config.S3,
		MinioConfig:   config.MinIO,
		URLPrefix:     config.URLPrefix,
		MaxFileSize:   config.UploadLimits.MaxFileSize,
	}
	
	return config, nil
}

// ensureDir 确保目录存在
func ensureDir(dirPath string, mode os.FileMode) error {
	if dirPath == "" {
		return fmt.Errorf("directory path is empty")
	}
	
	// 如果目录已存在，直接返回
	if _, err := os.Stat(dirPath); err == nil {
		return nil
	}
	
	// 创建目录
	return os.MkdirAll(dirPath, mode)
}

// IsAllowedFileType 检查文件类型是否允许
func (c *EnhancedStorageConfig) IsAllowedFileType(mimeType, extension, userType string) bool {
	// 如果没有配置文件类型限制，则允许所有类型
	if len(c.FileTypes.AllowedMimeTypes) == 0 && len(c.FileTypes.AllowedExtensions) == 0 {
		return true
	}
	
	// 检查扩展名是否在禁止列表中
	extension = strings.ToLower(extension)
	for _, denied := range c.FileTypes.DeniedExtensions {
		if extension == denied {
			return false
		}
	}
	
	// 获取用户类型允许的文件类别
	allowedCategories, ok := c.FileTypes.UserTypeRestrictions[userType]
	if !ok {
		// 如果没有为该用户类型配置限制，则允许所有类型
		return true
	}
	
	// 检查MIME类型是否在允许列表中
	for _, category := range allowedCategories {
		allowedMimes, ok := c.FileTypes.AllowedMimeTypes[category]
		if !ok {
			continue
		}
		
		for _, allowed := range allowedMimes {
			if mimeType == allowed {
				return true
			}
		}
	}
	
	// 检查扩展名是否在允许列表中
	for _, category := range allowedCategories {
		allowedExts, ok := c.FileTypes.AllowedExtensions[category]
		if !ok {
			continue
		}
		
		for _, allowed := range allowedExts {
			if extension == allowed {
				return true
			}
		}
	}
	
	return false
}

// GetFileURL 根据配置生成文件URL
func (c *EnhancedStorageConfig) GetFileURL(subPath, fileName string) string {
	var baseURL string
	
	// 检查是否使用自定义域名
	if c.Domains.UseCustomDomain && c.Domains.StaticDomain != "" {
		baseURL = c.Domains.StaticDomain
	} else if c.URLPrefix != "" {
		baseURL = c.URLPrefix
	} else {
		// 默认URL前缀
		baseURL = "/uploads"
	}
	
	// 构建完整URL
	filePath := filepath.Join(subPath, fileName)
	return fmt.Sprintf("%s/%s", baseURL, filePath)
}

// GetUploadURL 获取上传服务URL
func (c *EnhancedStorageConfig) GetUploadURL() string {
	if c.Domains.UseCustomDomain && c.Domains.UploadDomain != "" {
		return c.Domains.UploadDomain
	}
	return ""
}

// GetCDNURL 获取CDN URL
func (c *EnhancedStorageConfig) GetCDNURL(subPath, fileName string) string {
	if c.Domains.UseCustomDomain && c.Domains.CDNDomain != "" {
		filePath := filepath.Join(subPath, fileName)
		return fmt.Sprintf("%s/%s", c.Domains.CDNDomain, filePath)
	}
	return c.GetFileURL(subPath, fileName)
}

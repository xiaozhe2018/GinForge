package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"goweb/pkg/logger"
)

// SecurityService 安全服务
type SecurityService struct {
	config *EnhancedStorageConfig
	logger logger.Logger
}

// NewSecurityService 创建安全服务
func NewSecurityService(config *EnhancedStorageConfig, log logger.Logger) *SecurityService {
	return &SecurityService{
		config: config,
		logger: log,
	}
}

// CheckReferer 检查Referer
func (s *SecurityService) CheckReferer(r *http.Request) bool {
	// 如果未启用防盗链，则直接通过
	if !s.config.Security.RefererCheck {
		return true
	}

	// 获取Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		// 如果允许空Referer，则通过
		return s.config.Security.AllowEmptyReferer
	}

	// 解析Referer
	u, err := url.Parse(referer)
	if err != nil {
		s.logger.Warn("Failed to parse referer", "referer", referer, "error", err)
		return false
	}

	// 检查Referer是否在允许列表中
	host := u.Host
	for _, allowed := range s.config.Security.AllowedReferers {
		if allowed == host {
			return true
		}

		// 支持通配符匹配
		if strings.HasPrefix(allowed, "*.") && strings.HasSuffix(host, allowed[1:]) {
			return true
		}
	}

	s.logger.Warn("Referer not allowed", "referer", referer, "host", host)
	return false
}

// GenerateSignedURL 生成签名URL
func (s *SecurityService) GenerateSignedURL(path string, expireSeconds int) string {
	// 如果未启用签名URL，则直接返回原始路径
	if !s.config.Security.SignedURL {
		return path
	}

	// 如果未指定过期时间，则使用配置的默认值
	if expireSeconds <= 0 {
		expireSeconds = s.config.Security.URLExpire
	}

	// 计算过期时间戳
	expireTime := time.Now().Unix() + int64(expireSeconds)
	expireStr := strconv.FormatInt(expireTime, 10)

	// 生成签名
	signStr := path + "?e=" + expireStr
	signature := s.generateSignature(signStr)

	// 构建最终URL
	return path + "?e=" + expireStr + "&token=" + url.QueryEscape(signature)
}

// VerifySignedURL 验证签名URL
func (s *SecurityService) VerifySignedURL(r *http.Request) bool {
	// 如果未启用签名URL，则直接通过
	if !s.config.Security.SignedURL {
		return true
	}

	// 获取过期时间和签名
	query := r.URL.Query()
	expireStr := query.Get("e")
	token := query.Get("token")

	// 检查参数是否完整
	if expireStr == "" || token == "" {
		s.logger.Warn("Missing expire time or token")
		return false
	}

	// 检查URL是否过期
	expireTime, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil {
		s.logger.Warn("Invalid expire time", "expire", expireStr, "error", err)
		return false
	}

	if time.Now().Unix() > expireTime {
		s.logger.Warn("URL expired", "expire", expireTime, "now", time.Now().Unix())
		return false
	}

	// 重新构建签名字符串
	path := r.URL.Path
	signStr := path + "?e=" + expireStr

	// 生成签名并比较
	signature := s.generateSignature(signStr)
	if signature != token {
		s.logger.Warn("Invalid signature", "expected", signature, "got", token)
		return false
	}

	return true
}

// generateSignature 生成签名
func (s *SecurityService) generateSignature(data string) string {
	// 使用固定的密钥（实际应用中应该从配置中读取）
	key := []byte("your-secret-key")

	// 计算HMAC-SHA256签名
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// IsAllowedIP 检查IP是否在允许列表中
func (s *SecurityService) IsAllowedIP(ip string) bool {
	// 这里可以实现IP白名单/黑名单检查
	// 目前默认允许所有IP
	return true
}

// IsRateLimited 检查是否超过速率限制
func (s *SecurityService) IsRateLimited(userID uint, ip string) bool {
	// 这里可以实现速率限制检查
	// 目前默认不限制
	return false
}

// SanitizeFileName 清理文件名
func (s *SecurityService) SanitizeFileName(fileName string) string {
	// 移除路径信息
	fileName = filepath.Base(fileName)
	
	// 移除特殊字符
	fileName = strings.ReplaceAll(fileName, "..", "")
	fileName = strings.ReplaceAll(fileName, "/", "")
	fileName = strings.ReplaceAll(fileName, "\\", "")
	
	return fileName
}

// CheckFileType 检查文件类型是否允许
func (s *SecurityService) CheckFileType(fileName, mimeType, userType string) error {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(fileName))
	
	// 检查扩展名是否在禁止列表中
	for _, denied := range s.config.FileTypes.DeniedExtensions {
		if ext == denied {
			return fmt.Errorf("file extension not allowed: %s", ext)
		}
	}
	
	// 如果没有配置文件类型限制，则允许所有类型
	if len(s.config.FileTypes.AllowedMimeTypes) == 0 && len(s.config.FileTypes.AllowedExtensions) == 0 {
		return nil
	}
	
	// 检查用户类型是否有限制
	allowedCategories, ok := s.config.FileTypes.UserTypeRestrictions[userType]
	if !ok {
		// 如果没有为该用户类型配置限制，则允许所有类型
		return nil
	}
	
	// 检查MIME类型是否在允许列表中
	for _, category := range allowedCategories {
		allowedMimes, ok := s.config.FileTypes.AllowedMimeTypes[category]
		if !ok {
			continue
		}
		
		for _, allowed := range allowedMimes {
			if mimeType == allowed {
				return nil
			}
		}
	}
	
	// 检查扩展名是否在允许列表中
	for _, category := range allowedCategories {
		allowedExts, ok := s.config.FileTypes.AllowedExtensions[category]
		if !ok {
			continue
		}
		
		for _, allowed := range allowedExts {
			if ext == allowed {
				return nil
			}
		}
	}
	
	return fmt.Errorf("file type not allowed for user type %s: %s (%s)", userType, ext, mimeType)
}

// CheckUploadLimit 检查上传限制
func (s *SecurityService) CheckUploadLimit(userID uint, fileSize int64, uploadCount int) error {
	// 检查文件大小
	if s.config.UploadLimits.MaxFileSize > 0 && fileSize > s.config.UploadLimits.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size: %d > %d", fileSize, s.config.UploadLimits.MaxFileSize)
	}
	
	// 检查单次请求文件数量
	if s.config.UploadLimits.MaxFilesPerRequest > 0 && uploadCount > s.config.UploadLimits.MaxFilesPerRequest {
		return fmt.Errorf("too many files in request: %d > %d", uploadCount, s.config.UploadLimits.MaxFilesPerRequest)
	}
	
	// 这里可以添加更多的限制检查，例如：
	// - 每分钟上传数量
	// - 每天每用户上传容量
	// - 每天每用户上传数量
	// 这些检查需要与数据库交互，记录和查询用户的上传历史
	
	return nil
}

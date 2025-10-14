// Package security 提供文件安全验证和限制接口
package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FileValidator 文件验证接口
type FileValidator interface {
	// Validate 验证文件
	Validate(ctx context.Context, fileName, mimeType string, fileSize int64) error
	
	// ValidateType 验证文件类型
	ValidateType(ctx context.Context, fileName, mimeType string) error
	
	// ValidateSize 验证文件大小
	ValidateSize(ctx context.Context, fileSize int64) error
	
	// ValidateExtension 验证文件扩展名
	ValidateExtension(ctx context.Context, fileName string) error
}

// UploadLimiter 上传限制接口
type UploadLimiter interface {
	// CheckLimit 检查限制
	CheckLimit(ctx context.Context, userID uint, fileSize int64) error
	
	// RecordUpload 记录上传
	RecordUpload(ctx context.Context, userID uint, fileSize int64) error
	
	// GetUserStats 获取用户统计
	GetUserStats(ctx context.Context, userID uint) (map[string]interface{}, error)
}

// RefererChecker 防盗链检查接口
type RefererChecker interface {
	// CheckReferer 检查Referer
	CheckReferer(r *http.Request) bool
}

// URLSigner URL签名接口
type URLSigner interface {
	// SignURL 签名URL
	SignURL(path string, expireSeconds int) (string, error)
	
	// VerifyURL 验证URL
	VerifyURL(r *http.Request) bool
}

// SecurityService 安全服务
type SecurityService struct {
	validators []FileValidator
	limiters   []UploadLimiter
	refererChecker RefererChecker
	urlSigner   URLSigner
}

// NewSecurityService 创建安全服务
func NewSecurityService() *SecurityService {
	return &SecurityService{
		validators: make([]FileValidator, 0),
		limiters:   make([]UploadLimiter, 0),
	}
}

// AddValidator 添加验证器
func (s *SecurityService) AddValidator(validator FileValidator) {
	s.validators = append(s.validators, validator)
}

// AddLimiter 添加限制器
func (s *SecurityService) AddLimiter(limiter UploadLimiter) {
	s.limiters = append(s.limiters, limiter)
}

// SetRefererChecker 设置防盗链检查器
func (s *SecurityService) SetRefererChecker(checker RefererChecker) {
	s.refererChecker = checker
}

// SetURLSigner 设置URL签名器
func (s *SecurityService) SetURLSigner(signer URLSigner) {
	s.urlSigner = signer
}

// ValidateFile 验证文件
func (s *SecurityService) ValidateFile(ctx context.Context, fileName, mimeType string, fileSize int64) error {
	for _, validator := range s.validators {
		if err := validator.Validate(ctx, fileName, mimeType, fileSize); err != nil {
			return err
		}
	}
	return nil
}

// CheckUploadLimit 检查上传限制
func (s *SecurityService) CheckUploadLimit(ctx context.Context, userID uint, fileSize int64) error {
	for _, limiter := range s.limiters {
		if err := limiter.CheckLimit(ctx, userID, fileSize); err != nil {
			return err
		}
	}
	return nil
}

// RecordUpload 记录上传
func (s *SecurityService) RecordUpload(ctx context.Context, userID uint, fileSize int64) error {
	for _, limiter := range s.limiters {
		if err := limiter.RecordUpload(ctx, userID, fileSize); err != nil {
			return err
		}
	}
	return nil
}

// CheckReferer 检查Referer
func (s *SecurityService) CheckReferer(r *http.Request) bool {
	if s.refererChecker == nil {
		return true
	}
	return s.refererChecker.CheckReferer(r)
}

// SignURL 签名URL
func (s *SecurityService) SignURL(path string, expireSeconds int) (string, error) {
	if s.urlSigner == nil {
		return path, nil
	}
	return s.urlSigner.SignURL(path, expireSeconds)
}

// VerifyURL 验证URL
func (s *SecurityService) VerifyURL(r *http.Request) bool {
	if s.urlSigner == nil {
		return true
	}
	return s.urlSigner.VerifyURL(r)
}

// BasicFileValidator 基本文件验证器
type BasicFileValidator struct {
	maxFileSize   int64
	allowedTypes  map[string]bool
	allowedExts   map[string]bool
	deniedExts    map[string]bool
}

// NewBasicFileValidator 创建基本文件验证器
func NewBasicFileValidator(maxFileSize int64, allowedTypes, allowedExts, deniedExts []string) *BasicFileValidator {
	// 转换为map便于查找
	typesMap := make(map[string]bool)
	for _, t := range allowedTypes {
		typesMap[t] = true
	}
	
	extsMap := make(map[string]bool)
	for _, e := range allowedExts {
		extsMap[e] = true
	}
	
	deniedMap := make(map[string]bool)
	for _, e := range deniedExts {
		deniedMap[e] = true
	}
	
	return &BasicFileValidator{
		maxFileSize:   maxFileSize,
		allowedTypes:  typesMap,
		allowedExts:   extsMap,
		deniedExts:    deniedMap,
	}
}

// Validate 验证文件
func (v *BasicFileValidator) Validate(ctx context.Context, fileName, mimeType string, fileSize int64) error {
	// 验证文件大小
	if err := v.ValidateSize(ctx, fileSize); err != nil {
		return err
	}
	
	// 验证文件类型
	if err := v.ValidateType(ctx, fileName, mimeType); err != nil {
		return err
	}
	
	// 验证文件扩展名
	if err := v.ValidateExtension(ctx, fileName); err != nil {
		return err
	}
	
	return nil
}

// ValidateType 验证文件类型
func (v *BasicFileValidator) ValidateType(ctx context.Context, fileName, mimeType string) error {
	// 如果没有设置允许的类型，则允许所有类型
	if len(v.allowedTypes) == 0 {
		return nil
	}
	
	// 检查MIME类型是否在允许列表中
	if !v.allowedTypes[mimeType] {
		return fmt.Errorf("file type not allowed: %s", mimeType)
	}
	
	return nil
}

// ValidateSize 验证文件大小
func (v *BasicFileValidator) ValidateSize(ctx context.Context, fileSize int64) error {
	// 如果没有设置最大文件大小，则允许所有大小
	if v.maxFileSize <= 0 {
		return nil
	}
	
	// 检查文件大小是否超过限制
	if fileSize > v.maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size: %d > %d", fileSize, v.maxFileSize)
	}
	
	return nil
}

// ValidateExtension 验证文件扩展名
func (v *BasicFileValidator) ValidateExtension(ctx context.Context, fileName string) error {
	ext := filepath.Ext(fileName)
	
	// 检查扩展名是否在禁止列表中
	if v.deniedExts[ext] {
		return fmt.Errorf("file extension not allowed: %s", ext)
	}
	
	// 如果没有设置允许的扩展名，则允许所有扩展名
	if len(v.allowedExts) == 0 {
		return nil
	}
	
	// 检查扩展名是否在允许列表中
	if !v.allowedExts[ext] {
		return fmt.Errorf("file extension not allowed: %s", ext)
	}
	
	return nil
}

// BasicUploadLimiter 基本上传限制器
type BasicUploadLimiter struct {
	maxUploadsPerMinute       int
	maxDailyUploadSizePerUser int64
	maxDailyUploadsPerUser    int
	userUploads               map[uint]*UserUploadStats
	mu                        sync.Mutex
}

// UserUploadStats 用户上传统计
type UserUploadStats struct {
	LastMinuteUploads     int
	LastMinuteResetTime   time.Time
	DailyUploadSize       int64
	DailyUploads          int
	DailyStatsResetTime   time.Time
}

// NewBasicUploadLimiter 创建基本上传限制器
func NewBasicUploadLimiter(maxUploadsPerMinute int, maxDailyUploadSizePerUser int64, maxDailyUploadsPerUser int) *BasicUploadLimiter {
	return &BasicUploadLimiter{
		maxUploadsPerMinute:       maxUploadsPerMinute,
		maxDailyUploadSizePerUser: maxDailyUploadSizePerUser,
		maxDailyUploadsPerUser:    maxDailyUploadsPerUser,
		userUploads:               make(map[uint]*UserUploadStats),
	}
}

// CheckLimit 检查限制
func (l *BasicUploadLimiter) CheckLimit(ctx context.Context, userID uint, fileSize int64) error {
	// 如果用户ID为0，则不检查限制
	if userID == 0 {
		return nil
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// 获取用户上传统计
	stats, ok := l.userUploads[userID]
	if !ok {
		// 如果没有统计信息，则创建新的统计信息
		stats = &UserUploadStats{
			LastMinuteResetTime: time.Now(),
			DailyStatsResetTime: time.Now(),
		}
		l.userUploads[userID] = stats
	}
	
	// 检查是否需要重置分钟统计
	now := time.Now()
	if now.Sub(stats.LastMinuteResetTime) > time.Minute {
		stats.LastMinuteUploads = 0
		stats.LastMinuteResetTime = now
	}
	
	// 检查是否需要重置每日统计
	if now.Day() != stats.DailyStatsResetTime.Day() || now.Month() != stats.DailyStatsResetTime.Month() || now.Year() != stats.DailyStatsResetTime.Year() {
		stats.DailyUploadSize = 0
		stats.DailyUploads = 0
		stats.DailyStatsResetTime = now
	}
	
	// 检查每分钟上传数量
	if l.maxUploadsPerMinute > 0 && stats.LastMinuteUploads >= l.maxUploadsPerMinute {
		return fmt.Errorf("exceeded maximum uploads per minute: %d", l.maxUploadsPerMinute)
	}
	
	// 检查每日上传大小
	if l.maxDailyUploadSizePerUser > 0 && stats.DailyUploadSize+fileSize > l.maxDailyUploadSizePerUser {
		return fmt.Errorf("exceeded maximum daily upload size: %d > %d", stats.DailyUploadSize+fileSize, l.maxDailyUploadSizePerUser)
	}
	
	// 检查每日上传数量
	if l.maxDailyUploadsPerUser > 0 && stats.DailyUploads >= l.maxDailyUploadsPerUser {
		return fmt.Errorf("exceeded maximum daily uploads: %d", l.maxDailyUploadsPerUser)
	}
	
	return nil
}

// RecordUpload 记录上传
func (l *BasicUploadLimiter) RecordUpload(ctx context.Context, userID uint, fileSize int64) error {
	// 如果用户ID为0，则不记录上传
	if userID == 0 {
		return nil
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// 获取用户上传统计
	stats, ok := l.userUploads[userID]
	if !ok {
		// 如果没有统计信息，则创建新的统计信息
		stats = &UserUploadStats{
			LastMinuteResetTime: time.Now(),
			DailyStatsResetTime: time.Now(),
		}
		l.userUploads[userID] = stats
	}
	
	// 检查是否需要重置分钟统计
	now := time.Now()
	if now.Sub(stats.LastMinuteResetTime) > time.Minute {
		stats.LastMinuteUploads = 0
		stats.LastMinuteResetTime = now
	}
	
	// 检查是否需要重置每日统计
	if now.Day() != stats.DailyStatsResetTime.Day() || now.Month() != stats.DailyStatsResetTime.Month() || now.Year() != stats.DailyStatsResetTime.Year() {
		stats.DailyUploadSize = 0
		stats.DailyUploads = 0
		stats.DailyStatsResetTime = now
	}
	
	// 记录上传
	stats.LastMinuteUploads++
	stats.DailyUploadSize += fileSize
	stats.DailyUploads++
	
	return nil
}

// GetUserStats 获取用户统计
func (l *BasicUploadLimiter) GetUserStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	// 如果用户ID为0，则返回空统计
	if userID == 0 {
		return map[string]interface{}{
			"last_minute_uploads": 0,
			"daily_upload_size":   0,
			"daily_uploads":       0,
		}, nil
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// 获取用户上传统计
	stats, ok := l.userUploads[userID]
	if !ok {
		// 如果没有统计信息，则返回空统计
		return map[string]interface{}{
			"last_minute_uploads": 0,
			"daily_upload_size":   0,
			"daily_uploads":       0,
		}, nil
	}
	
	// 检查是否需要重置分钟统计
	now := time.Now()
	if now.Sub(stats.LastMinuteResetTime) > time.Minute {
		stats.LastMinuteUploads = 0
		stats.LastMinuteResetTime = now
	}
	
	// 检查是否需要重置每日统计
	if now.Day() != stats.DailyStatsResetTime.Day() || now.Month() != stats.DailyStatsResetTime.Month() || now.Year() != stats.DailyStatsResetTime.Year() {
		stats.DailyUploadSize = 0
		stats.DailyUploads = 0
		stats.DailyStatsResetTime = now
	}
	
	// 返回统计信息
	return map[string]interface{}{
		"last_minute_uploads": stats.LastMinuteUploads,
		"daily_upload_size":   stats.DailyUploadSize,
		"daily_uploads":       stats.DailyUploads,
	}, nil
}

// SimpleRefererChecker 简单防盗链检查器
type SimpleRefererChecker struct {
	allowedReferers   []string
	allowEmptyReferer bool
}

// NewSimpleRefererChecker 创建简单防盗链检查器
func NewSimpleRefererChecker(allowedReferers []string, allowEmptyReferer bool) *SimpleRefererChecker {
	return &SimpleRefererChecker{
		allowedReferers:   allowedReferers,
		allowEmptyReferer: allowEmptyReferer,
	}
}

// CheckReferer 检查Referer
func (c *SimpleRefererChecker) CheckReferer(r *http.Request) bool {
	// 获取Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		// 如果允许空Referer，则通过
		return c.allowEmptyReferer
	}
	
	// 解析Referer
	u, err := url.Parse(referer)
	if err != nil {
		return false
	}
	
	// 检查Referer是否在允许列表中
	host := u.Host
	for _, allowed := range c.allowedReferers {
		if allowed == host {
			return true
		}
		
		// 支持通配符匹配
		if strings.HasPrefix(allowed, "*.") && strings.HasSuffix(host, allowed[1:]) {
			return true
		}
	}
	
	return false
}

// SimpleURLSigner 简单URL签名器
type SimpleURLSigner struct {
	secretKey string
}

// NewSimpleURLSigner 创建简单URL签名器
func NewSimpleURLSigner(secretKey string) *SimpleURLSigner {
	return &SimpleURLSigner{
		secretKey: secretKey,
	}
}

// SignURL 签名URL
func (s *SimpleURLSigner) SignURL(path string, expireSeconds int) (string, error) {
	// 计算过期时间戳
	expireTime := time.Now().Unix() + int64(expireSeconds)
	expireStr := strconv.FormatInt(expireTime, 10)
	
	// 生成签名
	signStr := path + "?e=" + expireStr
	signature := s.generateSignature(signStr)
	
	// 构建最终URL
	return path + "?e=" + expireStr + "&token=" + url.QueryEscape(signature), nil
}

// VerifyURL 验证URL
func (s *SimpleURLSigner) VerifyURL(r *http.Request) bool {
	// 获取过期时间和签名
	query := r.URL.Query()
	expireStr := query.Get("e")
	token := query.Get("token")
	
	// 检查参数是否完整
	if expireStr == "" || token == "" {
		return false
	}
	
	// 检查URL是否过期
	expireTime, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil {
		return false
	}
	
	if time.Now().Unix() > expireTime {
		return false
	}
	
	// 重新构建签名字符串
	path := r.URL.Path
	signStr := path + "?e=" + expireStr
	
	// 生成签名并比较
	signature := s.generateSignature(signStr)
	return signature == token
}

// generateSignature 生成签名
func (s *SimpleURLSigner) generateSignature(data string) string {
	// 计算HMAC-SHA256签名
	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	
	return signature
}

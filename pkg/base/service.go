package base

import (
	"context"
	"goweb/pkg/logger"

	"go.uber.org/zap"
)

// BaseService 基础服务类
type BaseService struct {
	logger logger.Logger
}

// NewBaseService 创建基础服务
func NewBaseService(logger logger.Logger) *BaseService {
	return &BaseService{
		logger: logger,
	}
}

// SetLogger 设置日志器
func (s *BaseService) SetLogger(logger logger.Logger) {
	s.logger = logger
}

// GetLogger 获取日志器
func (s *BaseService) GetLogger() logger.Logger {
	return s.logger
}

// LogInfo 记录信息日志
func (s *BaseService) LogInfo(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Info(msg, fields...)
	}
}

// LogError 记录错误日志
func (s *BaseService) LogError(msg string, err error, fields ...interface{}) {
	if s.logger != nil {
		allFields := append([]interface{}{"error", err}, fields...)
		s.logger.Error(msg, allFields...)
	}
}

// LogWarn 记录警告日志
func (s *BaseService) LogWarn(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Warn(msg, fields...)
	}
}

// LogDebug 记录调试日志
func (s *BaseService) LogDebug(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Debug(msg, fields...)
	}
}

// WithContext 创建带上下文的日志
func (s *BaseService) WithContext(ctx context.Context) *BaseService {
	if s.logger != nil {
		if traceID := ctx.Value("trace_id"); traceID != nil {
			return &BaseService{
				logger: s.logger.With(zap.String("trace_id", traceID.(string))),
			}
		}
	}
	return s
}

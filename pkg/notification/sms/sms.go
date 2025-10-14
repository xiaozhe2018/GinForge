// Package sms 提供短信发送功能
package sms

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"goweb/pkg/logger"
)

// Config 短信配置
type Config struct {
	Provider    string            `mapstructure:"provider" json:"provider" yaml:"provider"`       // 服务提供商
	APIKey      string            `mapstructure:"api_key" json:"api_key" yaml:"api_key"`         // API密钥
	APISecret   string            `mapstructure:"api_secret" json:"api_secret" yaml:"api_secret"` // API密钥
	APIURL      string            `mapstructure:"api_url" json:"api_url" yaml:"api_url"`         // API地址
	Timeout     int               `mapstructure:"timeout" json:"timeout" yaml:"timeout"`         // 超时时间(秒)
	SignName    string            `mapstructure:"sign_name" json:"sign_name" yaml:"sign_name"`     // 签名
	ExtraParams map[string]string `mapstructure:"extra_params" json:"extra_params" yaml:"extra_params"` // 额外参数
}

// Message 短信消息
type Message struct {
	PhoneNumbers []string          // 手机号码
	TemplateCode string            // 模板代码
	TemplateParam map[string]string // 模板参数
	OutID        string            // 外部流水号
}

// SMSSender 短信发送接口
type SMSSender interface {
	// Send 发送短信
	Send(ctx context.Context, message *Message) error
	
	// SendAsync 异步发送短信
	SendAsync(ctx context.Context, message *Message) <-chan error
	
	// SendBatch 批量发送短信
	SendBatch(ctx context.Context, messages []*Message) []error
}

// Provider 短信服务提供商
type Provider string

// 支持的短信服务提供商
const (
	ProviderAliyun   Provider = "aliyun"   // 阿里云
	ProviderTencent  Provider = "tencent"  // 腾讯云
	ProviderHuawei   Provider = "huawei"   // 华为云
	ProviderNetease  Provider = "netease"  // 网易云
	ProviderTwilio   Provider = "twilio"   // Twilio
	ProviderCustom   Provider = "custom"   // 自定义
)

// BaseSender 基础短信发送器
type BaseSender struct {
	config Config
	logger logger.Logger
	client *http.Client
}

// NewBaseSender 创建基础短信发送器
func NewBaseSender(config Config, log logger.Logger) *BaseSender {
	if config.Timeout == 0 {
		config.Timeout = 10 // 默认10秒
	}
	
	return &BaseSender{
		config: config,
		logger: log,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// AliyunSender 阿里云短信发送器
type AliyunSender struct {
	*BaseSender
}

// NewAliyunSender 创建阿里云短信发送器
func NewAliyunSender(config Config, log logger.Logger) (*AliyunSender, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}
	
	if config.APISecret == "" {
		return nil, errors.New("API secret cannot be empty")
	}
	
	if config.APIURL == "" {
		config.APIURL = "https://dysmsapi.aliyuncs.com"
	}
	
	return &AliyunSender{
		BaseSender: NewBaseSender(config, log),
	}, nil
}

// Send 发送短信
func (s *AliyunSender) Send(ctx context.Context, message *Message) error {
	if len(message.PhoneNumbers) == 0 {
		return errors.New("phone numbers cannot be empty")
	}
	
	if message.TemplateCode == "" {
		return errors.New("template code cannot be empty")
	}
	
	// 构建请求参数
	params := url.Values{}
	params.Set("AccessKeyId", s.config.APIKey)
	params.Set("SignName", s.config.SignName)
	params.Set("TemplateCode", message.TemplateCode)
	params.Set("PhoneNumbers", strings.Join(message.PhoneNumbers, ","))
	
	// TODO: 实现阿里云短信API的签名和请求逻辑
	
	// 记录日志
	s.logger.Info("Sending SMS via Aliyun", 
		"phone", message.PhoneNumbers, 
		"template", message.TemplateCode)
	
	// 模拟发送
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendAsync 异步发送短信
func (s *AliyunSender) SendAsync(ctx context.Context, message *Message) <-chan error {
	resultCh := make(chan error, 1)
	
	go func() {
		err := s.Send(ctx, message)
		resultCh <- err
	}()
	
	return resultCh
}

// SendBatch 批量发送短信
func (s *AliyunSender) SendBatch(ctx context.Context, messages []*Message) []error {
	results := make([]error, len(messages))
	var wg sync.WaitGroup
	
	for i, message := range messages {
		wg.Add(1)
		go func(index int, msg *Message) {
			defer wg.Done()
			results[index] = s.Send(ctx, msg)
		}(i, message)
	}
	
	wg.Wait()
	return results
}

// TencentSender 腾讯云短信发送器
type TencentSender struct {
	*BaseSender
}

// NewTencentSender 创建腾讯云短信发送器
func NewTencentSender(config Config, log logger.Logger) (*TencentSender, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}
	
	if config.APISecret == "" {
		return nil, errors.New("API secret cannot be empty")
	}
	
	if config.APIURL == "" {
		config.APIURL = "https://sms.tencentcloudapi.com"
	}
	
	return &TencentSender{
		BaseSender: NewBaseSender(config, log),
	}, nil
}

// Send 发送短信
func (s *TencentSender) Send(ctx context.Context, message *Message) error {
	if len(message.PhoneNumbers) == 0 {
		return errors.New("phone numbers cannot be empty")
	}
	
	if message.TemplateCode == "" {
		return errors.New("template code cannot be empty")
	}
	
	// TODO: 实现腾讯云短信API的签名和请求逻辑
	
	// 记录日志
	s.logger.Info("Sending SMS via Tencent", 
		"phone", message.PhoneNumbers, 
		"template", message.TemplateCode)
	
	// 模拟发送
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendAsync 异步发送短信
func (s *TencentSender) SendAsync(ctx context.Context, message *Message) <-chan error {
	resultCh := make(chan error, 1)
	
	go func() {
		err := s.Send(ctx, message)
		resultCh <- err
	}()
	
	return resultCh
}

// SendBatch 批量发送短信
func (s *TencentSender) SendBatch(ctx context.Context, messages []*Message) []error {
	results := make([]error, len(messages))
	var wg sync.WaitGroup
	
	for i, message := range messages {
		wg.Add(1)
		go func(index int, msg *Message) {
			defer wg.Done()
			results[index] = s.Send(ctx, msg)
		}(i, message)
	}
	
	wg.Wait()
	return results
}

// MockSender 用于测试的模拟短信发送器
type MockSender struct {
	SentMessages []*Message
	logger       logger.Logger
	shouldFail   bool
}

// NewMockSender 创建模拟短信发送器
func NewMockSender(log logger.Logger, shouldFail bool) *MockSender {
	return &MockSender{
		SentMessages: make([]*Message, 0),
		logger:       log,
		shouldFail:   shouldFail,
	}
}

// Send 模拟发送短信
func (s *MockSender) Send(ctx context.Context, message *Message) error {
	if s.shouldFail {
		return errors.New("mock SMS sending failed")
	}
	
	s.SentMessages = append(s.SentMessages, message)
	s.logger.Info("Mock SMS sent", "to", message.PhoneNumbers, "template", message.TemplateCode)
	return nil
}

// SendAsync 模拟异步发送短信
func (s *MockSender) SendAsync(ctx context.Context, message *Message) <-chan error {
	resultCh := make(chan error, 1)
	
	go func() {
		err := s.Send(ctx, message)
		resultCh <- err
	}()
	
	return resultCh
}

// SendBatch 模拟批量发送短信
func (s *MockSender) SendBatch(ctx context.Context, messages []*Message) []error {
	results := make([]error, len(messages))
	
	for i, message := range messages {
		results[i] = s.Send(ctx, message)
	}
	
	return results
}

// Reset 重置已发送短信列表
func (s *MockSender) Reset() {
	s.SentMessages = make([]*Message, 0)
}

// SetShouldFail 设置是否应该失败
func (s *MockSender) SetShouldFail(shouldFail bool) {
	s.shouldFail = shouldFail
}

// Factory 短信发送器工厂
type Factory struct {
	logger logger.Logger
}

// NewFactory 创建短信发送器工厂
func NewFactory(log logger.Logger) *Factory {
	return &Factory{
		logger: log,
	}
}

// Create 创建短信发送器
func (f *Factory) Create(config Config) (SMSSender, error) {
	switch Provider(config.Provider) {
	case ProviderAliyun:
		return NewAliyunSender(config, f.logger)
	case ProviderTencent:
		return NewTencentSender(config, f.logger)
	case "":
		return nil, errors.New("SMS provider cannot be empty")
	default:
		return nil, fmt.Errorf("unsupported SMS provider: %s", config.Provider)
	}
}

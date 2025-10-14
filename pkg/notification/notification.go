// Package notification 提供统一的通知服务
package notification

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"goweb/pkg/logger"
	"goweb/pkg/notification/email"
	"goweb/pkg/notification/sms"
)

// NotificationType 通知类型
type NotificationType string

// 支持的通知类型
const (
	TypeEmail      NotificationType = "email"      // 邮件
	TypeSMS        NotificationType = "sms"        // 短信
	TypeWebSocket  NotificationType = "websocket"  // WebSocket
	TypePush       NotificationType = "push"       // 移动推送
	TypeWeChat     NotificationType = "wechat"     // 微信
	TypeDingTalk   NotificationType = "dingtalk"   // 钉钉
	TypeFeishu     NotificationType = "feishu"     // 飞书
)

// NotificationLevel 通知级别
type NotificationLevel string

// 通知级别
const (
	LevelInfo     NotificationLevel = "info"     // 信息
	LevelWarning  NotificationLevel = "warning"  // 警告
	LevelError    NotificationLevel = "error"    // 错误
	LevelCritical NotificationLevel = "critical" // 严重
)

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	Type      NotificationType  // 通知类型
	Name      string            // 模板名称
	Content   string            // 模板内容
	Variables []string          // 模板变量
	Metadata  map[string]string // 元数据
}

// NotificationRequest 通知请求
type NotificationRequest struct {
	Type       NotificationType           // 通知类型
	Level      NotificationLevel          // 通知级别
	Recipients []string                   // 接收者
	Template   string                     // 模板名称
	Data       map[string]interface{}     // 模板数据
	Subject    string                     // 主题（邮件）
	Metadata   map[string]string          // 元数据
	Options    map[NotificationType]interface{} // 类型特定选项
}

// NotificationResult 通知结果
type NotificationResult struct {
	Success    bool   // 是否成功
	Error      error  // 错误信息
	RecipientCount int // 接收者数量
	MessageID  string // 消息ID
}

// NotificationService 通知服务接口
type NotificationService interface {
	// Send 发送通知
	Send(ctx context.Context, request *NotificationRequest) (*NotificationResult, error)
	
	// SendAsync 异步发送通知
	SendAsync(ctx context.Context, request *NotificationRequest) <-chan *NotificationResult
	
	// AddTemplate 添加模板
	AddTemplate(template *NotificationTemplate) error
	
	// GetTemplate 获取模板
	GetTemplate(notificationType NotificationType, name string) (*NotificationTemplate, error)
	
	// SupportsType 是否支持指定通知类型
	SupportsType(notificationType NotificationType) bool
}

// Service 通知服务实现
type Service struct {
	logger      logger.Logger
	emailSender email.EmailSender
	smsSender   sms.SMSSender
	
	emailTemplates *email.TemplateManager
	smsTemplates   *sms.TemplateManager
	
	templates map[NotificationType]map[string]*NotificationTemplate
	mutex     sync.RWMutex
}

// NewService 创建通知服务
func NewService(log logger.Logger) *Service {
	return &Service{
		logger:    log,
		templates: make(map[NotificationType]map[string]*NotificationTemplate),
	}
}

// SetEmailSender 设置邮件发送器
func (s *Service) SetEmailSender(sender email.EmailSender, templates *email.TemplateManager) {
	s.emailSender = sender
	s.emailTemplates = templates
}

// SetSMSSender 设置短信发送器
func (s *Service) SetSMSSender(sender sms.SMSSender, templates *sms.TemplateManager) {
	s.smsSender = sender
	s.smsTemplates = templates
}

// AddTemplate 添加模板
func (s *Service) AddTemplate(template *NotificationTemplate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if template.Name == "" {
		return errors.New("template name cannot be empty")
	}
	
	if template.Content == "" {
		return errors.New("template content cannot be empty")
	}
	
	// 确保类型映射存在
	if _, exists := s.templates[template.Type]; !exists {
		s.templates[template.Type] = make(map[string]*NotificationTemplate)
	}
	
	// 存储模板
	s.templates[template.Type][template.Name] = template
	
	// 根据通知类型添加到对应的模板管理器
	switch template.Type {
	case TypeEmail:
		if s.emailTemplates != nil {
			if err := s.emailTemplates.AddTemplate(template.Name, template.Content); err != nil {
				return fmt.Errorf("failed to add email template: %w", err)
			}
		}
	case TypeSMS:
		if s.smsTemplates != nil {
			if err := s.smsTemplates.AddTemplate(template.Name, template.Content); err != nil {
				return fmt.Errorf("failed to add SMS template: %w", err)
			}
		}
	}
	
	return nil
}

// GetTemplate 获取模板
func (s *Service) GetTemplate(notificationType NotificationType, name string) (*NotificationTemplate, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	templates, exists := s.templates[notificationType]
	if !exists {
		return nil, fmt.Errorf("no templates for notification type %s", notificationType)
	}
	
	template, exists := templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found for notification type %s", name, notificationType)
	}
	
	return template, nil
}

// SupportsType 是否支持指定通知类型
func (s *Service) SupportsType(notificationType NotificationType) bool {
	switch notificationType {
	case TypeEmail:
		return s.emailSender != nil
	case TypeSMS:
		return s.smsSender != nil
	default:
		return false
	}
}

// Send 发送通知
func (s *Service) Send(ctx context.Context, request *NotificationRequest) (*NotificationResult, error) {
	if request == nil {
		return nil, errors.New("notification request cannot be nil")
	}
	
	if len(request.Recipients) == 0 {
		return nil, errors.New("recipients cannot be empty")
	}
	
	// 根据通知类型发送
	switch request.Type {
	case TypeEmail:
		return s.sendEmail(ctx, request)
	case TypeSMS:
		return s.sendSMS(ctx, request)
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", request.Type)
	}
}

// SendAsync 异步发送通知
func (s *Service) SendAsync(ctx context.Context, request *NotificationRequest) <-chan *NotificationResult {
	resultCh := make(chan *NotificationResult, 1)
	
	go func() {
		result, err := s.Send(ctx, request)
		if err != nil {
			resultCh <- &NotificationResult{
				Success: false,
				Error:   err,
			}
			return
		}
		
		resultCh <- result
	}()
	
	return resultCh
}

// sendEmail 发送邮件
func (s *Service) sendEmail(ctx context.Context, request *NotificationRequest) (*NotificationResult, error) {
	if s.emailSender == nil {
		return nil, errors.New("email sender not configured")
	}
	
	// 创建邮件消息
	message := &email.Message{
		To:      request.Recipients,
		Subject: request.Subject,
	}
	
	// 如果指定了模板，使用模板
	if request.Template != "" {
		if s.emailTemplates == nil {
			return nil, errors.New("email template manager not configured")
		}
		
		// 渲染模板
		body, err := s.emailTemplates.RenderTemplate(request.Template, request.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to render email template: %w", err)
		}
		
		message.Body = body
		message.ContentType = "text/html; charset=UTF-8"
	} else if content, ok := request.Data["content"].(string); ok {
		// 直接使用内容
		message.Body = content
	} else {
		return nil, errors.New("no content or template provided for email")
	}
	
	// 应用选项
	if options, ok := request.Options[TypeEmail]; ok {
		if emailOptions, ok := options.(map[string]interface{}); ok {
			// 处理附件
			if attachments, ok := emailOptions["attachments"].([]email.Attachment); ok {
				message.Attachments = attachments
			}
			
			// 处理抄送
			if cc, ok := emailOptions["cc"].([]string); ok {
				message.Cc = cc
			}
			
			// 处理密送
			if bcc, ok := emailOptions["bcc"].([]string); ok {
				message.Bcc = bcc
			}
			
			// 处理自定义头部
			if headers, ok := emailOptions["headers"].(map[string]string); ok {
				message.Headers = headers
			}
		}
	}
	
	// 发送邮件
	err := s.emailSender.Send(ctx, message)
	if err != nil {
		return &NotificationResult{
			Success: false,
			Error:   err,
		}, err
	}
	
	return &NotificationResult{
		Success:       true,
		RecipientCount: len(request.Recipients),
	}, nil
}

// sendSMS 发送短信
func (s *Service) sendSMS(ctx context.Context, request *NotificationRequest) (*NotificationResult, error) {
	if s.smsSender == nil {
		return nil, errors.New("SMS sender not configured")
	}
	
	// 创建短信消息
	message := &sms.Message{
		PhoneNumbers: request.Recipients,
	}
	
	// 如果指定了模板
	if request.Template != "" {
		// 使用服务商模板
		message.TemplateCode = request.Template
		
		// 转换模板参数
		templateParam := make(map[string]string)
		for k, v := range request.Data {
			if str, ok := v.(string); ok {
				templateParam[k] = str
			} else {
				templateParam[k] = fmt.Sprintf("%v", v)
			}
		}
		message.TemplateParam = templateParam
	} else {
		// 如果没有指定模板，但提供了内容，则使用默认模板
		if content, ok := request.Data["content"].(string); ok {
			if s.smsTemplates != nil {
				// 添加默认模板
				defaultTemplate := "default"
				if err := s.smsTemplates.AddTemplate(defaultTemplate, content); err != nil {
					return nil, fmt.Errorf("failed to add default SMS template: %w", err)
				}
				message.TemplateCode = defaultTemplate
			} else {
				return nil, errors.New("SMS template manager not configured")
			}
		} else {
			return nil, errors.New("no content or template provided for SMS")
		}
	}
	
	// 应用选项
	if options, ok := request.Options[TypeSMS]; ok {
		if smsOptions, ok := options.(map[string]interface{}); ok {
			// 处理外部流水号
			if outID, ok := smsOptions["out_id"].(string); ok {
				message.OutID = outID
			}
		}
	}
	
	// 发送短信
	err := s.smsSender.Send(ctx, message)
	if err != nil {
		return &NotificationResult{
			Success: false,
			Error:   err,
		}, err
	}
	
	return &NotificationResult{
		Success:       true,
		RecipientCount: len(request.Recipients),
	}, nil
}

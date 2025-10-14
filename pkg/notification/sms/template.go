package sms

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"text/template"

	"goweb/pkg/logger"
)

// TemplateData 模板数据接口
type TemplateData interface{}

// TemplateManager 短信模板管理器
type TemplateManager struct {
	templates    map[string]*template.Template
	logger       logger.Logger
	mutex        sync.RWMutex
}

// NewTemplateManager 创建短信模板管理器
func NewTemplateManager(log logger.Logger) *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*template.Template),
		logger:    log,
	}
}

// AddTemplate 添加模板
func (m *TemplateManager) AddTemplate(name string, content string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}

	m.templates[name] = tmpl
	return nil
}

// GetTemplate 获取模板
func (m *TemplateManager) GetTemplate(name string) (*template.Template, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tmpl, exists := m.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return tmpl, nil
}

// RenderTemplate 渲染模板
func (m *TemplateManager) RenderTemplate(name string, data TemplateData) (string, error) {
	tmpl, err := m.GetTemplate(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TemplateSender 模板短信发送器
type TemplateSender struct {
	sender    SMSSender
	templates *TemplateManager
	logger    logger.Logger
}

// NewTemplateSender 创建模板短信发送器
func NewTemplateSender(sender SMSSender, templates *TemplateManager, log logger.Logger) *TemplateSender {
	return &TemplateSender{
		sender:    sender,
		templates: templates,
		logger:    log,
	}
}

// SendWithTemplate 使用模板发送短信
func (s *TemplateSender) SendWithTemplate(ctx context.Context, templateName string, data TemplateData, phoneNumbers []string) error {
	if len(phoneNumbers) == 0 {
		return errors.New("phone numbers cannot be empty")
	}

	// 渲染模板
	content, err := s.templates.RenderTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// 创建短信消息
	message := &Message{
		PhoneNumbers: phoneNumbers,
		TemplateCode: templateName,
		TemplateParam: map[string]string{
			"content": content,
		},
	}

	// 发送短信
	return s.sender.Send(ctx, message)
}

// SendWithTemplateAsync 异步使用模板发送短信
func (s *TemplateSender) SendWithTemplateAsync(ctx context.Context, templateName string, data TemplateData, phoneNumbers []string) <-chan error {
	resultCh := make(chan error, 1)

	go func() {
		err := s.SendWithTemplate(ctx, templateName, data, phoneNumbers)
		resultCh <- err
	}()

	return resultCh
}

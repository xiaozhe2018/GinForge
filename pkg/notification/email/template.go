package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"goweb/pkg/logger"
)

// TemplateData 模板数据接口
type TemplateData interface{}

// TemplateManager 邮件模板管理器
type TemplateManager struct {
	templates    map[string]*template.Template
	templateDir  string
	logger       logger.Logger
	mutex        sync.RWMutex
	defaultFuncs template.FuncMap
}

// NewTemplateManager 创建邮件模板管理器
func NewTemplateManager(templateDir string, log logger.Logger) (*TemplateManager, error) {
	if templateDir == "" {
		templateDir = "templates/email" // 默认模板目录
	}

	manager := &TemplateManager{
		templates:   make(map[string]*template.Template),
		templateDir: templateDir,
		logger:      log,
		defaultFuncs: template.FuncMap{
			"upper": strings.ToUpper,
			"lower": strings.ToLower,
			"title": strings.Title,
			"trim":  strings.TrimSpace,
		},
	}

	// 加载模板目录
	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		if err := manager.LoadTemplates(); err != nil {
			return nil, err
		}
	} else {
		log.Warn("Template directory does not exist", "dir", templateDir)
	}

	return manager, nil
}

// LoadTemplates 加载所有模板
func (m *TemplateManager) LoadTemplates() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 清空现有模板
	m.templates = make(map[string]*template.Template)

	// 遍历模板目录
	return filepath.Walk(m.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理HTML和TXT文件
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".html") || strings.HasSuffix(info.Name(), ".txt")) {
			// 读取模板文件
			content, err := ioutil.ReadFile(path)
			if err != nil {
				m.logger.Error("Failed to read template file", "path", path, "error", err)
				return err
			}

			// 解析模板
			tmpl, err := template.New(info.Name()).Funcs(m.defaultFuncs).Parse(string(content))
			if err != nil {
				m.logger.Error("Failed to parse template", "path", path, "error", err)
				return err
			}

			// 获取模板名称（不带扩展名）
			name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			m.templates[name] = tmpl
			m.logger.Debug("Loaded email template", "name", name, "path", path)
		}

		return nil
	})
}

// AddTemplate 添加模板
func (m *TemplateManager) AddTemplate(name string, content string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	tmpl, err := template.New(name).Funcs(m.defaultFuncs).Parse(content)
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

// TemplateSender 模板邮件发送器
type TemplateSender struct {
	sender    EmailSender
	templates *TemplateManager
	logger    logger.Logger
}

// NewTemplateSender 创建模板邮件发送器
func NewTemplateSender(sender EmailSender, templates *TemplateManager, log logger.Logger) *TemplateSender {
	return &TemplateSender{
		sender:    sender,
		templates: templates,
		logger:    log,
	}
}

// SendWithTemplate 使用模板发送邮件
func (s *TemplateSender) SendWithTemplate(ctx context.Context, templateName string, data TemplateData, message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	// 渲染模板
	body, err := s.templates.RenderTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// 设置邮件内容
	message.Body = body

	// 如果没有设置内容类型，默认为HTML
	if message.ContentType == "" {
		message.ContentType = "text/html; charset=UTF-8"
	}

	// 发送邮件
	return s.sender.Send(ctx, message)
}

// SendWithTemplateAsync 异步使用模板发送邮件
func (s *TemplateSender) SendWithTemplateAsync(ctx context.Context, templateName string, data TemplateData, message *Message) <-chan error {
	resultCh := make(chan error, 1)

	go func() {
		err := s.SendWithTemplate(ctx, templateName, data, message)
		resultCh <- err
	}()

	return resultCh
}

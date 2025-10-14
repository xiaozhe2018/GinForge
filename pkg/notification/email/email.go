// Package email 提供邮件发送功能
package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"goweb/pkg/logger"
)

// Config 邮件配置
type Config struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`         // SMTP服务器地址
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`         // SMTP服务器端口
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 用户名/邮箱地址
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码/授权码
	From     string `mapstructure:"from" json:"from" yaml:"from"`         // 发件人
	FromName string `mapstructure:"from_name" json:"from_name" yaml:"from_name"` // 发件人名称
	Timeout  int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`   // 超时时间(秒)
	UseSSL   bool   `mapstructure:"use_ssl" json:"use_ssl" yaml:"use_ssl"`   // 是否使用SSL
}

// Message 邮件消息
type Message struct {
	To          []string          // 收件人
	Cc          []string          // 抄送
	Bcc         []string          // 密送
	Subject     string            // 主题
	Body        string            // 内容
	ContentType string            // 内容类型，默认为text/html
	Attachments []Attachment      // 附件
	Headers     map[string]string // 自定义头部
}

// Attachment 附件
type Attachment struct {
	Filename string // 文件名
	Content  []byte // 文件内容
	MimeType string // MIME类型
}

// EmailSender 邮件发送接口
type EmailSender interface {
	// Send 发送邮件
	Send(ctx context.Context, message *Message) error
	
	// SendAsync 异步发送邮件
	SendAsync(ctx context.Context, message *Message) <-chan error
}

// SMTPSender SMTP邮件发送实现
type SMTPSender struct {
	config Config
	logger logger.Logger
}

// NewSMTPSender 创建SMTP邮件发送器
func NewSMTPSender(config Config, log logger.Logger) (*SMTPSender, error) {
	if config.Host == "" {
		return nil, errors.New("SMTP host cannot be empty")
	}
	
	if config.Port == 0 {
		// 默认端口
		if config.UseSSL {
			config.Port = 465
		} else {
			config.Port = 25
		}
	}
	
	if config.Username == "" {
		return nil, errors.New("SMTP username cannot be empty")
	}
	
	if config.Password == "" {
		return nil, errors.New("SMTP password cannot be empty")
	}
	
	if config.From == "" {
		config.From = config.Username
	}
	
	if config.Timeout == 0 {
		config.Timeout = 10 // 默认10秒
	}
	
	return &SMTPSender{
		config: config,
		logger: log,
	}, nil
}

// Send 发送邮件
func (s *SMTPSender) Send(ctx context.Context, message *Message) error {
	if len(message.To) == 0 {
		return errors.New("recipient cannot be empty")
	}
	
	if message.Subject == "" {
		return errors.New("subject cannot be empty")
	}
	
	if message.Body == "" {
		return errors.New("email body cannot be empty")
	}
	
	// 设置默认内容类型
	contentType := message.ContentType
	if contentType == "" {
		contentType = "text/html; charset=UTF-8"
	}
	
	// 构建邮件头部
	var header strings.Builder
	
	// 发件人
	fromName := s.config.FromName
	if fromName == "" {
		fromName = s.config.From
	}
	header.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, s.config.From))
	
	// 收件人
	header.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	
	// 抄送
	if len(message.Cc) > 0 {
		header.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.Cc, ", ")))
	}
	
	// 主题
	header.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	
	// 日期
	header.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	
	// 内容类型
	header.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	
	// MIME版本
	header.WriteString("MIME-Version: 1.0\r\n")
	
	// 自定义头部
	for k, v := range message.Headers {
		header.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	
	// 如果有附件，需要设置为multipart
	var body string
	if len(message.Attachments) > 0 {
		boundary := fmt.Sprintf("_boundary_%d", time.Now().UnixNano())
		header.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		header.WriteString("\r\n")
		
		// 正文部分
		body += fmt.Sprintf("--%s\r\n", boundary)
		body += fmt.Sprintf("Content-Type: %s\r\n\r\n", contentType)
		body += message.Body
		body += "\r\n"
		
		// 附件部分
		for _, attachment := range message.Attachments {
			body += fmt.Sprintf("--%s\r\n", boundary)
			body += fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", attachment.MimeType, attachment.Filename)
			body += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.Filename)
			body += "Content-Transfer-Encoding: base64\r\n\r\n"
			// TODO: 实现base64编码
			body += "\r\n"
		}
		
		body += fmt.Sprintf("--%s--", boundary)
	} else {
		header.WriteString("\r\n")
		body = message.Body
	}
	
	// 完整邮件内容
	msg := header.String() + body
	
	// 获取所有收件人
	var recipients []string
	recipients = append(recipients, message.To...)
	recipients = append(recipients, message.Cc...)
	recipients = append(recipients, message.Bcc...)
	
	// 发送邮件
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()
	
	// 使用goroutine发送，以便可以通过上下文取消
	errCh := make(chan error, 1)
	go func() {
		var err error
		
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
		
		if s.config.UseSSL {
			// SSL连接
			tlsConfig := &tls.Config{
				ServerName: s.config.Host,
			}
			
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				errCh <- fmt.Errorf("failed to connect to SMTP server: %w", err)
				return
			}
			defer conn.Close()
			
			client, err := smtp.NewClient(conn, s.config.Host)
			if err != nil {
				errCh <- fmt.Errorf("failed to create SMTP client: %w", err)
				return
			}
			defer client.Close()
			
			if err = client.Auth(auth); err != nil {
				errCh <- fmt.Errorf("SMTP authentication failed: %w", err)
				return
			}
			
			if err = client.Mail(s.config.From); err != nil {
				errCh <- fmt.Errorf("failed to set sender: %w", err)
				return
			}
			
			for _, recipient := range recipients {
				if err = client.Rcpt(recipient); err != nil {
					errCh <- fmt.Errorf("failed to add recipient %s: %w", recipient, err)
					return
				}
			}
			
			w, err := client.Data()
			if err != nil {
				errCh <- fmt.Errorf("failed to open data writer: %w", err)
				return
			}
			
			_, err = w.Write([]byte(msg))
			if err != nil {
				errCh <- fmt.Errorf("failed to write email data: %w", err)
				return
			}
			
			err = w.Close()
			if err != nil {
				errCh <- fmt.Errorf("failed to close data writer: %w", err)
				return
			}
			
			err = client.Quit()
			if err != nil {
				errCh <- fmt.Errorf("failed to quit SMTP connection: %w", err)
				return
			}
		} else {
			// 普通SMTP
			err = smtp.SendMail(addr, auth, s.config.From, recipients, []byte(msg))
			if err != nil {
				errCh <- fmt.Errorf("failed to send email: %w", err)
				return
			}
		}
		
		errCh <- nil
	}()
	
	// 等待发送完成或上下文取消
	select {
	case err := <-errCh:
		if err != nil {
			s.logger.Error("Failed to send email", "error", err, "to", message.To)
			return err
		}
		s.logger.Info("Email sent successfully", "to", message.To, "subject", message.Subject)
		return nil
	case <-timeoutCtx.Done():
		return fmt.Errorf("sending email timed out after %d seconds", s.config.Timeout)
	}
}

// SendAsync 异步发送邮件
func (s *SMTPSender) SendAsync(ctx context.Context, message *Message) <-chan error {
	resultCh := make(chan error, 1)
	
	go func() {
		err := s.Send(ctx, message)
		resultCh <- err
	}()
	
	return resultCh
}

// MockSender 用于测试的模拟邮件发送器
type MockSender struct {
	SentMessages []*Message
	logger       logger.Logger
	shouldFail   bool
}

// NewMockSender 创建模拟邮件发送器
func NewMockSender(log logger.Logger, shouldFail bool) *MockSender {
	return &MockSender{
		SentMessages: make([]*Message, 0),
		logger:       log,
		shouldFail:   shouldFail,
	}
}

// Send 模拟发送邮件
func (s *MockSender) Send(ctx context.Context, message *Message) error {
	if s.shouldFail {
		return errors.New("mock email sending failed")
	}
	
	s.SentMessages = append(s.SentMessages, message)
	s.logger.Info("Mock email sent", "to", message.To, "subject", message.Subject)
	return nil
}

// SendAsync 模拟异步发送邮件
func (s *MockSender) SendAsync(ctx context.Context, message *Message) <-chan error {
	resultCh := make(chan error, 1)
	
	go func() {
		err := s.Send(ctx, message)
		resultCh <- err
	}()
	
	return resultCh
}

// Reset 重置已发送邮件列表
func (s *MockSender) Reset() {
	s.SentMessages = make([]*Message, 0)
}

// SetShouldFail 设置是否应该失败
func (s *MockSender) SetShouldFail(shouldFail bool) {
	s.shouldFail = shouldFail
}

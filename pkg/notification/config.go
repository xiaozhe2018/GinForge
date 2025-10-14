package notification

import (
	"goweb/pkg/logger"
	"goweb/pkg/notification/email"
	"goweb/pkg/notification/sms"
)

// Config 通知服务配置
type Config struct {
	Email *email.Config `mapstructure:"email" json:"email" yaml:"email"` // 邮件配置
	SMS   *sms.Config   `mapstructure:"sms" json:"sms" yaml:"sms"`     // 短信配置
}

// NewServiceFromConfig 从配置创建通知服务
func NewServiceFromConfig(config *Config, log logger.Logger) (*Service, error) {
	service := NewService(log)
	
	// 配置邮件服务
	if config.Email != nil && config.Email.Host != "" {
		// 创建邮件模板管理器
		emailTemplates, err := email.NewTemplateManager("templates/email", log)
		if err != nil {
			log.Warn("Failed to create email template manager", "error", err)
		}
		
		// 创建邮件发送器
		emailSender, err := email.NewSMTPSender(*config.Email, log)
		if err != nil {
			log.Warn("Failed to create email sender", "error", err)
		} else {
			service.SetEmailSender(emailSender, emailTemplates)
			log.Info("Email notification service configured", "host", config.Email.Host)
		}
	}
	
	// 配置短信服务
	if config.SMS != nil && config.SMS.Provider != "" {
		// 创建短信模板管理器
		smsTemplates := sms.NewTemplateManager(log)
		
		// 创建短信发送器工厂
		factory := sms.NewFactory(log)
		
		// 创建短信发送器
		smsSender, err := factory.Create(*config.SMS)
		if err != nil {
			log.Warn("Failed to create SMS sender", "error", err)
		} else {
			service.SetSMSSender(smsSender, smsTemplates)
			log.Info("SMS notification service configured", "provider", config.SMS.Provider)
		}
	}
	
	return service, nil
}

# 通知服务使用指南

## 1. 简介

GinForge框架的通知服务提供了统一的接口，用于发送各种类型的通知，包括邮件和短信。通知服务支持模板功能，可以方便地创建和管理通知模板，提高开发效率。

## 2. 功能特性

- **统一接口**：提供统一的接口发送不同类型的通知
- **多种通知渠道**：支持邮件、短信等多种通知渠道
- **模板支持**：支持通知模板，可以方便地创建和管理模板
- **异步发送**：支持异步发送通知，不阻塞主流程
- **批量发送**：支持批量发送通知，提高效率
- **可扩展性**：易于扩展，支持添加新的通知渠道

## 3. 快速开始

### 3.1 配置通知服务

首先，在配置文件中添加通知服务的配置：

```yaml
# configs/config.yaml
notification:
  email:
    host: smtp.example.com
    port: 587
    username: noreply@example.com
    password: your_password_here
    from: noreply@example.com
    from_name: GinForge系统
    timeout: 10
    use_ssl: false
  
  sms:
    provider: aliyun
    api_key: your_api_key_here
    api_secret: your_api_secret_here
    sign_name: GinForge
```

### 3.2 初始化通知服务

在应用启动时初始化通知服务：

```go
package main

import (
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
)

func main() {
	// 初始化日志
	log := logger.NewLogger()
	
	// 加载配置
	cfg := config.LoadConfig("configs/config.yaml")
	
	// 初始化通知服务
	var notificationConfig notification.Config
	if err := cfg.UnmarshalKey("notification", &notificationConfig); err != nil {
		log.Error("Failed to unmarshal notification config", "error", err)
	}
	
	notificationService, err := notification.NewServiceFromConfig(&notificationConfig, log)
	if err != nil {
		log.Error("Failed to create notification service", "error", err)
	}
	
	// 将通知服务注册到依赖注入容器
	// ...
}
```

### 3.3 发送邮件通知

```go
// 发送简单邮件
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Recipients: []string{"user@example.com"},
	Subject:    "测试邮件",
	Data: map[string]interface{}{
		"content": "这是一封测试邮件",
	},
})

// 使用模板发送邮件
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Recipients: []string{"user@example.com"},
	Subject:    "欢迎加入",
	Template:   "welcome",
	Data: map[string]interface{}{
		"Name":             "张三",
		"Username":         "zhangsan",
		"Email":            "user@example.com",
		"VerificationLink": "https://example.com/verify?token=abc123",
		"Year":             "2025",
	},
})
```

### 3.4 发送短信通知

```go
// 发送短信
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeSMS,
	Recipients: []string{"13800138000"},
	Template:   "verification_code",
	Data: map[string]interface{}{
		"code": "123456",
	},
})
```

### 3.5 异步发送通知

```go
// 异步发送通知
resultCh := notificationService.SendAsync(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Recipients: []string{"user@example.com"},
	Subject:    "异步通知",
	Template:   "notification",
	Data: map[string]interface{}{
		"Name":    "张三",
		"Title":   "系统通知",
		"Message": "这是一条异步通知",
		"Level":   "info",
		"Year":    "2025",
	},
})

// 处理结果
go func() {
	result := <-resultCh
	if result.Error != nil {
		log.Error("Failed to send notification", "error", result.Error)
	} else {
		log.Info("Notification sent successfully", "recipient_count", result.RecipientCount)
	}
}()
```

## 4. 邮件服务

### 4.1 配置邮件服务

```yaml
email:
  # SMTP服务器地址
  host: smtp.example.com
  # SMTP服务器端口
  port: 587
  # 用户名/邮箱地址
  username: noreply@example.com
  # 密码/授权码
  password: your_password_here
  # 发件人
  from: noreply@example.com
  # 发件人名称
  from_name: GinForge系统
  # 超时时间(秒)
  timeout: 10
  # 是否使用SSL
  use_ssl: false
```

### 4.2 创建邮件模板

邮件模板位于 `templates/email` 目录下，使用 HTML 格式编写：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>尊敬的 {{.Name}}，</p>
    <p>{{.Message}}</p>
</body>
</html>
```

### 4.3 发送带附件的邮件

```go
// 发送带附件的邮件
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Recipients: []string{"user@example.com"},
	Subject:    "带附件的邮件",
	Template:   "notification",
	Data: map[string]interface{}{
		"Name":    "张三",
		"Title":   "带附件的邮件",
		"Message": "这是一封带附件的邮件",
		"Year":    "2025",
	},
	Options: map[notification.NotificationType]interface{}{
		notification.TypeEmail: map[string]interface{}{
			"attachments": []email.Attachment{
				{
					Filename: "document.pdf",
					Content:  pdfContent,
					MimeType: "application/pdf",
				},
			},
		},
	},
})
```

## 5. 短信服务

### 5.1 配置短信服务

```yaml
sms:
  # 服务提供商: aliyun, tencent
  provider: aliyun
  # API密钥
  api_key: your_api_key_here
  # API密钥
  api_secret: your_api_secret_here
  # API地址 (可选，默认使用服务商标准地址)
  api_url: ""
  # 超时时间(秒)
  timeout: 5
  # 签名
  sign_name: GinForge
```

### 5.2 发送验证码短信

```go
// 发送验证码短信
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeSMS,
	Recipients: []string{"13800138000"},
	Template:   "SMS_123456789", // 服务商的模板ID
	Data: map[string]interface{}{
		"code": "123456",
	},
})
```

## 6. 高级用法

### 6.1 自定义通知模板

```go
// 添加自定义模板
err := notificationService.AddTemplate(&notification.NotificationTemplate{
	Type:    notification.TypeEmail,
	Name:    "custom_template",
	Content: "<html><body><h1>{{.Title}}</h1><p>{{.Content}}</p></body></html>",
	Variables: []string{"Title", "Content"},
})

// 使用自定义模板发送通知
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Recipients: []string{"user@example.com"},
	Subject:    "自定义模板",
	Template:   "custom_template",
	Data: map[string]interface{}{
		"Title":   "自定义模板示例",
		"Content": "这是使用自定义模板发送的邮件",
	},
})
```

### 6.2 设置通知级别

```go
// 发送不同级别的通知
result, err := notificationService.Send(ctx, &notification.NotificationRequest{
	Type:       notification.TypeEmail,
	Level:      notification.LevelWarning, // info, warning, error, critical
	Recipients: []string{"user@example.com"},
	Subject:    "警告通知",
	Template:   "notification",
	Data: map[string]interface{}{
		"Name":    "张三",
		"Title":   "警告通知",
		"Message": "系统检测到异常登录",
		"Level":   "warning",
		"Year":    "2025",
	},
})
```

### 6.3 批量发送通知

```go
// 批量发送通知
requests := []*notification.NotificationRequest{
	{
		Type:       notification.TypeEmail,
		Recipients: []string{"user1@example.com"},
		Subject:    "通知1",
		Template:   "notification",
		Data: map[string]interface{}{
			"Name":    "用户1",
			"Title":   "通知1",
			"Message": "这是通知1",
			"Year":    "2025",
		},
	},
	{
		Type:       notification.TypeEmail,
		Recipients: []string{"user2@example.com"},
		Subject:    "通知2",
		Template:   "notification",
		Data: map[string]interface{}{
			"Name":    "用户2",
			"Title":   "通知2",
			"Message": "这是通知2",
			"Year":    "2025",
		},
	},
}

// 并发发送
var wg sync.WaitGroup
for _, req := range requests {
	wg.Add(1)
	go func(request *notification.NotificationRequest) {
		defer wg.Done()
		resultCh := notificationService.SendAsync(ctx, request)
		result := <-resultCh
		if result.Error != nil {
			log.Error("Failed to send notification", "error", result.Error)
		}
	}(req)
}
wg.Wait()
```

## 7. 最佳实践

### 7.1 通知服务封装

为了更好地管理通知逻辑，建议创建一个通知服务封装：

```go
// notification_service.go
package service

import (
	"context"
	"time"

	"goweb/pkg/logger"
	"goweb/pkg/notification"
)

// UserNotificationService 用户通知服务
type UserNotificationService struct {
	notificationService notification.NotificationService
	logger              logger.Logger
}

// NewUserNotificationService 创建用户通知服务
func NewUserNotificationService(notificationService notification.NotificationService, log logger.Logger) *UserNotificationService {
	return &UserNotificationService{
		notificationService: notificationService,
		logger:              log,
	}
}

// SendWelcomeEmail 发送欢迎邮件
func (s *UserNotificationService) SendWelcomeEmail(ctx context.Context, user *User) error {
	req := &notification.NotificationRequest{
		Type:       notification.TypeEmail,
		Recipients: []string{user.Email},
		Subject:    "欢迎加入 GinForge",
		Template:   "welcome",
		Data: map[string]interface{}{
			"Name":             user.Name,
			"Username":         user.Username,
			"Email":            user.Email,
			"VerificationLink": "https://example.com/verify?token=" + user.VerificationToken,
			"Year":             time.Now().Format("2006"),
		},
	}

	_, err := s.notificationService.Send(ctx, req)
	if err != nil {
		s.logger.Error("Failed to send welcome email", "error", err, "user_id", user.ID)
		return err
	}

	s.logger.Info("Welcome email sent", "user_id", user.ID)
	return nil
}

// SendVerificationCode 发送验证码短信
func (s *UserNotificationService) SendVerificationCode(ctx context.Context, phone string, code string) error {
	req := &notification.NotificationRequest{
		Type:       notification.TypeSMS,
		Recipients: []string{phone},
		Template:   "SMS_123456789", // 服务商的模板ID
		Data: map[string]interface{}{
			"code": code,
		},
	}

	_, err := s.notificationService.Send(ctx, req)
	if err != nil {
		s.logger.Error("Failed to send verification code", "error", err, "phone", phone)
		return err
	}

	s.logger.Info("Verification code sent", "phone", phone)
	return nil
}
```

### 7.2 使用依赖注入

在实际应用中，建议使用依赖注入来管理通知服务：

```go
// 在依赖注入容器中注册通知服务
container.Register(func(cfg *config.Config, log logger.Logger) (notification.NotificationService, error) {
	var notificationConfig notification.Config
	if err := cfg.UnmarshalKey("notification", &notificationConfig); err != nil {
		return nil, err
	}
	
	return notification.NewServiceFromConfig(&notificationConfig, log)
})

// 在服务中使用通知服务
type UserService struct {
	notificationService notification.NotificationService
}

func NewUserService(notificationService notification.NotificationService) *UserService {
	return &UserService{
		notificationService: notificationService,
	}
}
```

### 7.3 异常处理

在发送通知时，应该妥善处理异常：

```go
// 发送通知并处理异常
func sendNotificationWithRetry(ctx context.Context, service notification.NotificationService, req *notification.NotificationRequest, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = service.Send(ctx, req)
		if err == nil {
			return nil
		}
		
		// 等待一段时间后重试
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return err
}
```

## 8. 常见问题

### 8.1 邮件发送失败

可能原因：
- SMTP服务器配置错误
- 网络连接问题
- 邮箱地址格式错误

解决方案：
- 检查SMTP服务器配置
- 确认网络连接正常
- 验证邮箱地址格式

### 8.2 短信发送失败

可能原因：
- API密钥配置错误
- 短信模板未审核通过
- 手机号格式错误

解决方案：
- 检查API密钥配置
- 确认短信模板已审核通过
- 验证手机号格式

### 8.3 模板渲染失败

可能原因：
- 模板语法错误
- 模板文件不存在
- 模板数据不匹配

解决方案：
- 检查模板语法
- 确认模板文件存在
- 验证模板数据是否匹配模板变量

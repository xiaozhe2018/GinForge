# 存储配置指南

## 1. 概述

GinForge框架的文件服务支持多种存储方式，包括本地存储和云存储。本文档将帮助你配置和使用不同的存储选项。

## 2. 基本配置

### 2.1 存储类型

文件服务支持以下存储类型：

- `local`: 本地文件系统存储
- `oss`: 阿里云对象存储服务
- `s3`: AWS S3对象存储
- `minio`: MinIO对象存储
- `cos`: 腾讯云对象存储

配置示例：

```yaml
storage:
  # 存储类型: local, oss, s3, minio, cos
  type: local
```

### 2.2 URL配置

```yaml
storage:
  # URL前缀（用于生成访问URL）
  url_prefix: http://localhost:8086/uploads
  
  # 是否使用自定义域名 (生产环境建议开启)
  use_custom_domain: false
  # 静态资源域名
  static_domain: "https://static.example.com"
  # 上传服务域名
  upload_service_domain: "https://upload.example.com"
```

## 3. 本地存储配置

本地存储是最简单的配置，适合开发环境或小型应用：

```yaml
storage:
  type: local
  local:
    # 存储基础路径
    base_path: ./uploads
    # 临时文件存储路径
    temp_path: ./temp_uploads
    # 自动创建目录
    auto_create_dir: true
    # 文件权限
    file_mode: 0644
    # 目录权限
    dir_mode: 0755
```

## 4. 云存储配置

### 4.1 阿里云OSS

```yaml
storage:
  type: oss
  oss:
    # OSS区域节点
    endpoint: oss-cn-hangzhou.aliyuncs.com
    # 访问密钥ID
    access_key_id: YOUR_ACCESS_KEY_ID
    # 访问密钥密码
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    # 存储桶名称
    bucket_name: your-bucket-name
    # 区域
    region: cn-hangzhou
```

### 4.2 AWS S3

```yaml
storage:
  type: s3
  s3:
    # S3区域节点
    endpoint: s3.amazonaws.com
    # 访问密钥ID
    access_key_id: YOUR_ACCESS_KEY_ID
    # 访问密钥密码
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    # 存储桶名称
    bucket_name: your-bucket-name
    # 区域
    region: us-west-1
```

### 4.3 MinIO

```yaml
storage:
  type: minio
  minio:
    # MinIO服务地址
    endpoint: play.min.io
    # 访问密钥ID
    access_key_id: YOUR_ACCESS_KEY_ID
    # 访问密钥密码
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    # 存储桶名称
    bucket_name: your-bucket-name
    # 使用SSL
    use_ssl: true
```

## 5. 上传限制配置

为了保护服务器资源，可以设置上传限制：

```yaml
storage:
  # 上传限制配置
  upload_limits:
    # 最大文件大小（字节）- 默认100MB
    max_file_size: 104857600
    # 单次请求最大文件数
    max_files_per_request: 10
    # 每分钟最大上传数量
    max_uploads_per_minute: 30
    # 每天每用户最大上传容量 (字节)
    max_upload_capacity_per_day: 1073741824
    # 每天每用户最大上传数量
    max_upload_count_per_day: 100
```

## 6. 文件类型限制

### 6.1 允许的文件类型

```yaml
storage:
  # 文件类型和扩展名限制
  allowed_types:
    # 图片类型
    image:
      mimes: ["image/jpeg", "image/png", "image/gif", "image/webp"]
      exts: [".jpg", ".jpeg", ".png", ".gif", ".webp"]
    # 文档类型
    document:
      mimes: ["application/pdf", "application/msword"]
      exts: [".pdf", ".doc", ".docx"]
```

### 6.2 禁止的扩展名

```yaml
storage:
  # 禁止的扩展名 (黑名单)
  forbidden_exts: [".php", ".exe", ".sh", ".bat", ".cmd"]
```

## 7. 分片上传配置

对于大文件上传，可以启用分片上传功能：

```yaml
storage:
  # 分片上传配置
  chunk_upload:
    # 是否启用分片上传
    enabled: true
    # 分片大小 (字节)
    chunk_size: 5242880 # 5MB
    # 分片上传会话超时时间
    timeout: 24h
    # 分片存储目录
    chunk_dir: chunks
```

## 8. 安全配置

### 8.1 防盗链

防止其他网站直接引用你的资源：

```yaml
storage:
  # 安全配置
  security:
    # 防盗链
    anti_leech:
      # 是否启用防盗链
      enabled: true
      # 允许的Referer列表
      allowed_referers: ["example.com", "*.example.com"]
      # 是否允许空Referer
      allow_empty_referer: false
```

### 8.2 签名URL

为敏感资源提供临时访问链接：

```yaml
storage:
  # 安全配置
  security:
    # 签名URL
    signed_url:
      # 是否启用签名URL
      enabled: true
      # 签名URL默认过期时间
      expire_time: 30m
```

## 9. 实际应用场景

### 9.1 开发环境

开发环境通常使用本地存储，配置简单：

```yaml
storage:
  type: local
  local:
    base_path: ./uploads
  url_prefix: http://localhost:8086/uploads
  upload_limits:
    max_file_size: 104857600 # 100MB
```

### 9.2 生产环境

生产环境建议使用云存储和自定义域名：

```yaml
storage:
  type: oss # 或其他云存储
  oss:
    endpoint: oss-cn-hangzhou.aliyuncs.com
    access_key_id: YOUR_ACCESS_KEY_ID
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    bucket_name: your-bucket-name
    region: cn-hangzhou
  
  use_custom_domain: true
  static_domain: "https://static.example.com"
  
  security:
    anti_leech:
      enabled: true
      allowed_referers: ["example.com", "*.example.com"]
```

### 9.3 高安全需求

对于高安全性需求的应用：

```yaml
storage:
  # 严格的文件类型限制
  allowed_types:
    document:
      mimes: ["application/pdf"]
      exts: [".pdf"]
  
  forbidden_exts: [".php", ".exe", ".sh", ".bat", ".cmd"]
  
  security:
    signed_url:
      enabled: true
      expire_time: 5m # 短期有效
```

### 9.4 大文件处理

对于需要处理大文件的应用：

```yaml
storage:
  upload_limits:
    max_file_size: 5368709120 # 5GB
  
  chunk_upload:
    enabled: true
    chunk_size: 10485760 # 10MB
    timeout: 48h
```

## 10. 常见问题

### 10.1 无法访问上传的文件

- 检查 `url_prefix` 或 `static_domain` 配置是否正确
- 确保文件权限正确
- 检查防盗链设置是否过于严格

### 10.2 上传失败

- 检查文件大小是否超过限制
- 检查文件类型是否被允许
- 检查存储路径权限

### 10.3 分片上传问题

- 检查 `chunk_upload.enabled` 是否为 `true`
- 确保临时目录有足够的权限和空间
- 检查分片超时设置是否合理

## 11. 最佳实践

1. **生产环境使用云存储**：更可靠，更易扩展
2. **启用自定义域名**：便于迁移存储提供商
3. **设置合理的文件类型限制**：防止恶意文件上传
4. **配置防盗链**：保护资源不被盗用
5. **对敏感文件使用签名URL**：增加访问控制
6. **定期清理临时文件**：防止磁盘空间占用过多
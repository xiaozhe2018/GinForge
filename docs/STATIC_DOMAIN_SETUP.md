# 静态资源域名配置指南

## 1. 概述

GinForge框架支持为文件服务配置专用域名，可以将静态资源和上传服务分离到不同的域名下，提高性能和安全性。本文档将指导你如何配置和使用静态资源专用域名。

## 2. 为什么需要静态资源专用域名

使用专用域名有以下优势：

1. **性能优化**：浏览器对同一域名的并发请求有限制，使用专用域名可以增加并发下载
2. **CDN加速**：静态资源域名可以轻松配置CDN加速
3. **安全隔离**：将静态资源与应用逻辑分离，降低安全风险
4. **缓存控制**：为静态资源设置更长的缓存时间
5. **跨域控制**：更精确地控制跨域资源访问

## 3. 配置步骤

### 3.1 域名准备

准备两个域名：

- **静态资源域名**：如 `static.example.com`，用于访问上传的文件
- **上传服务域名**：如 `upload.example.com`，用于文件上传API

### 3.2 配置文件设置

修改 `configs/file-api.yaml`：

```yaml
storage:
  # 是否使用自定义域名 (生产环境建议开启)
  use_custom_domain: true
  # 静态资源域名 (例如: https://static.example.com)
  static_domain: "https://static.example.com"
  # 上传服务域名 (例如: https://upload.example.com)
  upload_service_domain: "https://upload.example.com"
```

### 3.3 Nginx配置

#### 静态资源服务配置

创建 `static.conf` 文件：

```nginx
server {
    listen 80;
    server_name static.example.com;

    # 如果使用HTTPS
    # listen 443 ssl;
    # ssl_certificate /path/to/cert.pem;
    # ssl_certificate_key /path/to/key.pem;

    # 静态文件根目录
    root /path/to/GinForge/uploads;

    location / {
        # 允许跨域访问
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Methods' 'GET, OPTIONS';
        
        # 缓存设置
        expires 7d;
        add_header Cache-Control "public, max-age=604800";
        
        # 启用Gzip压缩
        gzip on;
        gzip_types text/plain text/css application/json application/javascript image/svg+xml;
        
        try_files $uri $uri/ =404;
    }
    
    # 错误页面
    error_page 404 /404.html;
    
    # 日志
    access_log /var/log/nginx/static.access.log;
    error_log /var/log/nginx/static.error.log warn;
}
```

#### 上传服务代理配置

创建 `upload.conf` 文件：

```nginx
server {
    listen 80;
    server_name upload.example.com;

    # 如果使用HTTPS
    # listen 443 ssl;
    # ssl_certificate /path/to/cert.pem;
    # ssl_certificate_key /path/to/key.pem;

    location / {
        # 代理到文件服务
        proxy_pass http://localhost:8086;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 允许大文件上传
        client_max_body_size 100M;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # 日志
    access_log /var/log/nginx/upload.access.log;
    error_log /var/log/nginx/upload.error.log warn;
}
```

### 3.4 DNS配置

在你的DNS服务商处添加以下记录：

1. 为 `static.example.com` 添加A记录，指向你的服务器IP
2. 为 `upload.example.com` 添加A记录，指向你的服务器IP

## 4. 使用方法

### 4.1 上传文件

上传文件时，使用上传服务域名：

```bash
curl -X POST https://upload.example.com/api/v1/files/upload \
  -F "file=@/path/to/your/file.jpg"
```

### 4.2 访问文件

文件上传成功后，返回的URL将使用静态资源域名：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "url": "https://static.example.com/images/example_1642342342_a1b2c3d4.jpg",
    // 其他字段...
  }
}
```

## 5. 本地开发配置

在本地开发环境中，可以通过修改hosts文件来模拟域名：

### 5.1 修改hosts文件

```
127.0.0.1 static.example.local
127.0.0.1 upload.example.local
```

### 5.2 配置文件设置

```yaml
storage:
  use_custom_domain: true
  static_domain: "http://static.example.local"
  upload_service_domain: "http://upload.example.local"
```

### 5.3 本地Nginx配置

与生产环境配置类似，只需将域名改为本地域名。

## 6. CDN集成

### 6.1 配置CDN

1. 在CDN服务商处添加加速域名 `static.example.com`
2. 设置源站为你的服务器IP或域名
3. 配置缓存规则，例如：
   - 图片、视频等静态资源缓存7天
   - 其他文件缓存1天

### 6.2 更新配置

如果使用CDN，需要更新配置文件：

```yaml
storage:
  use_custom_domain: true
  static_domain: "https://static-cdn.example.com" # CDN域名
  upload_service_domain: "https://upload.example.com"
```

## 7. 安全配置

### 7.1 防盗链设置

在Nginx中配置防盗链：

```nginx
location / {
    # 允许的来源域名
    valid_referers none blocked server_names
                   *.example.com example.com;
    
    # 如果来源不合法，返回403
    if ($invalid_referer) {
        return 403;
    }
    
    # 其他配置...
}
```

### 7.2 HTTPS配置

强烈建议为两个域名都配置HTTPS：

```nginx
server {
    listen 443 ssl;
    server_name static.example.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 其他配置...
}
```

## 8. 故障排查

### 8.1 无法访问静态资源

1. 检查DNS解析是否正确：`nslookup static.example.com`
2. 检查Nginx配置是否正确：`nginx -t`
3. 检查文件权限：确保Nginx用户可以读取上传目录
4. 检查防火墙设置：确保80/443端口开放

### 8.2 上传失败

1. 检查上传服务是否正常运行：`curl http://upload.example.com/healthz`
2. 检查Nginx日志：`tail -f /var/log/nginx/upload.error.log`
3. 检查文件服务日志：`tail -f /path/to/GinForge/logs/file-api.log`

## 9. 最佳实践

1. **使用HTTPS**：为所有域名配置HTTPS
2. **配置CDN**：为静态资源域名配置CDN加速
3. **设置合理的缓存**：根据文件类型设置不同的缓存时间
4. **启用Gzip压缩**：减少传输大小
5. **配置防盗链**：防止资源被盗用
6. **监控流量**：定期检查流量和访问日志
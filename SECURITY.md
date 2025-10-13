# 安全政策

## 支持的版本

我们目前为以下版本提供安全更新：

| 版本 | 支持状态 |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## 报告安全漏洞

我们非常重视 GinForge 的安全性。如果你发现了安全漏洞，请通过以下方式负责任地向我们报告。

### 报告流程

**请不要在公开的 Issue 中报告安全漏洞。**

#### 1. 私下报告

请通过以下方式之一私下联系我们：

- **GitHub Security Advisories**（推荐）
  - 访问：https://github.com/xiaozhe2018/GinForge/security/advisories/new
  - 点击 "Report a vulnerability"
  
- **Email**
  - 发送邮件至：security@example.com（替换为实际邮箱）
  - 邮件主题：`[SECURITY] GinForge 安全漏洞报告`

#### 2. 包含的信息

请在报告中包含以下信息：

- **漏洞类型**：如 SQL 注入、XSS、CSRF 等
- **漏洞位置**：受影响的文件、函数或端点
- **复现步骤**：详细的步骤说明如何触发漏洞
- **影响范围**：漏洞可能造成的影响
- **PoC 代码**：如果有的话，提供概念验证代码
- **建议修复**：如果你有修复建议
- **发现者**：你的姓名或昵称（如果希望被致谢）

#### 3. 示例报告

```markdown
## 漏洞描述
在用户登录功能中发现 SQL 注入漏洞。

## 漏洞位置
文件：services/admin-api/internal/repository/admin_repository.go
函数：FindUserByUsername

## 复现步骤
1. 访问登录页面
2. 在用户名输入框输入：`admin' OR '1'='1`
3. 提交登录表单
4. 成功绕过认证

## 影响范围
攻击者可以绕过认证系统，获取未授权访问

## PoC 代码
```bash
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin'\'' OR '\''1'\''='\''1","password":"any"}'
```

## 建议修复
使用参数化查询而非字符串拼接
```

### 响应时间

- **初步响应**：我们会在 **48 小时**内确认收到报告
- **评估时间**：我们会在 **7 天**内评估漏洞并提供初步反馈
- **修复时间**：根据漏洞严重程度，通常在 **14-30 天**内发布修复

### 漏洞等级

我们使用 CVSS 3.1 标准来评估漏洞严重程度：

| 等级 | CVSS 分数 | 响应时间 |
|------|-----------|----------|
| 🔴 严重 | 9.0-10.0 | 24 小时 |
| 🟠 高危 | 7.0-8.9  | 7 天 |
| 🟡 中危 | 4.0-6.9  | 30 天 |
| 🟢 低危 | 0.1-3.9  | 90 天 |

## 安全更新流程

### 1. 漏洞确认
- 安全团队确认漏洞
- 评估影响范围和严重程度
- 制定修复计划

### 2. 修复开发
- 开发修复补丁
- 进行内部测试
- 准备更新文档

### 3. 发布更新
- 发布新版本
- 发布安全公告
- 通知受影响用户

### 4. 公开披露
- 在修复发布后 **30 天**公开漏洞详情
- 在 GitHub Security Advisories 发布公告
- 更新 CHANGELOG.md

## 安全最佳实践

### 部署建议

#### 1. 配置安全
```yaml
# 生产环境必须修改的配置
jwt:
  secret: "使用强随机密钥，不要使用默认值！"
  expire_hours: 24

database:
  password: "使用强密码"
  
redis:
  password: "使用强密码"
```

#### 2. 网络安全
- 使用 HTTPS（TLS 1.2+）
- 配置防火墙规则
- 限制数据库访问
- 使用 VPN 或私有网络

#### 3. 应用安全
```bash
# 设置安全的环境变量
export JWT_SECRET=$(openssl rand -hex 32)
export DB_PASSWORD=$(openssl rand -hex 16)

# 限制文件权限
chmod 600 config.yaml
chmod 600 .env
```

#### 4. 定期更新
```bash
# 更新依赖
go get -u ./...
go mod tidy

# 检查已知漏洞
go list -json -m all | nancy sleuth
```

### 代码审查清单

- [ ] 所有用户输入都经过验证和清理
- [ ] 使用参数化查询防止 SQL 注入
- [ ] 所有敏感数据都已加密
- [ ] 密码使用 bcrypt 等强哈希算法
- [ ] JWT token 有过期时间
- [ ] API 有速率限制
- [ ] 错误信息不泄露敏感信息
- [ ] CORS 配置正确
- [ ] 文件上传有类型和大小限制
- [ ] 使用安全的随机数生成器

## 已知安全问题

### 已修复的漏洞

无（首次发布）

### 当前限制

1. **默认配置**
   - 默认 JWT secret 不安全
   - **解决**：生产环境必须修改

2. **开发模式**
   - Debug 模式会暴露详细错误信息
   - **解决**：生产环境关闭 Debug

## 安全工具

### 推荐的安全扫描工具

```bash
# 静态代码分析
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# 依赖漏洞扫描
go install github.com/sonatype-nexus-community/nancy@latest
go list -json -m all | nancy sleuth

# SQL 注入检测
go install github.com/stripe/safesql@latest
```

### CI/CD 集成

```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Gosec
        uses: securego/gosec@master
```

## 安全联系方式

- **安全邮箱**：security@example.com（替换为实际邮箱）
- **PGP 公钥**：https://example.com/pgp-key.asc（如果使用）
- **GitHub Security**：https://github.com/xiaozhe2018/GinForge/security

## 致谢

我们感谢以下安全研究人员对 GinForge 安全性的贡献：

- 暂无

如果你报告的安全漏洞被确认并修复，我们会在这里列出你的名字（如果你愿意）。

---

**安全是每个人的责任。感谢你帮助我们保持 GinForge 的安全！** 🔒


# 贡献指南

感谢你考虑为 GinForge 做出贡献！

## 📋 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [问题反馈](#问题反馈)

## 行为准则

本项目采用《贡献者公约》行为准则。参与本项目即表示你同意遵守其条款。

## 如何贡献

### 🐛 报告 Bug

如果你发现了 bug，请通过以下步骤报告：

1. **检查是否已存在相关 Issue**
2. **创建新 Issue** 并提供：
   - 清晰的标题和描述
   - 复现步骤
   - 预期行为和实际行为
   - 系统环境信息（Go 版本、操作系统等）
   - 相关日志或截图

### ✨ 提议新功能

我们欢迎新功能建议！请：

1. 创建一个 Feature Request Issue
2. 描述功能的用途和预期效果
3. 说明为什么这个功能对项目有价值
4. 等待维护者的反馈

### 🔧 提交代码

#### 第一次贡献？

1. Fork 本仓库
2. Clone 你的 Fork
   ```bash
   git clone https://github.com/your-username/GinForge.git
   cd GinForge
   ```
3. 添加上游仓库
   ```bash
   git remote add upstream https://github.com/xiaozhe2018/GinForge.git
   ```
4. 创建新分支
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### 开发流程

1. **保持同步**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **进行修改**
   - 编写代码
   - 添加测试
   - 更新文档

3. **运行测试**
   ```bash
   make test
   go test -v ./...
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能"
   ```

5. **推送到你的 Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 填写 PR 模板
   - 等待代码审查

## 代码规范

### Go 代码规范

#### 格式化
```bash
# 格式化代码
go fmt ./...

# 检查代码
go vet ./...

# 使用 golangci-lint（推荐）
golangci-lint run
```

#### 命名规范
- **包名**：小写，简短，有意义
  ```go
  package user
  package middleware
  ```

- **函数名**：大驼峰（导出）或小驼峰（私有）
  ```go
  func GetUser() {}      // 导出
  func getUserByID() {}  // 私有
  ```

- **变量名**：小驼峰，简洁明了
  ```go
  var userName string
  var userID int64
  ```

- **常量**：大驼峰或全大写（特殊情况）
  ```go
  const MaxRetries = 3
  const API_VERSION = "v1"
  ```

#### 注释规范
```go
// GetUser 获取用户信息
// 参数：
//   - userID: 用户ID
// 返回：
//   - *User: 用户对象
//   - error: 错误信息
func GetUser(userID int64) (*User, error) {
    // 实现...
}
```

#### 错误处理
```go
// 好的做法
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// 避免忽略错误
result, _ := someFunc()  // 不推荐
```

### TypeScript/Vue 代码规范

#### 格式化
```bash
cd web/admin
npm run lint
npm run type-check
```

#### 命名规范
- **组件名**：大驼峰
  ```typescript
  UserList.vue
  RoleForm.vue
  ```

- **变量/函数**：小驼峰
  ```typescript
  const userName = ref('')
  const getUserList = () => {}
  ```

- **常量**：全大写下划线分隔
  ```typescript
  const API_BASE_URL = 'http://localhost:8083'
  ```

#### 类型定义
```typescript
// 使用接口定义类型
interface User {
  id: number
  username: string
  email: string
}

// 导出类型
export type { User }
```

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

### 提交类型

- **feat**: 新功能
- **fix**: Bug 修复
- **docs**: 文档更新
- **style**: 代码格式（不影响功能）
- **refactor**: 重构（不是新功能也不是修复）
- **perf**: 性能优化
- **test**: 测试相关
- **chore**: 构建过程或辅助工具的变动
- **ci**: CI 配置文件和脚本的变动

### 提交格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 示例

```bash
# 简单提交
feat: 添加用户导出功能

# 详细提交
feat(admin): 添加用户批量删除功能

实现了用户批量删除功能，包括：
- 前端选择框
- 批量删除 API
- 删除确认对话框

Closes #123
```

## Pull Request 规范

### PR 标题
使用与 commit 相同的格式：
```
feat(admin): 添加用户导出功能
```

### PR 描述模板
```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化
- [ ] 测试
- [ ] 其他

## 变更说明
描述本次 PR 的主要变更...

## 相关 Issue
Closes #issue_number

## 测试
- [ ] 已添加单元测试
- [ ] 已添加集成测试
- [ ] 手动测试通过

## 截图（如适用）
贴上截图...

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已更新相关文档
- [ ] 所有测试通过
- [ ] 没有引入新的警告
```

## 问题反馈

### Issue 模板

#### Bug Report
```markdown
**描述问题**
简洁清晰地描述 bug。

**复现步骤**
1. 执行 '...'
2. 点击 '...'
3. 看到错误

**预期行为**
描述你期望发生的行为。

**实际行为**
描述实际发生的行为。

**环境信息**
- OS: [e.g. macOS 12.0]
- Go Version: [e.g. 1.21]
- GinForge Version: [e.g. 1.0.0]

**日志和截图**
如果适用，添加日志或截图。

**额外信息**
其他相关信息。
```

#### Feature Request
```markdown
**功能描述**
简洁清晰地描述你想要的功能。

**使用场景**
描述这个功能的使用场景和价值。

**建议的解决方案**
描述你期望的实现方式。

**备选方案**
描述你考虑过的其他方案。

**额外信息**
其他相关信息。
```

## 代码审查

所有提交都需要经过代码审查。审查者会检查：

- ✅ 代码质量和可读性
- ✅ 测试覆盖率
- ✅ 文档完整性
- ✅ 性能影响
- ✅ 安全性
- ✅ 向后兼容性

## 发布流程

1. 更新 CHANGELOG.md
2. 更新版本号
3. 创建 Git 标签
4. 发布到 GitHub Releases

## 获取帮助

如有疑问，可以：

- 📖 查看[文档](./docs/INDEX.md)
- 💬 在 [Issues](https://github.com/xiaozhe2018/GinForge/issues) 提问
- 🤝 加入社区讨论

## 致谢

感谢所有贡献者！你们让 GinForge 变得更好！

---

**让我们一起让 GinForge 变得更好！** 🚀


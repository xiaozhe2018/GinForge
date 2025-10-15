# GinForge CRUD 代码生成器使用指南

## 📖 简介

GinForge CRUD 代码生成器是一个强大的脚手架工具，可以根据数据库表结构自动生成完整的 CRUD 代码，包括：

✅ **后端代码**：
- Model（数据模型）
- Repository（数据访问层）
- Service（业务逻辑层）
- Handler（HTTP 处理层）
- Swagger 注释

✅ **前端代码**：
- TypeScript API 定义
- Vue 3 列表页面（带搜索、分页、排序）
- Vue 3 表单页面（带验证）

✅ **智能特性**：
- 自动识别字段类型并映射
- 自动生成验证规则
- 支持软删除、时间戳
- 支持搜索、分页、排序
- 自定义配置文件

---

## 🚀 快速开始

### 1. 查看所有数据库表

```bash
go run cmd/generator/main.go list:tables
```

输出示例：
```
🚀 GinForge 数据库表列表
================================

找到 12 个表:

  1. admin_users
  2. admin_roles
  3. admin_permissions
  4. admin_menus
  5. articles
  6. categories
  7. tags
  8. comments
  ...
```

### 2. 生成完整的 CRUD 代码

最简单的方式：

```bash
go run cmd/generator/main.go gen:crud --table=articles --module=admin
```

这个命令会：
1. 读取 `articles` 表的结构
2. 生成所有后端代码到 `services/admin-api/`
3. 生成所有前端代码到 `web/admin/src/`

### 3. 查看生成结果

```
✅ 代码生成完成！

📁 生成的文件:
  ✅ services/admin-api/internal/model/article.go
  ✅ services/admin-api/internal/repository/article_repository.go
  ✅ services/admin-api/internal/service/article_service.go
  ✅ services/admin-api/internal/handler/article_handler.go
  ✅ web/admin/src/api/article.ts
  ✅ web/admin/src/views/Article/index.vue
  ✅ web/admin/src/views/Article/Form.vue
```

---

## 📚 命令详解

### gen:crud - 生成完整CRUD代码

```bash
go run cmd/generator/main.go gen:crud [选项]
```

#### 选项

| 选项 | 简写 | 说明 | 必填 | 默认值 |
|------|------|------|------|--------|
| --table | -t | 数据库表名 | 是* | - |
| --module | -m | 模块名称 | 否 | admin |
| --config | -c | 配置文件路径 | 是* | - |
| --output | -o | 输出目录 | 否 | . |
| --frontend | - | 生成前端代码 | 否 | true |
| --force | -f | 强制覆盖已存在的文件 | 否 | false |
| --dry-run | - | 预览模式，不实际创建文件 | 否 | false |
| --verbose | -v | 显示详细输出 | 否 | false |

\* 注：`--table` 和 `--config` 二选一，必须提供其中一个

#### 示例

**示例 1：从数据库表生成**
```bash
# 基础用法
go run cmd/generator/main.go gen:crud --table=articles --module=admin

# 只生成后端代码
go run cmd/generator/main.go gen:crud --table=articles --module=admin --frontend=false

# 强制覆盖已存在的文件
go run cmd/generator/main.go gen:crud --table=articles --module=admin --force

# 预览生成结果（不实际创建文件）
go run cmd/generator/main.go gen:crud --table=articles --module=admin --dry-run

# 显示详细输出
go run cmd/generator/main.go gen:crud --table=articles --module=admin --verbose
```

**示例 2：从配置文件生成**
```bash
go run cmd/generator/main.go gen:crud --config=generator/articles.yaml
```

---

### gen:model - 只生成Model

如果只需要生成数据模型，可以使用这个命令：

```bash
go run cmd/generator/main.go gen:model --table=articles --module=admin
```

---

### init:config - 生成配置文件模板

创建一个配置文件模板，方便自定义生成规则：

```bash
go run cmd/generator/main.go init:config --table=articles
```

这会创建 `generator/articles.yaml` 文件，内容如下：

```yaml
table: articles
module: admin
model_name: Article
resource_name: articles
fields:
  - name: id
    type: bigint unsigned
    go_type: uint64
    ts_type: number
    is_primary_key: true
    auto_increment: true
    label: ID
    list_visible: true
    form_visible: false
    # ...更多字段配置

features:
  soft_delete: true
  timestamps: true
  pagination: true
  search: true
  sort: true

frontend:
  title: 文章管理
  icon: Document
  show_in_menu: true
```

然后您可以编辑这个文件来自定义生成规则，再使用：

```bash
go run cmd/generator/main.go gen:crud --config=generator/articles.yaml
```

---

### list:tables - 列出所有表

```bash
go run cmd/generator/main.go list:tables
```

---

## 🎨 配置文件详解

### 基础配置

```yaml
table: articles           # 数据库表名
module: admin            # 模块名称（admin/user/file）
model_name: Article      # 模型名称（PascalCase）
resource_name: articles  # 资源名称（复数形式，用于 URL）
```

### 字段配置

每个字段支持以下配置：

```yaml
fields:
  - name: title                  # 字段名（snake_case）
    type: varchar(255)           # 数据库类型
    go_type: string              # Go 类型
    ts_type: string              # TypeScript 类型
    nullable: false              # 是否可为空
    is_primary_key: false        # 是否主键
    auto_increment: false        # 是否自增
    default_value: ""            # 默认值
    comment: "文章标题"           # 字段注释
    
    # 验证规则
    validations:
      - required               # 必填
      - max:255                # 最大长度
    
    # UI 配置
    label: "标题"               # 显示标签
    form_type: input           # 表单类型
    list_visible: true         # 列表中是否显示
    form_visible: true         # 表单中是否显示
    searchable: true           # 是否可搜索
    sortable: true             # 是否可排序
```

### 表单类型 (form_type)

| 类型 | 说明 | 示例 |
|------|------|------|
| input | 单行输入框 | 标题、名称 |
| textarea | 多行输入框 | 描述、备注 |
| password | 密码输入框 | 密码 |
| email | 邮箱输入框 | 邮箱 |
| number | 数字输入框 | 年龄、数量 |
| switch | 开关 | 状态、是否启用 |
| select | 下拉选择 | 分类、类型 |
| date | 日期选择器 | 出生日期 |
| datetime | 日期时间选择器 | 创建时间 |
| upload | 文件上传 | 头像、图片 |
| editor | 富文本编辑器 | 文章内容 |

### 验证规则 (validations)

| 规则 | 说明 | 示例 |
|------|------|------|
| required | 必填 | validations: [required] |
| email | 邮箱格式 | validations: [required, email] |
| min:6 | 最小长度 | validations: [required, min:6] |
| max:255 | 最大长度 | validations: [required, max:255] |
| len:11 | 固定长度 | validations: [required, len:11] |
| url | URL 格式 | validations: [url] |

### 功能特性

```yaml
features:
  soft_delete: true    # 软删除（自动识别 deleted_at 字段）
  timestamps: true     # 时间戳（自动识别 created_at, updated_at）
  pagination: true     # 分页
  search: true         # 搜索
  sort: true           # 排序
  export: false        # 导出（暂未实现）
  import: false        # 导入（暂未实现）
  batch_delete: false  # 批量删除（暂未实现）
```

### 前端配置

```yaml
frontend:
  title: "文章管理"     # 页面标题
  icon: "Document"      # 菜单图标（Element Plus Icon）
  show_in_menu: true   # 是否显示在菜单
  menu_parent: ""      # 父菜单（暂未使用）
```

---

## 🎯 生成后的后续步骤

### 1. 注册路由

在 `services/{module}-api/internal/router/router.go` 中添加路由：

```go
// 初始化 Handler
articleHandler := handler.NewArticleHandler(articleService, log)

// 注册路由
auth.GET("/articles", articleHandler.List)
auth.GET("/articles/:id", articleHandler.Get)
auth.POST("/articles", articleHandler.Create)
auth.PUT("/articles/:id", articleHandler.Update)
auth.DELETE("/articles/:id", articleHandler.Delete)
```

### 2. 注册前端路由

在 `web/admin/src/router/index.ts` 中添加路由：

```typescript
{
  path: 'articles',
  name: 'ArticleList',
  component: () => import('@/views/Article/index.vue'),
  meta: { title: '文章管理', requiresAuth: true }
}
```

### 3. 添加菜单

在 `web/admin/src/layout/index.vue` 中添加菜单项：

```vue
<el-menu-item index="/dashboard/articles">
  <el-icon><Document /></el-icon>
  <span>文章管理</span>
</el-menu-item>
```

### 4. 测试功能

1. 重启后端服务
2. 刷新前端页面
3. 访问新生成的页面
4. 测试增删改查功能

---

## 💡 最佳实践

### 1. 使用配置文件进行自定义

对于复杂的业务场景，建议先生成配置文件，再根据需求调整：

```bash
# 1. 生成配置文件
go run cmd/generator/main.go init:config --table=articles

# 2. 编辑 generator/articles.yaml，调整字段配置

# 3. 从配置文件生成代码
go run cmd/generator/main.go gen:crud --config=generator/articles.yaml
```

### 2. 预览模式避免错误

使用 `--dry-run` 选项预览生成结果：

```bash
go run cmd/generator/main.go gen:crud --table=articles --module=admin --dry-run
```

### 3. 分模块管理

不同的业务模块使用不同的 module：

```bash
# 管理后台相关
go run cmd/generator/main.go gen:crud --table=admin_users --module=admin

# 用户相关
go run cmd/generator/main.go gen:crud --table=user_profiles --module=user

# 文件相关
go run cmd/generator/main.go gen:crud --table=files --module=file
```

### 4. 版本控制

将生成的配置文件加入版本控制：

```bash
git add generator/
git commit -m "Add generator config for articles"
```

---

## 🔧 常见问题

### Q1: 生成的代码可以修改吗？

**A:** 可以！生成的代码只是一个起点，您可以根据实际需求修改。建议：
- 首次生成后，立即调整代码以满足业务需求
- 之后不要使用 `--force` 选项重新生成，以免覆盖您的修改
- 如果需要重新生成，建议使用版本控制系统（Git）

### Q2: 如何自定义字段的表单类型？

**A:** 使用配置文件：

```bash
# 1. 生成配置文件
go run cmd/generator/main.go init:config --table=articles

# 2. 编辑配置文件，修改 form_type
# generator/articles.yaml
fields:
  - name: content
    form_type: editor  # 改为富文本编辑器

# 3. 重新生成
go run cmd/generator/main.go gen:crud --config=generator/articles.yaml --force
```

### Q3: 如何支持关联查询？

**A:** 生成器目前生成的是基础 CRUD，关联查询需要手动添加：

1. 在 Repository 中添加关联查询方法
2. 在 Service 中调用
3. 在 Model 中定义关联结构体

示例（文章关联作者）：

```go
// Model
type Article struct {
    ID       uint64 `json:"id"`
    Title    string `json:"title"`
    AuthorID uint64 `json:"author_id"`
    Author   *User  `json:"author" gorm:"foreignKey:AuthorID"`
}

// Repository
func (r *ArticleRepository) GetWithAuthor(id uint64) (*model.Article, error) {
    var article model.Article
    err := r.db.Preload("Author").First(&article, id).Error
    return &article, err
}
```

### Q4: 前端表单需要添加下拉选项怎么办？

**A:** 在生成的 Vue 文件中，找到 `TODO: 添加选项` 的注释，手动添加选项：

```vue
<el-select v-model="form.category_id" placeholder="请选择分类">
  <!-- 静态选项 -->
  <el-option label="技术" value="1" />
  <el-option label="生活" value="2" />
  
  <!-- 或从 API 动态加载 -->
  <el-option
    v-for="item in categoryOptions"
    :key="item.id"
    :label="item.name"
    :value="item.id"
  />
</el-select>
```

### Q5: 如何修改生成的代码风格？

**A:** 修改模板文件：

1. 模板文件位置：`pkg/generator/template_*.go`
2. 找到对应的模板（model/repository/service/handler/frontend）
3. 修改模板内容
4. 重新运行生成器

---

## 📊 字段类型映射表

### MySQL → Go 类型

| MySQL 类型 | Go 类型 | 备注 |
|-----------|---------|------|
| tinyint | int8 | 1字节整数 |
| smallint | int16 | 2字节整数 |
| int, integer | int | 4字节整数 |
| bigint | int64 | 8字节整数 |
| float | float32 | 单精度浮点 |
| double | float64 | 双精度浮点 |
| decimal | float64 | 高精度小数 |
| varchar, char, text | string | 字符串 |
| date, datetime, timestamp | time.Time | 时间 |
| json | string | JSON字符串 |
| blob | []byte | 二进制数据 |

### Go → TypeScript 类型

| Go 类型 | TypeScript 类型 |
|---------|----------------|
| int, int8, int16, int32, int64 | number |
| uint, uint8, uint16, uint32, uint64 | number |
| float32, float64 | number |
| string | string |
| bool | boolean |
| time.Time | string |
| []byte | string |

---

## 🎨 示例：完整的工作流程

### 场景：创建文章管理模块

#### 1. 准备数据库表

```sql
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `author_id` bigint unsigned NOT NULL COMMENT '作者ID',
  `category_id` bigint unsigned DEFAULT NULL COMMENT '分类ID',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0草稿,1已发布',
  `view_count` int unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';
```

#### 2. 生成配置文件

```bash
go run cmd/generator/main.go init:config --table=articles
```

#### 3. 调整配置文件

编辑 `generator/articles.yaml`：

```yaml
# 修改 content 字段的表单类型为富文本编辑器
fields:
  - name: content
    form_type: editor

# 修改 status 字段为下拉选择
  - name: status
    form_type: select
    
# 隐藏 view_count 在表单中显示
  - name: view_count
    form_visible: false
```

#### 4. 生成代码

```bash
go run cmd/generator/main.go gen:crud --config=generator/articles.yaml --verbose
```

#### 5. 注册路由

在 `services/admin-api/internal/router/router.go`:

```go
// 初始化服务和处理器
articleRepo := repository.NewArticleRepository(database)
articleService := service.NewArticleService(articleRepo, log)
articleHandler := handler.NewArticleHandler(articleService, log)

// 注册路由
auth.GET("/articles", articleHandler.List)
auth.GET("/articles/:id", articleHandler.Get)
auth.POST("/articles", articleHandler.Create)
auth.PUT("/articles/:id", articleHandler.Update)
auth.DELETE("/articles/:id", articleHandler.Delete)
```

#### 6. 添加前端路由

在 `web/admin/src/router/index.ts`:

```typescript
{
  path: 'articles',
  name: 'ArticleList',
  component: () => import('@/views/Article/index.vue'),
  meta: { title: '文章管理', requiresAuth: true }
}
```

#### 7. 添加菜单

在 `web/admin/src/layout/index.vue`:

```vue
<el-menu-item index="/dashboard/articles">
  <el-icon><Document /></el-icon>
  <span>文章管理</span>
</el-menu-item>
```

#### 8. 重启服务并测试

```bash
# 后端
cd services/admin-api
go run cmd/server/main.go

# 前端
cd web/admin
npm run dev
```

访问 `http://localhost:3000/dashboard/articles`，测试功能！

---

## 🚀 高级功能

### 自定义模板

如果需要修改生成的代码风格，可以修改模板文件：

```
pkg/generator/
├── template_model.go           # Model 模板
├── template_repository.go      # Repository 模板
├── template_service.go         # Service 模板
├── template_handler.go         # Handler 模板
├── template_frontend_api.go    # 前端 API 模板
├── template_frontend_list.go   # 前端列表页模板
└── template_frontend_form.go   # 前端表单页模板
```

修改模板后，重新运行生成器即可。

---

## 📝 总结

GinForge CRUD 代码生成器可以帮助您：

✅ **节省 80% 的重复代码编写时间**  
✅ **保证代码风格统一**  
✅ **减少人为错误**  
✅ **快速搭建项目原型**  
✅ **专注于业务逻辑开发**

开始使用代码生成器，让开发更高效！🎉

---

**如有问题或建议，欢迎提 Issue 或 PR！**


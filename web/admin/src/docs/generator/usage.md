# 代码生成器详细使用指南

## 📚 命令详解

### gen:crud - 生成完整CRUD代码

这是最常用的命令，可以生成完整的 CRUD 功能代码。

#### 基础用法

```bash
./bin/generator gen:crud --table=<表名> --module=<模块>
```

#### 完整选项

| 选项 | 简写 | 说明 | 必填 | 默认值 |
|------|------|------|------|--------|
| --table | -t | 数据库表名 | 是* | - |
| --module | -m | 模块名称 | 否 | admin |
| --config | -c | 配置文件路径 | 是* | - |
| --output | -o | 输出目录 | 否 | . |
| --frontend | - | 生成前端代码 | 否 | true |
| --force | -f | 强制覆盖已存在的文件 | 否 | false |
| --auto-register | -a | 自动注册路由和菜单 ⭐ | 否 | false |
| --dry-run | - | 预览模式，不实际创建文件 | 否 | false |
| --verbose | -v | 显示详细输出 | 否 | false |

\* 注：`--table` 和 `--config` 二选一，必须提供其中一个

#### 使用示例

**示例 1：一键生成（推荐）⭐**

```bash
# 生成代码并自动注册路由和菜单
./bin/generator gen:crud --table=articles --module=admin -a

# 只需重启服务即可使用！
```

**这是最快的方式！** 生成器会自动完成所有工作，包括注册路由和菜单。

**示例 2：预览模式（最佳实践）**

```bash
# 先预览会生成什么
./bin/generator gen:crud --table=articles --module=admin -a --dry-run

# 确认无误后正式生成
./bin/generator gen:crud --table=articles --module=admin -a
```

先预览会生成哪些文件和注册哪些路由，确认无误后再正式生成。

**示例 3：从配置文件生成**

```bash
./bin/generator gen:crud --config=generator/articles.yaml -a
```

**示例 4：只生成代码，手动注册**

```bash
./bin/generator gen:crud --table=articles --module=admin
```

不使用 `-a` 选项，需要手动注册路由和菜单。

**示例 5：只生成后端代码**

```bash
./bin/generator gen:crud --table=articles --module=admin --frontend=false
```

**示例 6：强制覆盖已存在的文件**

```bash
./bin/generator gen:crud --table=articles --module=admin -a --force
```

**示例 7：显示详细输出**

```bash
./bin/generator gen:crud --table=articles --module=admin -a --verbose
```

---

### gen:model - 只生成Model

如果只需要生成数据模型，不需要完整的 CRUD 功能：

```bash
./bin/generator gen:model --table=articles --module=admin
```

这个命令只会生成 Model 文件，包括：
- 数据模型结构体
- 请求结构体（CreateRequest、UpdateRequest、ListRequest）
- 响应结构体（Response）
- 转换方法

---

### init:config - 生成配置文件模板

创建一个配置文件模板，方便自定义生成规则：

```bash
./bin/generator init:config --table=articles
```

这会创建 `generator/articles.yaml` 文件，内容包括：

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
    label: "ID"
    form_type: input
    list_visible: true
    form_visible: false
    # ... 更多字段配置

features:
  soft_delete: true
  timestamps: true
  pagination: true
  search: true
  sort: true

frontend:
  title: "文章管理"
  icon: "Document"
  show_in_menu: true
```

然后您可以编辑这个文件来自定义生成规则，再使用：

```bash
./bin/generator gen:crud --config=generator/articles.yaml
```

---

### list:tables - 列出所有表

查看数据库中的所有表：

```bash
./bin/generator list:tables
```

输出：
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
  ...

💡 使用示例:
  generator gen:crud --table=<表名> --module=admin
```

---

## ⭐ 自动注册功能

### 什么是自动注册？

自动注册是代码生成器的**杀手级功能**，可以自动完成以下工作：

✅ **自动注册后端路由**
- 在 `router.go` 中添加 Handler 初始化代码
- 自动注册 5 个 CRUD 路由（GET、POST、PUT、DELETE）

✅ **自动注册前端路由**
- 在 `router/index.ts` 中添加页面路由配置
- 自动配置路由元信息

✅ **自动注册菜单**
- 在 `layout/index.vue` 中添加菜单项
- 自动导入需要的图标

### 使用方式

只需添加 `--auto-register` 或 `-a` 选项：

```bash
./bin/generator gen:crud --table=articles --module=admin -a
```

### 智能特性

- ✅ **防重复注册** - 智能检测是否已注册，避免重复
- ✅ **优雅的错误处理** - 部分失败不影响整体
- ✅ **预览模式** - 使用 `--dry-run` 预览不实际修改
- ✅ **详细输出** - 使用 `--verbose` 查看详细过程

### 自动修改的文件

| 文件 | 修改内容 |
|------|---------|
| `services/{module}-api/internal/router/router.go` | 添加 Handler 初始化和路由注册 |
| `web/admin/src/router/index.ts` | 添加页面路由配置 |
| `web/admin/src/layout/index.vue` | 添加菜单项和图标导入 |

### 注意事项

1. **使用版本控制** - 建议先提交当前代码，方便回滚
2. **首次使用先预览** - 使用 `--dry-run` 查看会修改哪些文件
3. **检查注册结果** - 生成后检查路由和菜单是否正确

### 效率对比

| 方式 | 耗时 | 步骤 |
|------|------|------|
| 手动注册 | 12 分钟 | 生成代码 + 手动编辑 3 个文件 |
| **自动注册** | **2 分钟** | **一条命令 + 重启服务** |
| **效率提升** | **6 倍** | - |

---

## 🎯 生成后的后续步骤

> **💡 提示**：如果使用了 `--auto-register` 或 `-a` 选项，以下步骤会自动完成，您只需重启服务即可！

### 手动注册流程（不使用 -a 时需要）

#### 1. 注册后端路由

在 `services/{module}-api/internal/router/router.go` 中添加路由：

```go
// 初始化 Repository、Service、Handler
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

## 💡 使用技巧

### 1. 使用别名简化命令

在 `~/.zshrc` 或 `~/.bashrc` 中添加：

```bash
alias gen='./bin/generator'
```

然后就可以这样使用：

```bash
gen list:tables
gen gen:crud --table=articles --module=admin
```

### 2. 批量生成多个表

创建一个脚本 `generate_all.sh`：

```bash
#!/bin/bash

tables=("articles" "categories" "tags" "comments")

for table in "${tables[@]}"
do
    echo "生成 $table..."
    ./bin/generator gen:crud --table=$table --module=admin
done

echo "全部生成完成！"
```

运行：

```bash
chmod +x generate_all.sh
./generate_all.sh
```

### 3. 使用配置文件版本控制

```bash
# 将配置文件加入版本控制
git add generator/*.yaml
git commit -m "Add generator configs"
```

### 4. 预览模式避免错误

使用 `--dry-run` 选项预览生成结果：

```bash
./bin/generator gen:crud --table=articles --module=admin --dry-run
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
./bin/generator init:config --table=articles

# 2. 编辑配置文件，修改 form_type
# generator/articles.yaml
fields:
  - name: content
    form_type: editor  # 改为富文本编辑器

# 3. 重新生成
./bin/generator gen:crud --config=generator/articles.yaml --force
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
4. 重新编译生成器：`go build -o bin/generator ./cmd/generator`

---

## 📊 生成的文件结构

### 后端代码结构

```
services/admin-api/
├── internal/
│   ├── model/
│   │   └── article.go                 # 数据模型
│   ├── repository/
│   │   └── article_repository.go      # 数据访问层
│   ├── service/
│   │   └── article_service.go         # 业务逻辑层
│   └── handler/
│       └── article_handler.go         # HTTP 处理层
```

### 前端代码结构

```
web/admin/src/
├── api/
│   └── article.ts                     # API 接口定义
└── views/
    └── Article/
        ├── index.vue                  # 列表页面
        └── Form.vue                   # 表单页面
```

---

## 🎯 最佳实践

### 1. 使用配置文件进行自定义

对于复杂的业务场景，建议先生成配置文件，再根据需求调整：

```bash
# 1. 生成配置文件
./bin/generator init:config --table=articles

# 2. 编辑 generator/articles.yaml，调整字段配置

# 3. 从配置文件生成代码
./bin/generator gen:crud --config=generator/articles.yaml
```

### 2. 预览模式避免错误

使用 `--dry-run` 选项预览生成结果：

```bash
./bin/generator gen:crud --table=articles --module=admin --dry-run
```

### 3. 分模块管理

不同的业务模块使用不同的 module：

```bash
# 管理后台相关
./bin/generator gen:crud --table=admin_users --module=admin

# 用户相关
./bin/generator gen:crud --table=user_profiles --module=user

# 文件相关
./bin/generator gen:crud --table=files --module=file
```

### 4. 版本控制

将生成的配置文件加入版本控制：

```bash
git add generator/
git commit -m "Add generator config for articles"
```

---

## 📖 下一步

- ⚙️ [配置文件详解](./configuration) - 了解所有配置选项
- 💼 [实战示例](./examples) - 完整的业务模块开发流程

**掌握这些命令，开发效率提升 10 倍！** 🚀


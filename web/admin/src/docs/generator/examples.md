# 实战示例：完整的业务模块开发

## 🎯 场景：创建文章管理模块

我们将从零到一创建一个完整的文章管理模块，包括前后端的全部功能。

---

## 🚀 快速入门（5分钟）

### 最快的方式（使用自动注册）⭐

如果您只想快速体验，使用以下简化流程：

```bash
# 1. 编译生成器（首次使用）
go build -o bin/generator ./cmd/generator

# 2. 创建数据库表（执行 SQL）
# ... 见下文的数据库表设计

# 3. 一键生成（自动注册所有内容）
./bin/generator gen:crud --table=articles --module=admin -a

# 4. 重启服务
cd services/admin-api && go run cmd/server/main.go
# 刷新前端浏览器

# 完成！立即可用！✅
```

**只需这 4 步，5 分钟完成！** 🎉

### 完整流程（自定义配置）

如果需要自定义字段配置，使用以下完整流程（约 80 分钟）：

1. 数据库表设计
2. 生成配置文件并编辑
3. 生成代码（使用自动注册）
4. 自定义表单
5. 扩展功能
6. 测试调试

**下面是详细的步骤说明。**

---

## 📋 需求分析

### 功能需求

- ✅ 文章列表（分页、搜索、排序）
- ✅ 创建文章
- ✅ 编辑文章
- ✅ 删除文章（软删除）
- ✅ 文章状态管理（草稿、已发布、已下线）
- ✅ 置顶功能
- ✅ 浏览次数统计

### 数据库表设计

```sql
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `author_id` bigint unsigned NOT NULL COMMENT '作者ID',
  `category_id` bigint unsigned DEFAULT NULL COMMENT '分类ID',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0草稿,1已发布,2已下线',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶',
  `view_count` int unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';
```

---

## 🚀 第一步：生成配置文件

### 1.1 创建配置文件

```bash
cd /Users/chaojidoudou/project/go/GinForge
./bin/generator init:config --table=articles
```

### 1.2 编辑配置文件

编辑 `generator/articles.yaml`，自定义字段配置：

```yaml
table: articles
module: admin
model_name: Article
resource_name: articles

fields:
  # ... id 字段保持默认

  # 标题 - 可搜索
  - name: title
    type: varchar(255)
    go_type: string
    ts_type: string
    nullable: false
    comment: "标题"
    validations:
      - required
      - max:255
    label: "标题"
    form_type: input
    list_visible: true
    form_visible: true
    searchable: true        # 可搜索
    sortable: true

  # 内容 - 富文本编辑器
  - name: content
    type: text
    go_type: string
    ts_type: string
    nullable: false
    comment: "内容"
    validations:
      - required
    label: "内容"
    form_type: editor       # 富文本编辑器
    list_visible: false     # 列表中不显示
    form_visible: true
    searchable: true
    sortable: false

  # 作者ID - 下拉选择
  - name: author_id
    type: bigint unsigned
    go_type: uint64
    ts_type: number
    nullable: false
    comment: "作者ID"
    validations:
      - required
    label: "作者"
    form_type: select       # 下拉选择
    list_visible: true
    form_visible: true
    searchable: false
    sortable: true

  # 分类ID - 下拉选择（可选）
  - name: category_id
    type: bigint unsigned
    go_type: "*uint64"
    ts_type: number
    nullable: true
    comment: "分类ID"
    validations: []
    label: "分类"
    form_type: select
    list_visible: true
    form_visible: true
    searchable: false
    sortable: true

  # 状态 - 下拉选择
  - name: status
    type: tinyint
    go_type: int8
    ts_type: number
    nullable: false
    default_value: "0"
    comment: "状态:0草稿,1已发布,2已下线"
    validations:
      - required
    label: "状态"
    form_type: select
    list_visible: true
    form_visible: true
    searchable: false
    sortable: true

  # 是否置顶 - 开关
  - name: is_top
    type: tinyint(1)
    go_type: bool
    ts_type: boolean
    nullable: false
    default_value: "0"
    comment: "是否置顶"
    validations: []
    label: "置顶"
    form_type: switch       # 开关
    list_visible: true
    form_visible: true
    searchable: false
    sortable: true

  # 浏览次数 - 只显示
  - name: view_count
    type: int unsigned
    go_type: uint
    ts_type: number
    nullable: false
    default_value: "0"
    comment: "浏览次数"
    validations: []
    label: "浏览次数"
    form_type: number
    list_visible: true
    form_visible: false     # 表单中不显示
    searchable: false
    sortable: true

  # ... 时间字段保持默认

features:
  soft_delete: true         # 启用软删除
  timestamps: true          # 启用时间戳
  pagination: true          # 启用分页
  search: true             # 启用搜索
  sort: true               # 启用排序

frontend:
  title: "文章管理"
  icon: "Document"
  show_in_menu: true
```

---

## ⚙️ 第二步：生成代码

### 2.1 预览生成结果

```bash
# 使用自动注册预览
./bin/generator gen:crud --config=generator/articles.yaml -a --dry-run
```

检查输出，确认要生成的文件和注册的路由。

### 2.2 正式生成（推荐使用自动注册）⭐

```bash
# 一键生成并自动注册
./bin/generator gen:crud --config=generator/articles.yaml -a --verbose
```

输出：
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

🔧 自动注册路由和菜单...
✅ 后端路由注册成功
✅ 前端路由注册成功
✅ 菜单注册成功
✅ 路由和菜单注册完成！
```

**使用自动注册后，可以跳过第三步和第四步，直接进行测试！**

---

## 🔧 第三步：注册后端路由

> **💡 提示**：如果使用了 `-a` 选项，此步骤已自动完成，可以跳过！

### 3.1 编辑路由文件（手动注册时需要）

编辑 `services/admin-api/internal/router/router.go`：

```go
package router

import (
    // ... 其他导入
    "goweb/services/admin-api/internal/repository"
    "goweb/services/admin-api/internal/service"
    "goweb/services/admin-api/internal/handler"
)

func NewRouter(/* ... 参数 */) *gin.Engine {
    // ... 其他代码
    
    // 初始化文章相关的 Repository、Service、Handler
    articleRepo := repository.NewArticleRepository(database)
    articleService := service.NewArticleService(articleRepo, log)
    articleHandler := handler.NewArticleHandler(articleService, log)
    
    // 注册路由
    api := r.Group("/api/v1/admin")
    {
        // ... 其他路由
        
        // 文章相关路由
        auth := api.Group("")
        auth.Use(middleware.JWTAuthWithRedis(redisClient, log))
        {
            // ... 其他路由
            
            // 文章路由
            auth.GET("/articles", articleHandler.List)
            auth.GET("/articles/:id", articleHandler.Get)
            auth.POST("/articles", articleHandler.Create)
            auth.PUT("/articles/:id", articleHandler.Update)
            auth.DELETE("/articles/:id", articleHandler.Delete)
        }
    }
    
    return r
}
```

### 3.2 重启后端服务

```bash
cd services/admin-api
go run cmd/server/main.go
```

### 3.3 测试 API

```bash
# 测试列表接口
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8083/api/v1/admin/articles?page=1&page_size=10"

# 测试创建接口
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试文章","content":"测试内容","author_id":1,"status":1}' \
  "http://localhost:8083/api/v1/admin/articles"
```

---

## 🎨 第四步：注册前端路由

> **💡 提示**：如果使用了 `-a` 选项，此步骤已自动完成，可以跳过！

### 4.1 编辑路由文件（手动注册时需要）

编辑 `web/admin/src/router/index.ts`：

```typescript
const routes = [
  // ... 其他路由
  {
    path: '/dashboard',
    component: () => import('@/layout/index.vue'),
    children: [
      // ... 其他路由
      
      // 文章管理
      {
        path: 'articles',
        name: 'ArticleList',
        component: () => import('@/views/Article/index.vue'),
        meta: { title: '文章管理', requiresAuth: true }
      },
    ]
  }
]
```

### 4.2 添加菜单

编辑 `web/admin/src/layout/index.vue`：

```vue
<template>
  <el-menu>
    <!-- ... 其他菜单 -->
    
    <!-- 文章管理 -->
    <el-menu-item index="/dashboard/articles">
      <el-icon><Document /></el-icon>
      <span>文章管理</span>
    </el-menu-item>
  </el-menu>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'
// ... 其他代码
</script>
```

---

## 🎯 第五步：自定义前端表单

### 5.1 添加下拉选项

编辑 `web/admin/src/views/Article/index.vue`，找到状态选择下拉框，添加选项：

```vue
<!-- 状态选择 -->
<el-form-item label="状态" prop="status">
  <el-select v-model="form.status" placeholder="请选择状态">
    <el-option label="草稿" :value="0" />
    <el-option label="已发布" :value="1" />
    <el-option label="已下线" :value="2" />
  </el-select>
</el-form-item>
```

### 5.2 添加作者选择

```vue
<el-form-item label="作者" prop="author_id">
  <el-select v-model="form.author_id" placeholder="请选择作者">
    <el-option
      v-for="author in authorOptions"
      :key="author.id"
      :label="author.name"
      :value="author.id"
    />
  </el-select>
</el-form-item>

<script setup lang="ts">
// 添加作者选项加载
const authorOptions = ref([])

const loadAuthors = async () => {
  try {
    const data = await userApi.getUserList({ page: 1, page_size: 100 })
    authorOptions.value = data.list
  } catch (error) {
    console.error('加载作者列表失败:', error)
  }
}

onMounted(() => {
  loadData()
  loadAuthors()
})
</script>
```

### 5.3 添加分类选择

类似地添加分类选择功能。

---

## 🧪 第六步：测试功能

### 6.1 刷新前端

访问 `http://localhost:3000`，登录后台。

### 6.2 测试列表功能

1. 点击"文章管理"菜单
2. 查看文章列表
3. 测试搜索功能
4. 测试分页功能
5. 测试排序功能

### 6.3 测试创建功能

1. 点击"新建文章"按钮
2. 填写表单
3. 提交创建
4. 检查是否创建成功

### 6.4 测试编辑功能

1. 点击"编辑"按钮
2. 修改内容
3. 保存
4. 检查是否更新成功

### 6.5 测试删除功能

1. 点击"删除"按钮
2. 确认删除
3. 检查是否删除成功（软删除）

---

## 💡 第七步：扩展功能

### 7.1 添加发布功能

在 Service 中添加发布方法：

```go
// Publish 发布文章
func (s *ArticleService) Publish(id uint64) error {
    article, err := s.repo.GetByID(id)
    if err != nil {
        return errors.New("文章不存在")
    }
    
    article.Status = 1
    now := time.Now()
    article.PublishedAt = &now
    
    if err := s.repo.Update(article); err != nil {
        s.logger.Error("发布文章失败", err, "id", id)
        return errors.New("发布文章失败")
    }
    
    return nil
}
```

在 Handler 中添加发布接口：

```go
// Publish 发布文章
func (h *ArticleHandler) Publish(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, 400, "ID 格式错误")
        return
    }
    
    if err := h.service.Publish(id); err != nil {
        response.Error(c, 500, err.Error())
        return
    }
    
    response.Success(c, nil)
}
```

注册路由：

```go
auth.PUT("/articles/:id/publish", articleHandler.Publish)
```

### 7.2 添加浏览量统计

在 Repository 中添加增加浏览量方法：

```go
// IncrementViewCount 增加浏览次数
func (r *ArticleRepository) IncrementViewCount(id uint64) error {
    return r.db.Model(&model.Article{}).
        Where("id = ?", id).
        UpdateColumn("view_count", gorm.Expr("view_count + 1")).
        Error
}
```

在获取文章详情时调用：

```go
func (s *ArticleService) GetByID(id uint64) (*model.Article, error) {
    article, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 异步增加浏览次数
    go func() {
        s.repo.IncrementViewCount(id)
    }()
    
    return article, nil
}
```

---

## 📊 完成效果

### 功能清单

- ✅ 文章列表（分页、搜索、排序） - 自动生成
- ✅ 创建文章 - 自动生成
- ✅ 编辑文章 - 自动生成
- ✅ 删除文章 - 自动生成
- ✅ 文章状态管理 - 手动添加下拉选项
- ✅ 置顶功能 - 自动生成（开关）
- ✅ 浏览次数统计 - 手动扩展
- ✅ 发布功能 - 手动扩展

### 时间统计

#### 使用自动注册（-a）⭐

| 阶段 | 时间 |
|------|------|
| 数据库表设计 | 15 分钟 |
| 生成配置文件 | 5 分钟 |
| 代码生成 + 自动注册 | 10 秒 |
| ~~注册路由~~ | ~~自动完成~~ |
| ~~添加菜单~~ | ~~自动完成~~ |
| 自定义表单 | 10 分钟 |
| 扩展功能 | 20 分钟 |
| 测试调试 | 30 分钟 |
| **总计** | **约 80 分钟** |

**相比传统手写（6-7 小时），效率提升 5-9 倍！** ⚡

#### 手动注册

| 阶段 | 时间 |
|------|------|
| 数据库表设计 | 15 分钟 |
| 生成配置文件 | 5 分钟 |
| 代码生成 | 10 秒 |
| 注册路由 | 5 分钟 |
| 添加菜单 | 3 分钟 |
| 自定义表单 | 10 分钟 |
| 扩展功能 | 20 分钟 |
| 测试调试 | 30 分钟 |
| **总计** | **约 90 分钟** |

**相比传统手写（6-7 小时），效率提升 5-7 倍！** ⚡

---

## 🎉 总结

### 学到的技能

通过这个完整的实战示例，我们学会了：

1. ✅ 如何设计数据库表
2. ✅ 如何生成和自定义配置文件
3. ✅ 如何使用代码生成器生成代码
4. ✅ 如何使用自动注册功能 ⭐
5. ✅ 如何自定义前端表单
6. ✅ 如何测试功能
7. ✅ 如何扩展新功能

### 两种开发模式对比

#### 🚀 快速模式（5-10 分钟）

**适用场景**：标准 CRUD，快速原型

```bash
# 一条命令完成所有工作
./bin/generator gen:crud --table=articles --module=admin -a

# 重启服务，立即使用
```

**优势**：
- ⚡ 超快速度（5-10 分钟）
- 🎯 自动注册所有内容
- 📦 开箱即用

#### 🎨 定制模式（80-90 分钟）

**适用场景**：复杂业务，需要自定义

```bash
# 1. 生成配置文件
./bin/generator init:config --table=articles

# 2. 编辑配置文件（自定义字段、验证规则等）

# 3. 生成代码并自动注册
./bin/generator gen:crud --config=generator/articles.yaml -a

# 4. 自定义表单（添加选项、关联等）

# 5. 扩展功能（添加业务逻辑）

# 6. 测试调试
```

**优势**：
- 🎨 高度可定制
- 💼 符合业务需求
- 🏗️ 可扩展

### 效率总结

| 开发模式 | 耗时 | 效率提升 | 适用场景 |
|---------|------|---------|---------|
| 传统手写 | 6-7 小时 | 基准 | - |
| 快速模式（自动注册） | 5-10 分钟 | **40-80 倍** ⚡ | 标准 CRUD |
| 定制模式（自动注册） | 80-90 分钟 | **5-9 倍** | 复杂业务 |

**现在您可以用同样的方法快速开发其他业务模块了！** 🚀

---

## 📖 更多示例

### 示例 2：用户管理模块

```bash
./bin/generator gen:crud --table=users --module=user
```

### 示例 3：分类管理模块

```bash
./bin/generator gen:crud --table=categories --module=admin
```

### 示例 4：评论管理模块

```bash
./bin/generator gen:crud --table=comments --module=admin
```

**掌握代码生成器，快速构建完整系统！** 💪


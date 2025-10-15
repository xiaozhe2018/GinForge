# CRUD 代码生成器介绍

## 🎯 什么是代码生成器？

**GinForge CRUD 代码生成器**是一个强大的脚手架工具，可以根据数据库表结构自动生成完整的 CRUD 代码，包括后端和前端的全部代码。

### 核心功能

✅ **自动生成后端代码**
- Model（数据模型 + 请求/响应结构）
- Repository（CRUD + 分页 + 搜索 + 排序）
- Service（业务逻辑 + 验证）
- Handler（HTTP 处理 + Swagger 注释）

✅ **自动生成前端代码**
- TypeScript API 接口定义
- Vue 3 列表页面（带搜索、分页、排序）
- Vue 3 表单页面（带数据验证）

✅ **智能特性**
- 字段类型自动映射（MySQL → Go → TypeScript）
- 表单类型智能识别（input/textarea/select/switch/date）
- 验证规则自动生成（required/email/min/max）
- 软删除支持（自动识别 deleted_at）
- 时间戳支持（自动识别 created_at/updated_at）

---

## 🚀 快速开始

### 1. 编译生成器

```bash
cd /Users/chaojidoudou/project/go/GinForge
go build -o bin/generator ./cmd/generator
```

### 2. 查看帮助

```bash
./bin/generator --help
```

输出：
```
GinForge 脚手架工具 - 快速生成 CRUD 代码

这个工具可以帮助您：
  • 从数据库表自动生成完整的 CRUD 代码
  • 生成后端代码：Model、Repository、Service、Handler
  • 生成前端代码：API、Vue 列表页、表单页
  • 支持自定义模板和配置文件
  • 大幅提升开发效率
```

### 3. 列出数据库表

```bash
./bin/generator list:tables
```

输出示例：
```
找到 12 个表:
  1. admin_users
  2. admin_roles
  3. admin_menus
  ...
```

### 4. 生成 CRUD 代码

#### 方式 1：一键生成（推荐）⭐

```bash
# 生成代码并自动注册路由和菜单
./bin/generator gen:crud --table=articles --module=admin --auto-register

# 或使用简写
./bin/generator gen:crud --table=articles --module=admin -a
```

**这是最快的方式！** 生成器会自动完成：
- ✅ 生成所有代码
- ✅ 注册后端路由
- ✅ 注册前端路由
- ✅ 注册菜单
- ✅ 导入图标

**只需重启服务即可使用！**

#### 方式 2：预览模式

```bash
./bin/generator gen:crud --table=articles --module=admin -a --dry-run
```

先预览会生成哪些文件和注册哪些路由，确认无误后再正式生成。

#### 方式 3：使用配置文件

```bash
# 1. 生成配置文件
./bin/generator init:config --table=articles

# 2. 编辑配置文件 generator/articles.yaml

# 3. 从配置文件生成并自动注册
./bin/generator gen:crud --config=generator/articles.yaml -a
```

---

## 📊 生成结果

### 使用自动注册（-a）

生成完成后，您会看到类似输出：

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

📌 后续步骤:
  ✅ 后端路由已自动注册
  ✅ 前端路由已自动注册
  ✅ 菜单已自动注册

  🚀 现在只需重启服务即可使用！
     后端: cd services/admin-api && go run cmd/server/main.go
     前端: 刷新浏览器
```

### 不使用自动注册

如果不使用 `-a` 选项，需要手动完成以下步骤：

```
📌 后续步骤:
  1. 在路由文件中注册新的 Handler
  2. 在前端路由中添加新页面
  3. 在菜单中添加入口
  4. 重启服务并测试功能

💡 提示: 使用 --auto-register 或 -a 选项可以自动完成上述步骤
```

---

## 💡 核心优势

### ⚡ 效率提升 12-14 倍

传统手写一个 CRUD 模块（包含前后端）需要 **6-7 小时**，使用代码生成器（自动注册）只需 **30 分钟**！

| 任务 | 传统手写 | 手动注册 | 自动注册 ⭐ |
|------|---------|---------|------------|
| Model | 30 分钟 | 10 秒 | 10 秒 |
| Repository | 45 分钟 | 10 秒 | 10 秒 |
| Service | 60 分钟 | 10 秒 | 10 秒 |
| Handler | 45 分钟 | 10 秒 | 10 秒 |
| 前端 API | 20 分钟 | 10 秒 | 10 秒 |
| 前端列表 | 90 分钟 | 10 秒 | 10 秒 |
| 前端表单 | 60 分钟 | 10 秒 | 10 秒 |
| 注册路由和菜单 | 30 分钟 | 10 分钟 | **自动** |
| 测试调试 | 30 分钟 | 30 分钟 | 30 分钟 |
| **总计** | **6-7 小时** | **45 分钟** | **30 分钟** |
| **效率提升** | **基准** | **8-9x** | **12-14x** |

### 🎯 代码质量

- ✅ 统一的代码风格
- ✅ 完整的错误处理
- ✅ 规范的 API 设计
- ✅ 符合最佳实践
- ✅ 包含 Swagger 注释

### 🔧 高度可定制

- ✅ 支持自定义配置文件
- ✅ 支持自定义模板
- ✅ 灵活的字段配置
- ✅ 丰富的验证规则

---

## 📖 可用命令

### gen:crud

生成完整的 CRUD 代码（后端 + 前端）

```bash
./bin/generator gen:crud --table=<表名> --module=<模块>
```

**选项**：
- `--table, -t` - 数据库表名（必填）
- `--module, -m` - 模块名称（默认：admin）
- `--config, -c` - 配置文件路径
- `--frontend` - 是否生成前端代码（默认：true）
- `--force, -f` - 强制覆盖已存在的文件
- `--dry-run` - 预览模式，不实际创建文件
- `--verbose, -v` - 显示详细输出

### gen:model

只生成 Model 数据模型

```bash
./bin/generator gen:model --table=<表名> --module=<模块>
```

### init:config

生成配置文件模板

```bash
./bin/generator init:config --table=<表名>
```

### list:tables

列出数据库中的所有表

```bash
./bin/generator list:tables
```

---

## 🎨 智能特性详解

### 1. 字段类型自动映射

生成器会自动将 MySQL 类型映射到 Go 和 TypeScript 类型：

| MySQL 类型 | Go 类型 | TypeScript 类型 | 表单类型 |
|-----------|---------|-----------------|---------|
| int | int | number | input |
| bigint | int64 | number | input |
| varchar | string | string | input |
| text | string | string | textarea |
| datetime | time.Time | string | datetime |
| tinyint(1) | bool | boolean | switch |

### 2. 表单类型智能识别

根据字段名自动识别最合适的表单类型：

- `password` → 密码输入框
- `email` → 邮箱输入框
- `content`, `description` → 多行文本框
- `is_*` → 开关（switch）
- `status`, `type`, `category` → 下拉选择
- `*_at` → 日期时间选择器

### 3. 验证规则自动生成

根据字段属性自动生成验证规则：

- 非空字段 → `required`
- email 字段 → `required, email`
- phone 字段 → `required, len:11`
- password 字段 → `required, min:6`
- varchar(255) → `required, max:255`

### 4. 软删除支持

如果表有 `deleted_at` 字段，自动启用软删除功能：

```go
// 自动生成软删除方法
func (r *Repository) Delete(id uint64) error {
    return r.db.Delete(&Model{}, id).Error
}

func (r *Repository) ForceDelete(id uint64) error {
    return r.db.Unscoped().Delete(&Model{}, id).Error
}

func (r *Repository) Restore(id uint64) error {
    return r.db.Model(&Model{}).Unscoped().
        Where("id = ?", id).
        Update("deleted_at", nil).Error
}
```

---

## 🎯 下一步

- 📖 [详细使用指南](./usage) - 了解所有命令和选项
- ⚙️ [配置文件详解](./configuration) - 自定义生成规则
- 💼 [实战示例](./examples) - 完整的业务模块开发

**开始使用代码生成器，让开发效率提升 10 倍！** 🚀


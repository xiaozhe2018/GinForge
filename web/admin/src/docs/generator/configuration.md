# 配置文件详解

## 📋 配置文件概述

配置文件使用 YAML 格式，可以精确控制代码生成的各个方面。

### 生成配置文件

```bash
./bin/generator init:config --table=articles
```

这会创建 `generator/articles.yaml` 文件。

---

## 🔧 配置文件结构

### 完整配置示例

```yaml
# ========== 基础配置 ==========
table: articles             # 数据库表名
module: admin              # 模块名称（admin/user/file）
model_name: Article        # 模型名称（PascalCase）
resource_name: articles    # 资源名称（复数形式，用于 URL）

# ========== 字段配置 ==========
fields:
  - name: id
    type: bigint unsigned
    go_type: uint64
    ts_type: number
    nullable: false
    is_primary_key: true
    auto_increment: true
    comment: "文章ID"
    validations: []
    label: "ID"
    form_type: input
    list_visible: true
    form_visible: false
    searchable: false
    sortable: true

# ========== 功能特性 ==========
features:
  soft_delete: true       # 启用软删除
  timestamps: true        # 启用时间戳
  pagination: true        # 启用分页
  search: true           # 启用搜索
  sort: true             # 启用排序
  export: false          # 暂不支持导出
  import: false          # 暂不支持导入
  batch_delete: false    # 暂不支持批量删除

# ========== 前端配置 ==========
frontend:
  title: "文章管理"       # 页面标题
  icon: "Document"        # 菜单图标（Element Plus Icon）
  show_in_menu: true     # 是否显示在菜单
  menu_parent: ""        # 父菜单（可选）

# ========== 生成选项 ==========
options:
  output_dir: "."        # 输出目录
  with_frontend: true    # 生成前端代码
  force: false           # 是否强制覆盖
  dry_run: false         # 是否预览模式
  verbose: false         # 是否详细输出
```

---

## 📝 字段配置详解

每个字段支持以下配置：

### 基本属性

| 属性 | 类型 | 说明 | 示例 |
|------|------|------|------|
| name | string | 字段名（snake_case） | `title` |
| type | string | 数据库类型 | `varchar(255)` |
| go_type | string | Go 类型 | `string` |
| ts_type | string | TypeScript 类型 | `string` |
| nullable | boolean | 是否可为空 | `false` |
| is_primary_key | boolean | 是否主键 | `false` |
| auto_increment | boolean | 是否自增 | `false` |
| default_value | string | 默认值 | `""` |
| comment | string | 字段注释 | `"文章标题"` |

### 验证规则

```yaml
validations:
  - required      # 必填
  - email         # 邮箱格式
  - min:6         # 最小长度
  - max:255       # 最大长度
  - len:11        # 固定长度
  - url           # URL 格式
```

**可用的验证规则：**

| 规则 | 说明 | 示例 |
|------|------|------|
| required | 必填 | `validations: [required]` |
| email | 邮箱格式 | `validations: [required, email]` |
| min:N | 最小长度 | `validations: [required, min:6]` |
| max:N | 最大长度 | `validations: [required, max:255]` |
| len:N | 固定长度 | `validations: [required, len:11]` |
| url | URL 格式 | `validations: [url]` |

### UI 配置

```yaml
label: "标题"              # 显示标签
form_type: input          # 表单类型
list_visible: true        # 列表中是否显示
form_visible: true        # 表单中是否显示
searchable: true          # 是否可搜索
sortable: true            # 是否可排序
```

**表单类型 (form_type)：**

| 类型 | 说明 | 适用场景 |
|------|------|---------|
| input | 单行输入框 | 标题、名称、短文本 |
| textarea | 多行输入框 | 描述、备注、长文本 |
| password | 密码输入框 | 密码字段 |
| email | 邮箱输入框 | 邮箱字段 |
| number | 数字输入框 | 年龄、数量、价格 |
| switch | 开关 | 状态、是否启用 |
| select | 下拉选择 | 分类、类型、状态 |
| date | 日期选择器 | 出生日期、日期 |
| datetime | 日期时间选择器 | 创建时间、发布时间 |
| upload | 文件上传 | 头像、图片、文件 |
| editor | 富文本编辑器 | 文章内容、详情 |

---

## 🎨 字段配置示例

### 示例 1：标题字段

```yaml
- name: title
  type: varchar(255)
  go_type: string
  ts_type: string
  nullable: false
  comment: "文章标题"
  validations:
    - required
    - max:255
  label: "标题"
  form_type: input
  list_visible: true
  form_visible: true
  searchable: true
  sortable: true
```

### 示例 2：内容字段（富文本）

```yaml
- name: content
  type: text
  go_type: string
  ts_type: string
  nullable: false
  comment: "文章内容"
  validations:
    - required
  label: "内容"
  form_type: editor         # 富文本编辑器
  list_visible: false       # 列表中不显示
  form_visible: true
  searchable: true
  sortable: false
```

### 示例 3：状态字段（下拉选择）

```yaml
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
```

### 示例 4：是否置顶（开关）

```yaml
- name: is_top
  type: tinyint(1)
  go_type: bool
  ts_type: boolean
  nullable: false
  default_value: "0"
  comment: "是否置顶"
  validations: []
  label: "置顶"
  form_type: switch
  list_visible: true
  form_visible: true
  searchable: false
  sortable: true
```

### 示例 5：浏览次数（只显示）

```yaml
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
  form_visible: false       # 表单中不显示
  searchable: false
  sortable: true
```

### 示例 6：创建时间（自动管理）

```yaml
- name: created_at
  type: datetime
  go_type: "*time.Time"
  ts_type: string
  nullable: true
  comment: "创建时间"
  validations: []
  label: "创建时间"
  form_type: datetime
  list_visible: true
  form_visible: false       # 表单中不显示（自动填充）
  searchable: false
  sortable: true
```

---

## 🚀 功能特性配置

### soft_delete - 软删除

```yaml
features:
  soft_delete: true
```

启用后，会自动生成软删除相关方法：

```go
// Delete - 软删除
func (r *Repository) Delete(id uint64) error {
    return r.db.Delete(&Model{}, id).Error
}

// ForceDelete - 永久删除
func (r *Repository) ForceDelete(id uint64) error {
    return r.db.Unscoped().Delete(&Model{}, id).Error
}

// Restore - 恢复已删除
func (r *Repository) Restore(id uint64) error {
    return r.db.Model(&Model{}).Unscoped().
        Where("id = ?", id).
        Update("deleted_at", nil).Error
}
```

### timestamps - 时间戳

```yaml
features:
  timestamps: true
```

启用后，GORM 会自动管理 `created_at` 和 `updated_at` 字段。

### pagination - 分页

```yaml
features:
  pagination: true
```

启用后，会生成分页相关代码：

```go
// 分页
if req.Page > 0 && req.PageSize > 0 {
    offset := (req.Page - 1) * req.PageSize
    db = db.Offset(offset).Limit(req.PageSize)
}
```

### search - 搜索

```yaml
features:
  search: true
```

启用后，会生成搜索相关代码：

```go
// 搜索
if req.Keyword != "" {
    db = db.Where("title LIKE ? OR content LIKE ?", 
        "%"+req.Keyword+"%", 
        "%"+req.Keyword+"%")
}
```

### sort - 排序

```yaml
features:
  sort: true
```

启用后，会生成排序相关代码：

```go
// 排序
if req.SortBy != "" {
    order := req.SortBy
    if req.SortOrder == "desc" {
        order += " DESC"
    }
    db = db.Order(order)
} else {
    db = db.Order("id DESC")
}
```

---

## 🎨 前端配置

### title - 页面标题

```yaml
frontend:
  title: "文章管理"
```

用于：
- 页面标题
- 表格标题
- 对话框标题

### icon - 菜单图标

```yaml
frontend:
  icon: "Document"
```

可用的 Element Plus 图标：
- `Document` - 文档
- `User` - 用户
- `Setting` - 设置
- `Files` - 文件
- `Menu` - 菜单
- `Lock` - 权限
- `Position` - 分类
- `PriceTag` - 标签
- `ChatDotRound` - 评论
- 更多图标参考：[Element Plus Icons](https://element-plus.org/zh-CN/component/icon.html)

### show_in_menu - 是否显示在菜单

```yaml
frontend:
  show_in_menu: true
```

如果设置为 `false`，则不会在菜单中显示，但仍可通过路由访问。

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

## 💡 配置技巧

### 1. 继承默认配置

生成配置文件时，会自动填充所有字段的默认配置。您只需要修改需要自定义的部分。

### 2. 批量修改字段属性

使用文本编辑器的查找替换功能，可以批量修改字段属性：

```yaml
# 批量设置所有 varchar 字段为可搜索
查找：   type: varchar
替换为： type: varchar\n    searchable: true
```

### 3. 复用配置

对于相似的表，可以复制配置文件并修改：

```bash
cp generator/articles.yaml generator/news.yaml
# 编辑 generator/news.yaml
```

### 4. 版本控制

将配置文件加入版本控制，方便团队协作和历史追溯：

```bash
git add generator/
git commit -m "Add generator configs"
```

---

## 📖 下一步

- 💼 [实战示例](./examples) - 完整的业务模块开发流程

**掌握配置文件，让代码生成更加灵活！** 🚀


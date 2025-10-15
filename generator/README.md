# CRUD 代码生成器配置文件

这个目录用于存放 CRUD 代码生成器的配置文件。

## 📁 文件说明

- `example.yaml` - 完整的配置文件示例
- 其他 `.yaml` 文件 - 您自定义的配置文件

## 🚀 快速开始

### 1. 生成配置文件

```bash
go run cmd/generator/main.go init:config --table=your_table_name
```

这会在当前目录创建 `your_table_name.yaml` 配置文件。

### 2. 编辑配置文件

根据需要修改字段配置、验证规则、表单类型等。

### 3. 生成代码

```bash
go run cmd/generator/main.go gen:crud --config=generator/your_table_name.yaml
```

## 📖 配置文件结构

```yaml
# 基础配置
table: table_name          # 数据库表名
module: module_name        # 模块名（admin/user/file）
model_name: ModelName      # 模型名（PascalCase）
resource_name: resources   # 资源名（复数，URL用）

# 字段配置
fields:
  - name: field_name       # 字段名
    type: varchar(255)     # 数据库类型
    go_type: string        # Go 类型
    ts_type: string        # TS 类型
    label: "字段标签"       # 显示标签
    form_type: input       # 表单类型
    validations:           # 验证规则
      - required
      - max:255
    list_visible: true     # 列表中显示
    form_visible: true     # 表单中显示
    searchable: true       # 可搜索
    sortable: true         # 可排序

# 功能特性
features:
  soft_delete: true        # 软删除
  timestamps: true         # 时间戳
  pagination: true         # 分页
  search: true            # 搜索
  sort: true              # 排序

# 前端配置
frontend:
  title: "页面标题"
  icon: "Document"
  show_in_menu: true
```

## 💡 提示

1. 参考 `example.yaml` 了解完整配置
2. 使用 `--dry-run` 选项预览生成结果
3. 生成后的代码可以自由修改
4. 将配置文件加入版本控制

## 📚 更多信息

查看完整文档：[GENERATOR_GUIDE.md](../docs/GENERATOR_GUIDE.md)


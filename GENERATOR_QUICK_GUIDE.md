# 🚀 GinForge 代码生成器快速使用指南

## ⚡ 5 分钟快速体验

### 1. 编译生成器（首次使用）

```bash
cd /Users/chaojidoudou/project/go/GinForge
go build -o bin/generator ./cmd/generator
```

### 2. 查看数据库表

```bash
./bin/generator list:tables
```

### 3. 一键生成 CRUD 模块

```bash
# 推荐：使用自动注册（-a）
./bin/generator gen:crud --table=articles --module=admin -a

# 或先预览
./bin/generator gen:crud --table=articles --module=admin -a --dry-run
```

### 4. 重启服务

```bash
# 后端
cd services/admin-api && go run cmd/server/main.go

# 前端
# 刷新浏览器
```

### 5. 访问新功能

访问 `http://localhost:3000/dashboard/articles`

**完成！✅**

---

## 📖 常用命令

### 列出所有表

```bash
./bin/generator list:tables
```

### 生成 CRUD（推荐）⭐

```bash
./bin/generator gen:crud --table=<表名> --module=admin -a
```

### 生成配置文件

```bash
./bin/generator init:config --table=<表名>
```

### 从配置文件生成

```bash
./bin/generator gen:crud --config=generator/<表名>.yaml -a
```

### 预览模式

```bash
./bin/generator gen:crud --table=<表名> --module=admin -a --dry-run
```

---

## 💡 核心选项

| 选项 | 简写 | 说明 |
|------|------|------|
| --table | -t | 数据库表名 |
| --module | -m | 模块名称（admin/user/file） |
| --auto-register | -a | **自动注册路由和菜单** ⭐ |
| --config | -c | 配置文件路径 |
| --force | -f | 强制覆盖已存在的文件 |
| --dry-run | - | 预览模式，不实际创建文件 |
| --verbose | -v | 显示详细输出 |
| --frontend | - | 生成前端代码（默认：true） |

---

## 🎯 使用场景

### 场景 1：快速原型（5-10 分钟）

```bash
# 一键生成
./bin/generator gen:crud --table=articles --module=admin -a

# 重启服务
# 完成！
```

**适用于**：标准 CRUD，快速演示

### 场景 2：标准业务（30 分钟）

```bash
# 1. 生成代码
./bin/generator gen:crud --table=articles --module=admin -a

# 2. 自定义前端表单（添加下拉选项等）
# 3. 测试功能
# 完成！
```

**适用于**：常规业务模块

### 场景 3：复杂定制（80 分钟）

```bash
# 1. 生成配置文件
./bin/generator init:config --table=articles

# 2. 编辑配置文件
# 3. 生成代码
./bin/generator gen:crud --config=generator/articles.yaml -a

# 4. 自定义业务逻辑
# 5. 扩展功能
# 6. 测试调试
# 完成！
```

**适用于**：复杂业务需求

---

## 📊 效率对比

| 开发模式 | 耗时 | 效率提升 |
|---------|------|---------|
| 传统手写 | 6-7 小时 | 基准 |
| 生成器（手动注册） | 45 分钟 | 8-9 倍 |
| **生成器（自动注册）** | **30 分钟** | **12-14 倍** ⚡ |
| **快速模式** | **5-10 分钟** | **40-80 倍** ⚡⚡ |

---

## 📚 完整文档

### 在线教程

访问 `http://localhost:3000` → 文档中心 → CRUD 代码生成器

### 本地文档

- **完整指南**：`docs/GENERATOR_GUIDE.md`
- **自动注册**：`docs/GENERATOR_AUTO_REGISTER.md`
- **快速上手**：`generator/QUICK_START.md`
- **配置示例**：`generator/example.yaml`

---

## 🎉 总结

使用 GinForge 代码生成器，您可以：

✅ **5-10 分钟** 完成一个标准 CRUD 模块
✅ **30 分钟** 完成一个常规业务模块
✅ **80 分钟** 完成一个复杂定制模块

**效率提升 12-14 倍，快速模式高达 40-80 倍！**

**现在就开始使用，让开发效率飞跃！** 🚀

---

**快速参考**：

```bash
# 编译
go build -o bin/generator ./cmd/generator

# 列表
./bin/generator list:tables

# 一键生成（推荐）
./bin/generator gen:crud --table=<表名> --module=admin -a

# 预览
./bin/generator gen:crud --table=<表名> --module=admin -a --dry-run

# 帮助
./bin/generator --help
```

**GinForge - 真正的一键生成，开箱即用！** 🎊


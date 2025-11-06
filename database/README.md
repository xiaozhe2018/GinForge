# GinForge 数据库管理

## 目录结构

```
database/
├── README.md                    # 本文件
└── migrations/                  # 数据库迁移文件
    └── init.sql                # 统一初始化脚本（包含所有表结构）
```

**注意：** 所有表都使用 `gf_` 前缀（如 `gf_admin_users`、`gf_articles`）

## 快速开始

### 前置条件

1. **安装 MySQL** 并启动服务
2. **配置环境变量**（推荐）或修改 `configs/config.yaml`
   ```bash
   # 复制 env.example 为 .env
   cp env.example .env
   
   # 编辑 .env，配置数据库信息
   GINFORGE_DATABASE_HOST=localhost
   GINFORGE_DATABASE_PORT=3306
   GINFORGE_DATABASE_DATABASE=gin_forge
   GINFORGE_DATABASE_USERNAME=root
   GINFORGE_DATABASE_PASSWORD=your_password
   ```
3. **创建数据库**
   ```bash
   mysql -uroot -p -e "CREATE DATABASE IF NOT EXISTS gin_forge DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   ```

### 方式1：使用 Makefile（推荐）

```bash
# 初始化数据库（推荐使用）
make init

# 或者使用完整命令
make db-init

# 重置数据库（删除并重新创建）
make db-reset

# 查看数据库状态
make db-status
```

### 方式2：手动执行

```bash
# 1. 确保 MySQL 已启动并配置好环境变量或 configs/config.yaml

# 2. 创建数据库
mysql -uroot -p123456 -e "CREATE DATABASE IF NOT EXISTS gin_forge DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 3. 执行初始化脚本
mysql -uroot -p123456 gin_forge < database/migrations/init.sql
```

## 数据库表说明

### 管理后台相关表（前缀：gf_）
- `gf_admin_users` - 管理员用户表
- `gf_admin_roles` - 角色表
- `gf_admin_permissions` - 权限表
- `gf_admin_menus` - 菜单表
- `gf_admin_user_roles` - 用户角色关联表
- `gf_admin_role_permissions` - 角色权限关联表
- `gf_admin_role_menus` - 角色菜单关联表

### 系统配置相关表（前缀：gf_）
- `gf_admin_system_configs` - 系统配置表
- `gf_admin_operation_logs` - 操作日志表

### 业务表（前缀：gf_）
- `gf_articles` - 文章表（示例）

**默认数据：**
- 默认管理员账号：`admin` / `admin123`
- 超级管理员角色：`super_admin`
- 18个默认权限
- 6个默认菜单（一级）+ 6个二级菜单
- 系统配置数据
- 3条示例文章数据

## 配置说明

数据库配置可以通过以下方式设置：

### 方式1：环境变量（推荐）

复制 `env.example` 为 `.env` 并修改：

```bash
GINFORGE_DATABASE_DRIVER=mysql
GINFORGE_DATABASE_HOST=localhost
GINFORGE_DATABASE_PORT=3306
GINFORGE_DATABASE_DATABASE=gin_forge
GINFORGE_DATABASE_USERNAME=root
GINFORGE_DATABASE_PASSWORD=123456
GINFORGE_DATABASE_CHARSET=utf8mb4
GINFORGE_DATABASE_TIMEZONE=Asia/Shanghai
```

### 方式2：配置文件

在 `configs/config.yaml` 中配置：

```yaml
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
  charset: "utf8mb4"
  timezone: "Asia/Shanghai"
```

**注意：** 环境变量优先级高于配置文件

## 添加新的表

1. 在 `database/migrations/init.sql` 中添加新的表结构
2. **重要：** 所有表名必须使用 `gf_` 前缀
3. 执行 `make db-init` 或 `make db-reset` 应用更改

示例：

```sql
-- 在 init.sql 中添加新表
CREATE TABLE IF NOT EXISTS `gf_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
```

**表命名规范：**
- 所有表名必须以 `gf_` 开头
- 使用小写字母和下划线
- 例如：`gf_products`、`gf_orders`、`gf_user_profiles`

## 常见问题

### Q: 执行失败，提示表已存在？
A: 使用 `make db-reset` 重置数据库，或手动删除表后重新执行。

### Q: 如何修改默认管理员密码？
A: 修改 `001_create_admin_tables.sql` 中的密码哈希值，或登录后通过管理后台修改。

### Q: 如何只执行某个迁移文件？
A: 直接执行对应的 SQL 文件：
```bash
mysql -uroot -p123456 gin_forge < database/migrations/001_create_admin_tables.sql
```

## 注意事项

1. **执行顺序**：迁移文件必须按数字顺序执行
2. **数据备份**：生产环境执行前请先备份数据
3. **外键约束**：某些迁移文件可能依赖其他表，注意执行顺序
4. **字符集**：确保使用 `utf8mb4` 字符集以支持 emoji 等特殊字符

## 默认账号

- **用户名**：`admin`
- **密码**：`admin123`
- **邮箱**：`admin@ginforge.com`

**注意：生产环境请务必修改默认密码！**


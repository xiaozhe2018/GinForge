# 🚀 快速上手指南

## 5 分钟快速体验

### 1. 编译生成器（首次使用）

```bash
cd /Users/chaojidoudou/project/go/GinForge
go build -o bin/generator ./cmd/generator
```

### 2. 查看可用的数据库表

```bash
./bin/generator list:tables
```

输出：
```
找到 12 个表:
  1. admin_menus
  2. admin_users
  3. admin_roles
  ...
```

### 3. 选择一个表，预览生成结果

```bash
./bin/generator gen:crud --table=admin_menus --module=admin --dry-run
```

查看会生成哪些文件，确认无误后继续。

### 4. 正式生成代码

```bash
./bin/generator gen:crud --table=admin_menus --module=admin
```

生成的文件：
```
✅ services/admin-api/internal/model/menus.go
✅ services/admin-api/internal/repository/menus_repository.go
✅ services/admin-api/internal/service/menus_service.go
✅ services/admin-api/internal/handler/menus_handler.go
✅ web/admin/src/api/menus.ts
✅ web/admin/src/views/Menus/index.vue
✅ web/admin/src/views/Menus/Form.vue
```

### 5. 注册路由（按提示操作）

根据命令行的提示，在 `services/admin-api/internal/router/router.go` 中添加：

```go
// 初始化
menusRepo := repository.NewMenusRepository(database)
menusService := service.NewMenusService(menusRepo, log)
menusHandler := handler.NewMenusHandler(menusService, log)

// 注册路由
auth.GET("/menuses", menusHandler.List)
auth.GET("/menuses/:id", menusHandler.Get)
auth.POST("/menuses", menusHandler.Create)
auth.PUT("/menuses/:id", menusHandler.Update)
auth.DELETE("/menuses/:id", menusHandler.Delete)
```

### 6. 添加前端路由

在 `web/admin/src/router/index.ts` 中添加：

```typescript
{
  path: 'menuses',
  name: 'MenusList',
  component: () => import('@/views/Menus/index.vue'),
  meta: { title: '菜单管理', requiresAuth: true }
}
```

### 7. 添加菜单入口

在 `web/admin/src/layout/index.vue` 中添加：

```vue
<el-menu-item index="/dashboard/menuses">
  <el-icon><Document /></el-icon>
  <span>菜单管理</span>
</el-menu-item>
```

### 8. 重启服务并测试

```bash
# 重启后端
cd services/admin-api
go run cmd/server/main.go

# 刷新前端
```

访问 `http://localhost:3000/dashboard/menuses`，查看生成的页面！

---

## 进阶使用

### 使用配置文件自定义生成

#### 1. 生成配置文件

```bash
./bin/generator init:config --table=articles
```

#### 2. 编辑配置文件

编辑 `generator/articles.yaml`：

```yaml
# 修改字段配置
fields:
  - name: content
    form_type: editor      # 改为富文本编辑器
  
  - name: category_id
    form_type: select      # 改为下拉选择
  
  - name: view_count
    list_visible: true     # 在列表中显示
    form_visible: false    # 在表单中隐藏
```

#### 3. 从配置文件生成

```bash
./bin/generator gen:crud --config=generator/articles.yaml
```

---

## 常用命令

### 生成 CRUD

```bash
# 基础用法
./bin/generator gen:crud --table=articles --module=admin

# 只生成后端
./bin/generator gen:crud --table=articles --module=admin --frontend=false

# 强制覆盖
./bin/generator gen:crud --table=articles --module=admin --force

# 预览模式
./bin/generator gen:crud --table=articles --module=admin --dry-run

# 详细输出
./bin/generator gen:crud --table=articles --module=admin --verbose
```

### 只生成 Model

```bash
./bin/generator gen:model --table=articles --module=admin
```

### 生成配置文件

```bash
./bin/generator init:config --table=articles
```

### 列出数据库表

```bash
./bin/generator list:tables
```

---

## 小技巧

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

---

## 下一步

- 📖 阅读[完整文档](../docs/GENERATOR_GUIDE.md)
- 🎨 查看[配置示例](./example.yaml)
- 💡 学习[最佳实践](../docs/GENERATOR_GUIDE.md#最佳实践)

**祝您开发愉快！** 🎉


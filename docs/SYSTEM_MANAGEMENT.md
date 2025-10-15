# 系统管理功能使用指南

## 1. 概述

系统管理模块提供了对GinForge框架的全局配置、监控和日志管理功能。通过该模块，管理员可以方便地配置系统参数、监控系统运行状态、查看操作日志等。

## 2. 系统配置

系统配置分为以下几个部分：

### 2.1 基本配置

基本配置包括系统名称、版本、描述、Logo和默认语言等。

**API接口：**

```
GET /api/v1/admin/system/configs?group=basic  // 获取基本配置列表
PUT /api/v1/admin/system/configs              // 更新配置
```

更新配置请求示例：

```json
{
  "key": "system.name",
  "value": "GinForge 管理系统"
}
```

### 2.2 安全配置

安全配置包括密码策略、登录失败处理和会话超时等。

**API接口：**

```
GET /api/v1/admin/system/configs?group=security  // 获取安全配置列表
PUT /api/v1/admin/system/configs                 // 更新配置
```

更新密码最小长度示例：

```json
{
  "key": "security.min_password_length",
  "value": "8"
}
```

### 2.3 邮件配置

邮件配置包括SMTP服务器、端口、发件人邮箱等。

**API接口：**

```
GET /api/v1/admin/system/configs?group=email  // 获取邮件配置列表
PUT /api/v1/admin/system/configs              // 更新配置
POST /api/v1/admin/system/email/test          // 测试邮件发送
```

测试邮件发送请求示例：

```json
{
  "email": "test@example.com"
}
```

### 2.4 缓存配置

缓存配置包括缓存类型、Redis连接参数和过期时间等。

**API接口：**

```
GET /api/v1/admin/system/configs?group=cache  // 获取缓存配置列表
PUT /api/v1/admin/system/configs              // 更新配置
POST /api/v1/admin/system/cache/test          // 测试缓存连接
POST /api/v1/admin/system/cache/clear         // 清空缓存
```

## 3. 系统监控

系统监控提供了对系统运行状态的实时监控，包括在线用户数、CPU使用率、内存使用率、磁盘使用率等。

**API接口：**

```
GET /api/v1/admin/system/info      // 获取系统信息
GET /api/v1/admin/system/runtime   // 获取运行时信息
GET /api/v1/admin/system/health    // 系统健康检查
```

系统信息响应示例：

```json
{
  "online_users": 10,
  "cpu_usage": 25,
  "memory_usage": 40,
  "disk_usage": 30,
  "network_in": 1024000,
  "network_out": 512000,
  "uptime": "3天5小时20分钟",
  "version": "1.0.0",
  "environment": "production"
}
```

## 4. 日志管理

日志管理提供了对系统操作日志的查询和管理功能。

**API接口：**

```
GET /api/v1/admin/system/logs       // 获取日志列表
POST /api/v1/admin/system/logs/clear // 清空日志
```

日志查询参数：

- `page`: 页码
- `page_size`: 每页数量
- `user_id`: 用户ID
- `username`: 用户名
- `method`: 请求方法
- `path`: 请求路径
- `ip`: IP地址
- `status_code`: 状态码
- `start_time`: 开始时间
- `end_time`: 结束时间

## 5. 数据库表结构

### 5.1 系统配置表

```sql
CREATE TABLE `admin_system_configs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL COMMENT '配置键',
  `value` text DEFAULT NULL COMMENT '配置值',
  `type` varchar(20) DEFAULT 'string' COMMENT '配置类型:string,number,boolean,json',
  `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
  `group` varchar(50) DEFAULT 'default' COMMENT '配置分组',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_group` (`group`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';
```

### 5.2 操作日志表

```sql
CREATE TABLE `admin_operation_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `path` varchar(255) NOT NULL COMMENT '请求路径',
  `ip` varchar(45) DEFAULT NULL COMMENT '请求IP',
  `user_agent` text DEFAULT NULL COMMENT '用户代理',
  `request_data` json DEFAULT NULL COMMENT '请求数据',
  `response_data` json DEFAULT NULL COMMENT '响应数据',
  `status_code` int(11) DEFAULT 200 COMMENT '响应状态码',
  `duration` int(11) DEFAULT 0 COMMENT '请求耗时(毫秒)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_method` (`method`),
  KEY `idx_path` (`path`),
  KEY `idx_ip` (`ip`),
  KEY `idx_status_code` (`status_code`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';
```

## 6. 前端集成

### 6.1 获取系统配置

```typescript
// 获取系统配置
const getSystemConfig = async (group: string = '') => {
  try {
    const params = group ? { group } : {};
    const response = await axios.get('/api/v1/admin/system/configs', { params });
    return response.data.data;
  } catch (error) {
    console.error('获取系统配置失败:', error);
    throw error;
  }
};

// 更新系统配置
const updateSystemConfig = async (key: string, value: string) => {
  try {
    const response = await axios.put('/api/v1/admin/system/configs', { key, value });
    return response.data;
  } catch (error) {
    console.error('更新系统配置失败:', error);
    throw error;
  }
};
```

### 6.2 系统监控组件

```vue
<template>
  <div class="system-monitor">
    <el-row :gutter="20">
      <el-col :span="6" v-for="(item, index) in monitorItems" :key="index">
        <el-card class="monitor-card">
          <div class="card-content">
            <div class="card-icon" :style="{ backgroundColor: item.color }">
              <el-icon size="24" color="#fff"><component :is="item.icon" /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">{{ item.title }}</div>
              <div class="card-value">{{ item.value }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { User, Monitor, Coin, Warning } from '@element-plus/icons-vue';
import axios from 'axios';

const monitorItems = ref([
  { title: '在线用户', value: 0, icon: User, color: '#409EFF' },
  { title: 'CPU使用率', value: '0%', icon: Monitor, color: '#67C23A' },
  { title: '内存使用率', value: '0%', icon: Coin, color: '#E6A23C' },
  { title: '错误日志', value: 0, icon: Warning, color: '#F56C6C' }
]);

let timer: number;

const fetchSystemInfo = async () => {
  try {
    const response = await axios.get('/api/v1/admin/system/info');
    const data = response.data.data;
    
    monitorItems.value[0].value = data.online_users;
    monitorItems.value[1].value = `${data.cpu_usage}%`;
    monitorItems.value[2].value = `${data.memory_usage}%`;
    monitorItems.value[3].value = data.error_count;
  } catch (error) {
    console.error('获取系统信息失败:', error);
  }
};

onMounted(() => {
  fetchSystemInfo();
  timer = window.setInterval(fetchSystemInfo, 30000);
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>
```

## 7. 最佳实践

1. **定期备份配置**：系统配置是整个应用的核心，建议定期备份配置数据。

2. **监控阈值设置**：为系统监控设置合理的阈值，当超过阈值时发送告警通知。

3. **日志轮转**：设置日志轮转策略，避免日志表数据过多影响性能。

4. **敏感信息处理**：在操作日志中，对密码等敏感信息进行脱敏处理。

5. **权限控制**：只允许超级管理员访问系统管理功能，普通管理员只能查看不能修改。

## 8. 常见问题

### 8.1 配置更新后不生效

可能原因：
- 配置缓存未更新
- 应用需要重启才能加载新配置

解决方案：
- 清空缓存
- 重启应用服务

### 8.2 系统监控数据不准确

可能原因：
- 监控代理未正确配置
- 数据收集间隔过长

解决方案：
- 检查监控代理配置
- 调整数据收集间隔

### 8.3 日志查询性能问题

可能原因：
- 日志表数据量过大
- 查询条件未使用索引

解决方案：
- 定期归档历史日志
- 优化查询条件，确保使用索引字段

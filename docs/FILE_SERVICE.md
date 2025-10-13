# 文件上传微服务使用指南

## 1. 简介

文件上传微服务是 GinForge 框架的一个独立服务，提供文件上传、下载、管理等功能。该服务支持多种存储方式，包括本地存储、阿里云 OSS、AWS S3 和 MinIO 等。

### 1.1 核心特性

- **多种存储支持**：本地存储、阿里云 OSS、AWS S3、MinIO
- **完整的文件管理**：上传、下载、删除、查询
- **元数据管理**：文件信息、统计数据
- **高级功能**：文件秒传、断点续传
- **安全控制**：文件类型验证、大小限制
- **性能优化**：异步处理、缓存支持

## 2. 快速开始

### 2.1 启动服务

```bash
# 方式一：单独启动文件服务
go run ./services/file-api/cmd/server

# 方式二：通过 Makefile 启动所有服务
make run
```

文件服务默认在 `8086` 端口启动，可通过配置文件修改。

### 2.2 基本配置

配置文件位于 `configs/file-api.yaml`，主要配置项：

```yaml
storage:
  # 存储类型: local, oss, s3, minio
  type: local
  
  # 本地存储配置
  local:
    base_path: ./uploads
  
  # URL前缀（用于生成访问URL）
  url_prefix: http://localhost:8086/uploads
  
  # 最大文件大小（字节）- 100MB
  max_file_size: 104857600
```

## 3. API 使用示例

### 3.1 上传文件

```bash
curl -X POST http://localhost:8086/api/v1/files/upload \
  -H "Content-Type: multipart/form-data" \
  -F "file=@/path/to/your/file.jpg" \
  -F "description=示例图片" \
  -F "tags=示例,图片" \
  -F "user_id=1" \
  -F "user_type=admin"
```

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "original_name": "example.jpg",
    "file_name": "example_1642342342_a1b2c3d4.jpg",
    "file_size": 1024000,
    "file_type": "image",
    "mime_type": "image/jpeg",
    "url": "http://localhost:8086/uploads/images/example_1642342342_a1b2c3d4.jpg",
    "hash": "d41d8cd98f00b204e9800998ecf8427e",
    "upload_time": "2025-01-15T10:30:45Z"
  }
}
```

### 3.2 下载文件

```bash
curl -X GET http://localhost:8086/api/v1/files/1/download \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o downloaded_file.jpg
```

### 3.3 获取文件列表

```bash
curl -X GET "http://localhost:8086/api/v1/files?page=1&page_size=10&file_type=image" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "list": [
      {
        "id": 1,
        "original_name": "example1.jpg",
        "file_name": "example1_1642342342_a1b2c3d4.jpg",
        "file_size": 1024000,
        "file_type": "image",
        "mime_type": "image/jpeg",
        "url": "http://localhost:8086/uploads/images/example1_1642342342_a1b2c3d4.jpg",
        "hash": "d41d8cd98f00b204e9800998ecf8427e",
        "upload_time": "2025-01-15T10:30:45Z"
      },
      // 更多文件...
    ]
  }
}
```

### 3.4 删除文件

```bash
curl -X DELETE http://localhost:8086/api/v1/files/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3.5 获取文件统计信息

```bash
curl -X GET http://localhost:8086/api/v1/files/statistics \
  -H "Authorization: Bearer YOUR_TOKEN"
```

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_count": 156,
    "total_size": 256000000,
    "type_statistics": {
      "image": 78,
      "document": 45,
      "video": 12,
      "other": 21
    }
  }
}
```

## 4. 前端集成示例

### 4.1 使用 Axios 上传文件

```javascript
import axios from 'axios';

const uploadFile = async (file, description, tags, userId, userType) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('description', description);
  formData.append('tags', tags);
  formData.append('user_id', userId);
  formData.append('user_type', userType);
  
  try {
    const response = await axios.post(
      'http://localhost:8086/api/v1/files/upload',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${token}`
        }
      }
    );
    
    return response.data;
  } catch (error) {
    console.error('上传文件失败:', error);
    throw error;
  }
};
```

### 4.2 使用 Vue 3 组件

```vue
<template>
  <div>
    <input type="file" @change="handleFileChange" />
    <input v-model="description" placeholder="文件描述" />
    <input v-model="tags" placeholder="标签（逗号分隔）" />
    <button @click="uploadFile" :disabled="!selectedFile">上传</button>
    
    <div v-if="uploadResult">
      <h3>上传成功</h3>
      <p>文件名: {{ uploadResult.original_name }}</p>
      <p>文件大小: {{ formatSize(uploadResult.file_size) }}</p>
      <p>文件类型: {{ uploadResult.file_type }}</p>
      <img v-if="isImage(uploadResult.mime_type)" :src="uploadResult.url" />
      <a :href="uploadResult.url" target="_blank">查看文件</a>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';

const selectedFile = ref(null);
const description = ref('');
const tags = ref('');
const uploadResult = ref(null);

const handleFileChange = (event) => {
  selectedFile.value = event.target.files[0];
};

const uploadFile = async () => {
  if (!selectedFile.value) return;
  
  const formData = new FormData();
  formData.append('file', selectedFile.value);
  formData.append('description', description.value);
  formData.append('tags', tags.value);
  formData.append('user_id', '1'); // 从用户状态获取
  formData.append('user_type', 'admin'); // 从用户状态获取
  
  try {
    const response = await axios.post(
      'http://localhost:8086/api/v1/files/upload',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }
    );
    
    uploadResult.value = response.data.data;
  } catch (error) {
    console.error('上传文件失败:', error);
    alert('上传失败: ' + error.message);
  }
};

const formatSize = (size) => {
  if (size < 1024) return size + ' B';
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB';
  return (size / (1024 * 1024)).toFixed(2) + ' MB';
};

const isImage = (mimeType) => {
  return mimeType && mimeType.startsWith('image/');
};
</script>
```

## 5. 高级配置

### 5.1 切换到阿里云 OSS

修改 `configs/file-api.yaml`：

```yaml
storage:
  type: oss
  oss:
    endpoint: oss-cn-hangzhou.aliyuncs.com
    access_key_id: YOUR_ACCESS_KEY_ID
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    bucket_name: your-bucket-name
    region: cn-hangzhou
  url_prefix: https://your-bucket-name.oss-cn-hangzhou.aliyuncs.com
```

### 5.2 配置文件类型限制

```yaml
storage:
  # 允许的文件类型
  allowed_types:
    - image/jpeg
    - image/png
    - application/pdf
  
  # 允许的文件扩展名
  allowed_exts:
    - .jpg
    - .jpeg
    - .png
    - .pdf
```

## 6. 性能优化建议

1. **使用对象存储**：生产环境建议使用对象存储而非本地存储
2. **配置 CDN**：对于图片等静态资源，配置 CDN 加速
3. **启用文件秒传**：通过文件哈希检测重复文件
4. **合理设置文件大小限制**：根据业务需求设置合理的文件大小限制
5. **定期清理临时文件**：设置定时任务清理临时文件和无效文件

## 7. 常见问题

### 7.1 文件上传失败

可能原因：
- 文件大小超过限制
- 文件类型不被允许
- 存储路径权限问题
- 存储空间不足

解决方案：
- 检查配置中的 `max_file_size` 设置
- 检查 `allowed_types` 和 `allowed_exts` 设置
- 确保存储目录有正确的读写权限
- 检查磁盘空间或云存储容量

### 7.2 无法访问上传的文件

可能原因：
- URL 前缀配置错误
- 文件权限问题
- 跨域限制

解决方案：
- 检查 `url_prefix` 配置是否正确
- 确保文件目录有正确的读取权限
- 配置正确的 CORS 设置

## 8. API 文档

完整的 API 文档可通过 Swagger 访问：

```
http://localhost:8086/swagger/index.html
```

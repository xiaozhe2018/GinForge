# 文件上传微服务使用指南

## 1. 简介

文件上传微服务是 GinForge 框架的一个独立服务，提供文件上传、下载、管理等功能。该服务支持多种存储方式，包括本地存储、阿里云 OSS、AWS S3 和 MinIO 等。

### 1.1 核心特性

- **多种存储支持**：本地存储、阿里云 OSS、AWS S3、MinIO
- **完整的文件管理**：上传、下载、删除、查询
- **元数据管理**：文件信息、统计数据
- **高级功能**：分片上传、防盗链、签名URL
- **安全控制**：文件类型验证、大小限制

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

### 3.2 分片上传大文件

对于大文件，可以使用分片上传功能：

#### 3.2.1 初始化上传

```bash
curl -X POST http://localhost:8086/api/v1/files/chunks/init \
  -F "file_name=large_video.mp4" \
  -F "total_chunks=10" \
  -F "total_size=104857600" \
  -F "chunk_size=10485760"
```

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": "550e8400-e29b-41d4-a716-446655440000" // 上传ID
}
```

#### 3.2.2 上传分片

```bash
curl -X POST http://localhost:8086/api/v1/files/chunks/upload \
  -F "upload_id=550e8400-e29b-41d4-a716-446655440000" \
  -F "chunk_index=0" \
  -F "total_chunks=10" \
  -F "file=@chunk_0.part" \
  -F "file_name=large_video.mp4"
```

#### 3.2.3 合并分片

```bash
curl -X POST http://localhost:8086/api/v1/files/chunks/merge \
  -H "Content-Type: application/json" \
  -d '{
    "upload_id": "550e8400-e29b-41d4-a716-446655440000",
    "file_name": "large_video.mp4"
  }'
```

### 3.3 下载文件

```bash
curl -X GET http://localhost:8086/api/v1/files/1/download \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o downloaded_file.jpg
```

### 3.4 获取文件列表

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

### 3.5 删除文件

```bash
curl -X DELETE http://localhost:8086/api/v1/files/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3.6 获取文件统计信息

```bash
curl -X GET http://localhost:8086/api/v1/files/statistics \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 4. 前端集成

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

### 4.3 分片上传组件

对于大文件，可以使用以下分片上传组件：

```javascript
class ChunkUploader {
  constructor(options) {
    this.options = Object.assign({
      url: 'http://localhost:8086/api/v1/files',
      chunkSize: 5 * 1024 * 1024, // 5MB
      threads: 3,
      retries: 3,
      onProgress: null,
      onSuccess: null,
      onError: null
    }, options);
    
    this.uploadId = null;
    this.chunks = [];
    this.uploadedChunks = new Set();
    this.uploading = false;
  }
  
  // 上传文件
  async upload(file, metadata = {}) {
    if (this.uploading) {
      throw new Error('Upload already in progress');
    }
    
    this.uploading = true;
    this.file = file;
    this.metadata = metadata;
    
    try {
      // 1. 初始化上传
      await this.initUpload();
      
      // 2. 准备分片
      this.prepareChunks();
      
      // 3. 上传分片
      await this.uploadChunks();
      
      // 4. 合并分片
      const result = await this.mergeChunks();
      
      if (this.options.onSuccess) {
        this.options.onSuccess(result);
      }
      
      return result;
    } catch (error) {
      if (this.options.onError) {
        this.options.onError(error);
      }
      throw error;
    } finally {
      this.uploading = false;
    }
  }
  
  // 初始化上传
  async initUpload() {
    const formData = new FormData();
    formData.append('file_name', this.file.name);
    formData.append('total_size', this.file.size);
    formData.append('chunk_size', this.options.chunkSize);
    
    // 计算分片数
    const totalChunks = Math.ceil(this.file.size / this.options.chunkSize);
    formData.append('total_chunks', totalChunks);
    
    // 添加元数据
    if (this.metadata.sub_path) {
      formData.append('sub_path', this.metadata.sub_path);
    }
    
    const response = await fetch(`${this.options.url}/chunks/init`, {
      method: 'POST',
      body: formData
    });
    
    const result = await response.json();
    this.uploadId = result.data;
    return this.uploadId;
  }
  
  // 其他方法...
}

// 使用示例
const uploader = new ChunkUploader({
  url: 'http://localhost:8086/api/v1/files',
  onProgress: (progress) => {
    document.getElementById('progress').value = progress;
  }
});

uploader.upload(file, { sub_path: 'videos' });
```

## 5. 配置详解

### 5.1 存储配置

文件服务支持多种存储方式，可以通过配置文件进行设置：

#### 本地存储

```yaml
storage:
  type: local
  local:
    base_path: ./uploads
    temp_path: ./temp_uploads
  url_prefix: http://localhost:8086/uploads
```

#### 阿里云 OSS

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

#### AWS S3

```yaml
storage:
  type: s3
  s3:
    endpoint: s3.amazonaws.com
    access_key_id: YOUR_ACCESS_KEY_ID
    access_key_secret: YOUR_ACCESS_KEY_SECRET
    bucket_name: your-bucket-name
    region: us-west-1
  url_prefix: https://your-bucket-name.s3.amazonaws.com
```

### 5.2 安全配置

#### 文件类型限制

```yaml
storage:
  # 允许的文件类型
  allowed_types:
    # 图片类型
    image:
      mimes: ["image/jpeg", "image/png", "image/gif"]
      exts: [".jpg", ".jpeg", ".png", ".gif"]
    # 文档类型
    document:
      mimes: ["application/pdf", "application/msword"]
      exts: [".pdf", ".doc", ".docx"]
  
  # 禁止的扩展名
  forbidden_exts: [".php", ".exe", ".sh", ".bat"]
```

#### 防盗链设置

```yaml
storage:
  security:
    # 防盗链
    anti_leech:
      enabled: true
      allowed_referers: ["example.com", "*.example.com"]
      allow_empty_referer: false
```

#### 签名URL

```yaml
storage:
  security:
    # 签名URL
    signed_url:
      enabled: true
      expire_time: 30m
```

### 5.3 分片上传配置

```yaml
storage:
  # 分片上传配置
  chunk_upload:
    enabled: true
    chunk_size: 5242880 # 5MB
    timeout: 24h
    chunk_dir: chunks
```

## 6. 常见问题

### 6.1 文件上传失败

可能原因：
- 文件大小超过限制
- 文件类型不被允许
- 存储路径权限问题

解决方案：
- 检查配置中的 `max_file_size` 设置
- 检查 `allowed_types` 和 `forbidden_exts` 设置
- 确保存储目录有正确的读写权限

### 6.2 无法访问上传的文件

可能原因：
- URL 前缀配置错误
- 文件权限问题
- 防盗链限制

解决方案：
- 检查 `url_prefix` 配置是否正确
- 确保文件目录有正确的读取权限
- 检查防盗链设置

### 6.3 分片上传问题

可能原因：
- 分片上传未启用
- 分片大小设置不合理
- 临时目录权限问题

解决方案：
- 确保 `chunk_upload.enabled` 设置为 `true`
- 调整分片大小（推荐5MB）
- 检查临时目录权限

## 7. API 文档

完整的 API 文档可通过 Swagger 访问：

```
http://localhost:8086/swagger/index.html
```

## 8. 最佳实践

1. **使用合适的存储方式**：生产环境建议使用对象存储
2. **设置合理的文件大小限制**：根据业务需求设置
3. **启用分片上传**：对于大文件，使用分片上传
4. **配置防盗链**：防止资源被盗用
5. **使用签名URL**：对敏感文件使用签名URL
6. **定期清理临时文件**：设置定时任务清理临时文件
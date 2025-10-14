# 分片上传使用指南

## 1. 什么是分片上传

分片上传是将大文件分成多个小块（分片）进行上传的技术，适用于以下场景：

- 上传大文件（通常大于10MB）
- 网络环境不稳定的情况
- 需要断点续传功能
- 需要显示上传进度

## 2. 使用场景

- 视频上传
- 大型文档上传
- 数据备份
- 批量图片上传

## 3. 分片上传流程

分片上传的基本流程如下：

1. **初始化上传**：获取上传ID
2. **上传分片**：将文件分成多个分片，逐个上传
3. **合并分片**：所有分片上传完成后，请求服务器合并分片
4. **完成上传**：获取最终的文件信息

## 4. API 使用方法

### 4.1 初始化上传

```http
POST /api/v1/files/chunks/init
Content-Type: multipart/form-data

file_name=large_file.mp4
total_chunks=10
total_size=104857600
chunk_size=10485760
sub_path=videos
```

参数说明：

- `file_name`: 文件名称
- `total_chunks`: 总分片数
- `total_size`: 文件总大小（字节）
- `chunk_size`: 每个分片的大小（字节）
- `sub_path`: 子路径（可选）

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": "550e8400-e29b-41d4-a716-446655440000" // 上传ID
}
```

### 4.2 上传分片

```http
POST /api/v1/files/chunks/upload
Content-Type: multipart/form-data

upload_id=550e8400-e29b-41d4-a716-446655440000
chunk_index=0
total_chunks=10
file=@chunk_0.part
file_name=large_file.mp4
```

参数说明：

- `upload_id`: 上传ID（初始化时获取）
- `chunk_index`: 分片索引（从0开始）
- `total_chunks`: 总分片数
- `file`: 分片文件
- `file_name`: 文件名称

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "upload_id": "550e8400-e29b-41d4-a716-446655440000",
    "chunk_index": 0,
    "total_chunks": 10,
    "completed": false
  }
}
```

### 4.3 合并分片

```http
POST /api/v1/files/chunks/merge
Content-Type: application/json

{
  "upload_id": "550e8400-e29b-41d4-a716-446655440000",
  "file_name": "large_file.mp4"
}
```

参数说明：

- `upload_id`: 上传ID
- `file_name`: 文件名称

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "original_name": "large_file.mp4",
    "file_name": "large_file_1642342342_a1b2c3d4.mp4",
    "file_size": 104857600,
    "file_type": "video",
    "mime_type": "video/mp4",
    "url": "http://localhost:8086/uploads/videos/large_file_1642342342_a1b2c3d4.mp4",
    "hash": "d41d8cd98f00b204e9800998ecf8427e",
    "upload_time": "2025-01-15T10:30:45Z"
  }
}
```

### 4.4 获取上传状态

```http
GET /api/v1/files/chunks/status/550e8400-e29b-41d4-a716-446655440000
```

响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "upload_id": "550e8400-e29b-41d4-a716-446655440000",
    "file_name": "large_file.mp4",
    "total_chunks": 10,
    "chunks": {
      "0": true,
      "1": true,
      "2": true
    }
  }
}
```

### 4.5 中止上传

```http
DELETE /api/v1/files/chunks/abort/550e8400-e29b-41d4-a716-446655440000
```

## 5. 前端实现

### 5.1 基本实现

```javascript
// 分片上传类
class ChunkUploader {
  constructor(options) {
    this.options = {
      url: 'http://localhost:8086/api/v1/files',
      chunkSize: 5 * 1024 * 1024, // 5MB
      threads: 3,
      ...options
    };
  }
  
  // 上传文件
  async upload(file) {
    // 1. 初始化上传
    const uploadId = await this.initUpload(file);
    
    // 2. 准备分片
    const chunks = this.prepareChunks(file);
    
    // 3. 上传分片
    await this.uploadChunks(uploadId, chunks);
    
    // 4. 合并分片
    return await this.mergeChunks(uploadId, file.name);
  }
  
  // 初始化上传
  async initUpload(file) {
    const formData = new FormData();
    formData.append('file_name', file.name);
    formData.append('total_size', file.size);
    formData.append('chunk_size', this.options.chunkSize);
    formData.append('total_chunks', Math.ceil(file.size / this.options.chunkSize));
    
    const response = await fetch(`${this.options.url}/chunks/init`, {
      method: 'POST',
      body: formData
    });
    
    const result = await response.json();
    return result.data;
  }
  
  // 准备分片
  prepareChunks(file) {
    const chunks = [];
    const totalChunks = Math.ceil(file.size / this.options.chunkSize);
    
    for (let i = 0; i < totalChunks; i++) {
      const start = i * this.options.chunkSize;
      const end = Math.min(file.size, start + this.options.chunkSize);
      chunks.push({
        index: i,
        blob: file.slice(start, end)
      });
    }
    
    return chunks;
  }
  
  // 上传分片
  async uploadChunks(uploadId, chunks) {
    const promises = [];
    const totalChunks = chunks.length;
    
    for (const chunk of chunks) {
      const promise = this.uploadChunk(uploadId, chunk, totalChunks);
      promises.push(promise);
      
      // 控制并发数
      if (promises.length >= this.options.threads) {
        await Promise.race(promises);
      }
    }
    
    return Promise.all(promises);
  }
  
  // 上传单个分片
  async uploadChunk(uploadId, chunk, totalChunks) {
    const formData = new FormData();
    formData.append('upload_id', uploadId);
    formData.append('chunk_index', chunk.index);
    formData.append('total_chunks', totalChunks);
    formData.append('file', chunk.blob);
    
    const response = await fetch(`${this.options.url}/chunks/upload`, {
      method: 'POST',
      body: formData
    });
    
    return response.json();
  }
  
  // 合并分片
  async mergeChunks(uploadId, fileName) {
    const response = await fetch(`${this.options.url}/chunks/merge`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        upload_id: uploadId,
        file_name: fileName
      })
    });
    
    const result = await response.json();
    return result.data;
  }
}

// 使用示例
const uploader = new ChunkUploader();
const fileInput = document.getElementById('file');
fileInput.addEventListener('change', async (e) => {
  const file = e.target.files[0];
  if (file) {
    try {
      const result = await uploader.upload(file);
      console.log('Upload complete:', result);
    } catch (error) {
      console.error('Upload failed:', error);
    }
  }
});
```

### 5.2 带进度和断点续传的实现

```javascript
class AdvancedChunkUploader {
  constructor(options) {
    this.options = {
      url: 'http://localhost:8086/api/v1/files',
      chunkSize: 5 * 1024 * 1024, // 5MB
      threads: 3,
      retries: 3,
      onProgress: null,
      ...options
    };
    
    this.uploadId = null;
    this.uploadedChunks = new Set();
  }
  
  // 上传文件
  async upload(file) {
    // 1. 初始化上传
    this.uploadId = await this.initUpload(file);
    
    // 2. 获取已上传的分片
    await this.getUploadedChunks();
    
    // 3. 准备分片
    const chunks = this.prepareChunks(file);
    
    // 4. 上传分片
    await this.uploadChunks(chunks);
    
    // 5. 合并分片
    return await this.mergeChunks(file.name);
  }
  
  // 获取已上传的分片
  async getUploadedChunks() {
    try {
      const response = await fetch(`${this.options.url}/chunks/list/${this.uploadId}`);
      const result = await response.json();
      this.uploadedChunks = new Set(result.data || []);
    } catch (error) {
      console.warn('Failed to get uploaded chunks:', error);
    }
  }
  
  // 上传分片（带进度和重试）
  async uploadChunks(chunks) {
    const totalChunks = chunks.length;
    let completedChunks = this.uploadedChunks.size;
    
    // 过滤出未上传的分片
    const pendingChunks = chunks.filter(chunk => !this.uploadedChunks.has(chunk.index));
    
    // 更新进度
    if (this.options.onProgress) {
      this.options.onProgress((completedChunks / totalChunks) * 100);
    }
    
    // 并发上传
    const queue = [...pendingChunks];
    const activeUploads = new Set();
    
    return new Promise((resolve, reject) => {
      const startUpload = async () => {
        if (queue.length === 0 && activeUploads.size === 0) {
          resolve();
          return;
        }
        
        while (queue.length > 0 && activeUploads.size < this.options.threads) {
          const chunk = queue.shift();
          activeUploads.add(chunk.index);
          
          this.uploadChunkWithRetry(chunk, totalChunks)
            .then(() => {
              activeUploads.delete(chunk.index);
              this.uploadedChunks.add(chunk.index);
              completedChunks++;
              
              // 更新进度
              if (this.options.onProgress) {
                this.options.onProgress((completedChunks / totalChunks) * 100);
              }
              
              startUpload();
            })
            .catch(error => {
              activeUploads.delete(chunk.index);
              reject(error);
            });
        }
      };
      
      startUpload();
    });
  }
  
  // 带重试的分片上传
  async uploadChunkWithRetry(chunk, totalChunks, retryCount = 0) {
    try {
      await this.uploadChunk(this.uploadId, chunk, totalChunks);
    } catch (error) {
      if (retryCount < this.options.retries) {
        // 延迟重试
        await new Promise(resolve => setTimeout(resolve, 1000 * Math.pow(2, retryCount)));
        return this.uploadChunkWithRetry(chunk, totalChunks, retryCount + 1);
      }
      throw error;
    }
  }
  
  // 保存上传状态（断点续传）
  saveUploadState(file) {
    const state = {
      uploadId: this.uploadId,
      fileName: file.name,
      fileSize: file.size,
      uploadedChunks: Array.from(this.uploadedChunks),
      timestamp: Date.now()
    };
    
    localStorage.setItem('chunkUploadState', JSON.stringify(state));
  }
  
  // 加载上传状态
  loadUploadState() {
    const stateJson = localStorage.getItem('chunkUploadState');
    if (!stateJson) return null;
    
    try {
      const state = JSON.parse(stateJson);
      this.uploadId = state.uploadId;
      this.uploadedChunks = new Set(state.uploadedChunks);
      return state;
    } catch (error) {
      return null;
    }
  }
}
```

## 6. 最佳实践

### 6.1 选择合适的分片大小

分片大小的选择需要考虑以下因素：

- **网络环境**：网络越不稳定，分片越小越好
- **服务器配置**：服务器处理能力强，可以选择较大的分片
- **文件大小**：文件越大，分片可以相应增大

推荐的分片大小：

- 较差网络环境：1-2MB
- 普通网络环境：5MB
- 良好网络环境：10MB

### 6.2 并发上传控制

并发上传可以提高上传速度，但也需要控制并发数量：

- 普通客户端：2-3个并发
- 高性能客户端：4-6个并发

### 6.3 错误处理和重试

实现健壮的错误处理和重试机制：

- 对每个分片设置独立的重试次数
- 实现指数退避算法（exponential backoff）
- 记录失败的分片，以便断点续传

### 6.4 断点续传

实现断点续传的关键步骤：

1. 在本地存储保存上传ID和进度
2. 页面刷新或重新打开时，检查是否有未完成的上传
3. 如果有，获取已上传的分片列表
4. 只上传未完成的分片

## 7. 配置

在 `configs/file-api.yaml` 中配置分片上传：

```yaml
storage:
  # 分片上传配置
  chunk_upload:
    # 是否启用分片上传
    enabled: true
    # 分片大小 (字节)
    chunk_size: 5242880 # 5MB
    # 分片上传会话超时时间
    timeout: 24h
    # 分片存储目录
    chunk_dir: chunks
```

## 8. 常见问题

### 8.1 上传失败

可能原因：
- 网络不稳定
- 分片大小设置不合理
- 服务器超时设置过短

解决方案：
- 实现重试机制
- 调整分片大小
- 增加服务器超时设置

### 8.2 合并失败

可能原因：
- 部分分片上传失败
- 服务器临时目录权限问题
- 磁盘空间不足

解决方案：
- 验证所有分片都已上传成功
- 检查临时目录权限
- 确保服务器有足够的磁盘空间

### 8.3 上传进度不准确

可能原因：
- 分片上传成功但未正确记录
- 并发上传导致进度计算错误

解决方案：
- 定期获取服务器上的分片状态
- 改进进度计算逻辑

## 9. 性能优化

1. **预加载分片**：提前加载下一个要上传的分片
2. **动态调整并发数**：根据网络状况动态调整并发上传数量
3. **压缩大文件**：在客户端对大文件进行压缩后再上传
4. **使用Web Workers**：在后台线程处理分片计算和上传

## 10. 安全考虑

1. **验证上传权限**：每个分片上传请求都应验证用户权限
2. **文件类型验证**：服务器应对每个分片进行基本验证
3. **限制上传大小**：设置合理的上传限制

## 11. 总结

分片上传是处理大文件上传的有效方案，GinForge框架提供了完整的分片上传支持。通过合理配置和使用本文档中的最佳实践，可以实现高效、可靠的大文件上传功能。
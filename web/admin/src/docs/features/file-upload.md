# 文件上传

GinForge 提供了完整的文件上传解决方案，支持本地存储和云存储。

## 🎯 功能特性

- ✅ **多种存储**：本地、阿里云 OSS、AWS S3、MinIO
- ✅ **分片上传**：支持大文件分片上传
- ✅ **文件管理**：上传、下载、删除、查询
- ✅ **安全控制**：文件类型和大小验证
- ✅ **元数据管理**：文件信息和统计

## 🔧 配置

### 本地存储配置

```yaml
# configs/file-api.yaml
storage:
  type: local
  local:
    base_path: ./uploads
  url_prefix: http://localhost:8086/uploads
  max_file_size: 104857600  # 100MB
```

### 云存储配置（阿里云 OSS）

```yaml
storage:
  type: oss
  oss:
    endpoint: oss-cn-hangzhou.aliyuncs.com
    access_key_id: ${OSS_ACCESS_KEY}
    access_key_secret: ${OSS_ACCESS_SECRET}
    bucket_name: my-bucket
```

## 📤 基础上传

### 后端实现

```go
// handler/file_handler.go
func (h *FileHandler) Upload(c *gin.Context) {
    // 获取上传的文件
    file, err := c.FormFile("file")
    if err != nil {
        response.BadRequest(c, "请选择文件")
        return
    }
    
    // 验证文件大小
    if file.Size > h.maxFileSize {
        response.BadRequest(c, "文件大小超过限制")
        return
    }
    
    // 保存文件
    result, err := h.uploadService.SaveFile(file, c.GetString("user_id"))
    if err != nil {
        response.InternalError(c, "文件上传失败")
        return
    }
    
    response.Success(c, result)
}
```

### 前端实现

```vue
<template>
  <el-upload
    action="/api/v1/files/upload"
    :headers="{ Authorization: `Bearer ${token}` }"
    :on-success="handleSuccess"
    :on-error="handleError"
    :before-upload="beforeUpload"
  >
    <el-button type="primary">点击上传</el-button>
  </el-upload>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'

const token = localStorage.getItem('admin_token')

// 上传前验证
const beforeUpload = (file: File) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

// 上传成功
const handleSuccess = (response: any) => {
  ElMessage.success('上传成功')
  console.log('文件 URL:', response.data.url)
}

// 上传失败
const handleError = (error: any) => {
  ElMessage.error('上传失败')
}
</script>
```

## 📦 分片上传（大文件）

### 使用场景

- 文件大小超过 100MB
- 网络不稳定，需要断点续传
- 需要显示上传进度

### 实现流程

```
1. 初始化上传
   POST /api/v1/files/chunks/init
   ↓
2. 获取 upload_id
   ↓
3. 循环上传每个分片
   POST /api/v1/files/chunks/upload
   ↓
4. 所有分片上传完成
   ↓
5. 合并分片
   POST /api/v1/files/chunks/merge
   ↓
6. 获取最终文件 URL
```

### 前端示例

```typescript
// 分片上传工具类
class ChunkUploader {
  private chunkSize = 10 * 1024 * 1024  // 10MB per chunk
  
  async upload(file: File) {
    // 1. 计算分片信息
    const totalChunks = Math.ceil(file.size / this.chunkSize)
    
    // 2. 初始化上传
    const uploadId = await this.initUpload(file.name, totalChunks, file.size)
    
    // 3. 上传每个分片
    for (let i = 0; i < totalChunks; i++) {
      const start = i * this.chunkSize
      const end = Math.min(start + this.chunkSize, file.size)
      const chunk = file.slice(start, end)
      
      await this.uploadChunk(uploadId, i, totalChunks, chunk, file.name)
      
      // 更新进度
      const progress = Math.round((i + 1) / totalChunks * 100)
      this.onProgress(progress)
    }
    
    // 4. 合并分片
    const result = await this.mergeChunks(uploadId, file.name)
    return result
  }
}
```

## 📥 文件下载

### 后端实现

```go
func (h *FileHandler) Download(c *gin.Context) {
    fileID := c.Param("id")
    
    // 查询文件信息
    fileInfo, err := h.fileService.GetFile(fileID)
    if err != nil {
        response.NotFound(c, "文件不存在")
        return
    }
    
    // 设置响应头
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.OriginalName))
    c.Header("Content-Type", fileInfo.MimeType)
    
    // 发送文件
    c.File(fileInfo.FilePath)
}
```

### 前端实现

```typescript
// 下载文件
const downloadFile = async (fileId: number, fileName: string) => {
  const response = await request.get(`/api/v1/files/${fileId}/download`, {
    responseType: 'blob'
  })
  
  // 创建下载链接
  const url = window.URL.createObjectURL(new Blob([response]))
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', fileName)
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(url)
}
```

## 🖼️ 图片上传（头像示例）

### 完整示例

```vue
<template>
  <el-upload
    class="avatar-uploader"
    action="/api/v1/files/upload"
    :headers="uploadHeaders"
    :show-file-list="false"
    :on-success="handleAvatarSuccess"
    :before-upload="beforeAvatarUpload"
  >
    <img v-if="imageUrl" :src="imageUrl" class="avatar">
    <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
  </el-upload>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const imageUrl = ref('')
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('admin_token')}`
}

const beforeAvatarUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

const handleAvatarSuccess = (response: any) => {
  imageUrl.value = response.data.url
  ElMessage.success('上传成功')
}
</script>

<style scoped>
.avatar-uploader .avatar {
  width: 178px;
  height: 178px;
  display: block;
  border-radius: 6px;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
}
</style>
```

## 📊 文件管理

### 获取文件列表

```go
// handler/file_handler.go
func (h *FileHandler) List(c *gin.Context) {
    var req model.FileListRequest
    c.ShouldBindQuery(&req)
    
    // 设置默认值
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    
    // 查询文件列表
    result, err := h.fileService.List(&req)
    if err != nil {
        response.InternalError(c, "获取文件列表失败")
        return
    }
    
    response.Success(c, result)
}
```

### 删除文件

```go
func (h *FileHandler) Delete(c *gin.Context) {
    fileID := c.Param("id")
    
    if err := h.fileService.DeleteFile(fileID); err != nil {
        response.InternalError(c, "删除文件失败")
        return
    }
    
    response.Success(c, gin.H{"message": "删除成功"})
}
```

## 🔐 安全控制

### 文件类型验证

```go
var allowedTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
    "application/pdf": true,
}

func validateFileType(mimeType string) bool {
    return allowedTypes[mimeType]
}
```

### 文件大小限制

```go
func (h *FileHandler) Upload(c *gin.Context) {
    file, _ := c.FormFile("file")
    
    // 检查文件大小
    maxSize := h.config.GetInt64("storage.max_file_size")
    if file.Size > maxSize {
        response.BadRequest(c, fmt.Sprintf("文件大小不能超过 %dMB", maxSize/1024/1024))
        return
    }
    
    // ...
}
```

### 文件命名策略

```go
// 生成安全的文件名：原文件名_时间戳_随机字符串.扩展名
func generateFileName(originalName string) string {
    ext := filepath.Ext(originalName)
    baseName := strings.TrimSuffix(originalName, ext)
    timestamp := time.Now().UnixNano()
    randomStr := utils.RandomString(8)
    
    return fmt.Sprintf("%s_%d_%s%s", baseName, timestamp, randomStr, ext)
}
```

## 📍 目录结构

```
uploads/
├── avatars/           # 用户头像
│   └── avatar_xxx.jpg
├── documents/         # 文档文件
│   └── doc_xxx.pdf
├── images/            # 普通图片
│   └── image_xxx.png
├── chunks/            # 分片临时目录
│   └── upload_id/
│       ├── chunk_0
│       ├── chunk_1
│       └── ...
└── temp/              # 临时文件
```

## 💡 最佳实践

### 1. 分类存储

根据文件类型存储到不同目录：

```go
func getUploadPath(fileType string) string {
    switch fileType {
    case "image":
        return "uploads/images"
    case "document":
        return "uploads/documents"
    case "avatar":
        return "uploads/avatars"
    default:
        return "uploads/others"
    }
}
```

### 2. 文件清理

定期清理过期的临时文件：

```go
// 清理 7 天前的临时文件
func CleanTempFiles() error {
    tempDir := "uploads/temp"
    cutoff := time.Now().AddDate(0, 0, -7)
    
    files, _ := ioutil.ReadDir(tempDir)
    for _, file := range files {
        if file.ModTime().Before(cutoff) {
            os.Remove(filepath.Join(tempDir, file.Name()))
        }
    }
    
    return nil
}
```

### 3. 访问控制

```go
// 验证用户是否有权限访问文件
func (h *FileHandler) Download(c *gin.Context) {
    fileID := c.Param("id")
    userID := c.GetString("user_id")
    
    // 检查权限
    if !h.fileService.CanAccess(fileID, userID) {
        response.Forbidden(c, "无权访问此文件")
        return
    }
    
    // 下载文件...
}
```

## 📚 完整示例

查看完整实现：

- **文件服务**: `services/file-api/`
- **分片上传**: `pkg/upload/chunked/`
- **存储工厂**: `pkg/storage/factory/`
- **详细文档**: `docs/FILE_SERVICE.md`
- **分片上传文档**: `docs/CHUNK_UPLOAD.md`

## 🎯 下一步

- [缓存系统](./cache) - 使用 Redis 缓存
- [WebSocket](./websocket) - 实时文件上传进度推送
- [存储配置](../../STORAGE_CONFIG.md) - 详细的存储配置说明

---

**提示**: 生产环境建议使用云存储（OSS/S3），避免服务器存储压力。


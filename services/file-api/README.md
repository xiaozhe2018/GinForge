# File-API 文件上传服务

独立的文件上传微服务，端口：8086

## 快速启动

```bash
# 编译
go build -o ../../bin/file-api ./cmd/server

# 运行
../../bin/file-api
```

## 主要功能

- 文件上传/下载
- 文件列表查询
- 文件统计分析
- 支持本地存储和云存储（OSS/S3/MinIO）

## API接口

| 接口 | 方法 | 路径 |
|------|------|------|
| 上传 | POST | `/api/v1/files/upload` |
| 下载 | GET | `/api/v1/files/:id/download` |
| 列表 | GET | `/api/v1/files` |
| 详情 | GET | `/api/v1/files/:id` |
| 删除 | DELETE | `/api/v1/files/:id` |
| 统计 | GET | `/api/v1/files/statistics` |

## 配置

配置文件：`configs/file-api.yaml`

默认存储路径：`./uploads/`

## 文档

- 详细文档：[../../docs/FILE_SERVICE.md](../../docs/FILE_SERVICE.md)
- API文档：`http://localhost:8086/swagger/index.html`

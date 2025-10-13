# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建所有服务
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./services/user-api/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o merchant-api ./services/merchant-api/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o admin-api ./services/admin-api/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./services/gateway/cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# 复制构建的二进制文件
COPY --from=builder /app/user-api .
COPY --from=builder /app/merchant-api .
COPY --from=builder /app/admin-api .
COPY --from=builder /app/gateway .

# 创建数据目录
RUN mkdir -p data

EXPOSE 8080

# 默认启动网关服务
CMD ["./gateway"]

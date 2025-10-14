# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git make

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建所有服务到 bin/ 目录
RUN mkdir -p bin && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/user-api ./services/user-api/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/merchant-api ./services/merchant-api/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/admin-api ./services/admin-api/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/gateway ./services/gateway/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/gateway-worker ./services/gateway-worker/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/file-api ./services/file-api/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/demo ./services/demo/cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/ginforge ./cmd/cli

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/ ./bin/

# 复制配置文件
COPY --from=builder /app/configs/ ./configs/

# 创建必要的目录
RUN mkdir -p data uploads logs

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080 8081 8082 8083 8084 8086

# 默认启动网关服务
CMD ["./bin/gateway"]

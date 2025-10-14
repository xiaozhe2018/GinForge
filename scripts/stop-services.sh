#!/bin/bash

# GinForge 微服务停止脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  停止所有 GinForge 微服务${NC}"
echo -e "${YELLOW}========================================${NC}"

# 端口列表
PORTS=(8080 8081 8082 8083 8084 8085 8086)

# 停止函数
stop_port() {
    local port=$1
    local pid=$(lsof -ti :$port 2>/dev/null)
    
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}正在停止端口 $port 的服务 [PID: $pid]...${NC}"
        kill -9 $pid 2>/dev/null
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ 端口 $port 服务已停止${NC}"
        else
            echo -e "${RED}✗ 停止端口 $port 服务失败${NC}"
        fi
    else
        echo -e "  端口 $port 无服务运行"
    fi
}

echo ""
# 停止所有端口
for port in "${PORTS[@]}"; do
    stop_port $port
done

# 停止所有go run进程
echo ""
echo -e "${YELLOW}清理所有 go run 进程...${NC}"
pkill -f "go run ./services/.*/cmd/server" 2>/dev/null
pkill -f "services/.*/cmd/server" 2>/dev/null

echo ""
echo -e "${GREEN}所有服务已停止！${NC}"
echo ""


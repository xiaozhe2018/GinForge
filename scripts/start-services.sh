#!/bin/bash

# GinForge 微服务启动脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  GinForge 微服务启动脚本${NC}"
echo -e "${GREEN}========================================${NC}"

# 检查bin目录是否存在
if [ ! -d "$BIN_DIR" ]; then
    echo -e "${YELLOW}bin目录不存在，正在编译服务...${NC}"
    cd "$PROJECT_ROOT"
    make build
fi

# 服务配置
declare -A SERVICES=(
    ["user-api"]="8081"
    ["merchant-api"]="8082"
    ["admin-api"]="8083"
    ["gateway-worker"]="8084"
    ["demo"]="8085"
    ["file-api"]="8086"
    ["gateway"]="8080"
)

# 启动函数
start_service() {
    local service=$1
    local port=$2
    local bin_file="$BIN_DIR/$service"
    
    # 检查二进制文件是否存在
    if [ ! -f "$bin_file" ]; then
        echo -e "${RED}✗ $service 二进制文件不存在: $bin_file${NC}"
        return 1
    fi
    
    # 检查端口是否被占用
    if lsof -i :$port >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠ $service 端口 $port 已被占用，跳过启动${NC}"
        return 1
    fi
    
    # 启动服务
    cd "$PROJECT_ROOT"
    nohup "$bin_file" > "logs/$service.log" 2>&1 &
    local pid=$!
    
    # 等待服务启动
    sleep 2
    
    # 检查服务是否启动成功
    if kill -0 $pid 2>/dev/null; then
        echo -e "${GREEN}✓ $service 启动成功 [PID: $pid, Port: $port]${NC}"
        return 0
    else
        echo -e "${RED}✗ $service 启动失败${NC}"
        return 1
    fi
}

# 创建必要的目录
mkdir -p "$PROJECT_ROOT/logs"
mkdir -p "$PROJECT_ROOT/uploads"
mkdir -p "$PROJECT_ROOT/data"

echo ""
echo -e "${GREEN}开始启动所有服务...${NC}"
echo ""

# 启动所有服务
for service in "${!SERVICES[@]}"; do
    port="${SERVICES[$service]}"
    start_service "$service" "$port"
done

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  服务启动完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${GREEN}可用服务：${NC}"
echo -e "  API网关:       http://localhost:8080"
echo -e "  用户端API:     http://localhost:8081"
echo -e "  商户端API:     http://localhost:8082"
echo -e "  管理后台API:   http://localhost:8083"
echo -e "  网关工作器:    http://localhost:8084"
echo -e "  演示服务:      http://localhost:8085"
echo -e "  ${YELLOW}文件服务:      http://localhost:8086${NC}"
echo ""
echo -e "${GREEN}API文档：${NC}"
echo -e "  管理后台:      http://localhost:8083/swagger/index.html"
echo -e "  ${YELLOW}文件服务:      http://localhost:8086/swagger/index.html${NC}"
echo ""
echo -e "${GREEN}查看日志：${NC}"
echo -e "  tail -f logs/file-api.log"
echo ""
echo -e "${GREEN}停止所有服务：${NC}"
echo -e "  make stop  或  ./scripts/stop-services.sh"
echo ""


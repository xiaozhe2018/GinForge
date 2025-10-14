#!/bin/bash

# GinForge 生产环境一键部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   🚀 GinForge 生产环境部署脚本                            ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEPLOY_DIR="$PROJECT_ROOT/deployments"

# 检查必要工具
check_requirements() {
    echo -e "${BLUE}📋 检查环境要求...${NC}"
    
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker: $(docker --version)${NC}"
    
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}❌ Docker Compose 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker Compose: $(docker-compose --version)${NC}"
    
    if ! command -v node &> /dev/null; then
        echo -e "${YELLOW}⚠️  Node.js 未安装，前端需要预先构建${NC}"
    else
        echo -e "${GREEN}✅ Node.js: $(node --version)${NC}"
    fi
    echo ""
}

# 检查配置文件
check_config() {
    echo -e "${BLUE}🔍 检查配置文件...${NC}"
    
    if [ ! -f "$DEPLOY_DIR/.env.production" ]; then
        echo -e "${YELLOW}⚠️  .env.production 不存在，正在创建...${NC}"
        cp "$DEPLOY_DIR/env.production.example" "$DEPLOY_DIR/.env.production"
        echo -e "${RED}❗ 请先编辑 deployments/.env.production 配置文件！${NC}"
        echo -e "${RED}   必须修改以下配置：${NC}"
        echo -e "${RED}   - MYSQL_PASSWORD${NC}"
        echo -e "${RED}   - REDIS_PASSWORD${NC}"
        echo -e "${RED}   - JWT_SECRET${NC}"
        echo -e "${RED}   - CORS_ORIGINS${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ 环境配置文件存在${NC}"
    
    if [ ! -f "$PROJECT_ROOT/web/admin/dist/index.html" ]; then
        echo -e "${YELLOW}⚠️  前端未构建${NC}"
        BUILD_FRONTEND=true
    else
        echo -e "${GREEN}✅ 前端已构建${NC}"
        BUILD_FRONTEND=false
    fi
    echo ""
}

# 构建前端
build_frontend() {
    if [ "$BUILD_FRONTEND" = true ]; then
        echo -e "${BLUE}🔨 构建前端项目...${NC}"
        cd "$PROJECT_ROOT/web/admin"
        
        if [ ! -d "node_modules" ]; then
            echo -e "${YELLOW}安装前端依赖...${NC}"
            npm install
        fi
        
        echo -e "${YELLOW}构建生产版本...${NC}"
        npm run build
        echo -e "${GREEN}✅ 前端构建完成${NC}"
        echo ""
    fi
}

# 启动服务
deploy_services() {
    echo -e "${BLUE}🚀 启动 Docker 服务...${NC}"
    cd "$DEPLOY_DIR"
    
    # 拉取镜像
    echo -e "${YELLOW}拉取基础镜像...${NC}"
    docker-compose -f docker-compose.prod.yml --env-file .env.production pull
    
    # 构建服务
    echo -e "${YELLOW}构建服务镜像...${NC}"
    docker-compose -f docker-compose.prod.yml --env-file .env.production build --no-cache
    
    # 启动服务
    echo -e "${YELLOW}启动所有服务...${NC}"
    docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
    
    echo -e "${GREEN}✅ 服务启动命令已执行${NC}"
    echo ""
}

# 等待服务就绪
wait_for_services() {
    echo -e "${BLUE}⏳ 等待服务启动...${NC}"
    sleep 15
    
    echo -e "${YELLOW}检查服务健康状态...${NC}"
    docker-compose -f "$DEPLOY_DIR/docker-compose.prod.yml" ps
    echo ""
}

# 验证部署
verify_deployment() {
    echo -e "${BLUE}✅ 验证部署...${NC}"
    
    # 检查健康端点
    if curl -f -s http://localhost/healthz > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 健康检查通过${NC}"
    else
        echo -e "${RED}❌ 健康检查失败${NC}"
    fi
    
    # 检查前端
    if curl -f -s http://localhost > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 前端页面可访问${NC}"
    else
        echo -e "${YELLOW}⚠️  前端页面访问失败${NC}"
    fi
    echo ""
}

# 显示访问信息
show_info() {
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   🎉 部署完成！                                           ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}🌐 访问地址：${NC}"
    echo -e "   前端管理后台: ${GREEN}http://localhost${NC}"
    echo -e "   API 网关:    ${GREEN}http://localhost:8080${NC}"
    echo -e "   健康检查:    ${GREEN}http://localhost/healthz${NC}"
    echo ""
    echo -e "${BLUE}🔑 默认账号：${NC}"
    echo -e "   用户名: ${GREEN}admin${NC}"
    echo -e "   密码:   ${GREEN}admin123${NC}"
    echo ""
    echo -e "${BLUE}📊 管理命令：${NC}"
    echo -e "   查看状态: ${YELLOW}docker-compose -f deployments/docker-compose.prod.yml ps${NC}"
    echo -e "   查看日志: ${YELLOW}docker-compose -f deployments/docker-compose.prod.yml logs -f${NC}"
    echo -e "   停止服务: ${YELLOW}docker-compose -f deployments/docker-compose.prod.yml down${NC}"
    echo ""
    echo -e "${GREEN}✨ 部署成功！${NC}"
}

# 主流程
main() {
    check_requirements
    check_config
    build_frontend
    deploy_services
    wait_for_services
    verify_deployment
    show_info
}

# 执行部署
main


#!/bin/bash

# GinForge 项目初始化脚本
# 用于首次设置项目环境

set -e

echo "🚀 GinForge 项目初始化"
echo "===================="

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 打印成功信息
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# 打印警告信息
print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# 打印错误信息
print_error() {
    echo -e "${RED}✗${NC} $1"
}

echo ""
echo "1️⃣ 检查环境..."

# 检查 Go
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go 已安装: $GO_VERSION"
else
    print_error "Go 未安装，请先安装 Go 1.20+"
    exit 1
fi

# 检查 Node.js
if command_exists node; then
    NODE_VERSION=$(node --version)
    print_success "Node.js 已安装: $NODE_VERSION"
else
    print_warning "Node.js 未安装，前端功能将不可用"
fi

# 检查 Docker (可选)
if command_exists docker; then
    print_success "Docker 已安装"
else
    print_warning "Docker 未安装，容器化功能将不可用"
fi

echo ""
echo "2️⃣ 安装 Go 依赖..."
go mod download
go mod tidy
print_success "Go 依赖安装完成"

echo ""
echo "3️⃣ 创建配置文件..."
if [ ! -f ".env" ]; then
    cp env.example .env
    print_success "创建 .env 文件"
    print_warning "请修改 .env 文件中的配置"
else
    print_warning ".env 文件已存在，跳过"
fi

echo ""
echo "4️⃣ 初始化数据库..."
if [ ! -f "goweb.db" ]; then
    print_success "将在首次运行时自动创建 SQLite 数据库"
else
    print_warning "数据库文件已存在"
fi

echo ""
echo "5️⃣ 安装前端依赖..."
if command_exists npm; then
    cd web/admin
    if [ ! -d "node_modules" ]; then
        npm install
        print_success "前端依赖安装完成"
    else
        print_warning "node_modules 已存在，跳过安装"
    fi
    cd ../..
else
    print_warning "跳过前端依赖安装（npm 不可用）"
fi

echo ""
echo "6️⃣ 生成 Swagger 文档..."
if command_exists swag; then
    make swagger 2>/dev/null || true
    print_success "Swagger 文档生成完成"
else
    print_warning "swag 未安装，跳过文档生成"
    echo "  安装命令: go install github.com/swaggo/swag/cmd/swag@latest"
fi

echo ""
echo "7️⃣ 创建必要的目录..."
mkdir -p bin
mkdir -p logs
mkdir -p uploads
print_success "目录创建完成"

echo ""
echo "===================="
echo "✨ 初始化完成！"
echo ""
echo "📚 下一步："
echo "  1. 修改 .env 文件中的配置"
echo "  2. 启动后端服务: go run ./services/admin-api/cmd/server"
echo "  3. 启动前端服务: cd web/admin && npm run dev"
echo "  4. 访问管理后台: http://localhost:3000"
echo "  5. 默认账号: admin / admin123"
echo ""
echo "📖 查看完整文档: docs/INDEX.md"
echo "🆘 遇到问题: docs/TROUBLESHOOTING.md"
echo ""


#!/bin/bash

# 初始化設定腳本

echo "=== 微服務專案初始化 ==="
echo ""

# 檢查 Docker
echo "檢查 Docker..."
if ! command -v docker &> /dev/null; then
    echo "✗ Docker 未安裝，請先安裝 Docker"
    exit 1
fi
echo "✓ Docker 已安裝"
echo ""

# 檢查 Docker Compose
echo "檢查 Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    echo "✗ Docker Compose 未安裝，請先安裝 Docker Compose"
    exit 1
fi
echo "✓ Docker Compose 已安裝"
echo ""

# 創建環境變數文件
echo "設定環境變數..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✓ .env 文件已創建"
else
    echo "⚠ .env 文件已存在，跳過"
fi
echo ""

# 下載 Go 依賴
echo "下載 Go 依賴..."
if command -v go &> /dev/null; then
    cd api-gateway && go mod download && cd ..
    cd services/user-service && go mod download && cd ../..
    echo "✓ Go 依賴已下載"
else
    echo "⚠ Go 未安裝，將在 Docker 容器中下載依賴"
fi
echo ""

# 安裝前端依賴
echo "安裝前端依賴..."
if command -v npm &> /dev/null; then
    cd frontend && npm install && cd ..
    echo "✓ 前端依賴已安裝"
else
    echo "⚠ npm 未安裝，將在 Docker 容器中安裝依賴"
fi
echo ""

# 建構 Docker 映像
echo "建構 Docker 映像..."
docker-compose build
echo "✓ Docker 映像建構完成"
echo ""

# 啟動服務
echo "啟動服務..."
docker-compose up -d
echo "✓ 服務已啟動"
echo ""

# 等待服務啟動
echo "等待服務啟動（30秒）..."
sleep 30
echo ""

# 健康檢查
echo "執行健康檢查..."
bash scripts/health-check.sh
echo ""

echo "=== 初始化完成 ==="
echo ""
echo "服務訪問地址："
echo "  Frontend:     http://localhost:3000"
echo "  API Gateway:  http://localhost:8080"
echo "  User Service: http://localhost:8081"
echo ""
echo "管理命令："
echo "  查看日誌:     docker-compose logs -f"
echo "  停止服務:     docker-compose down"
echo "  重啟服務:     docker-compose restart"
echo ""

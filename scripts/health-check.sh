#!/bin/bash

# 健康檢查腳本

echo "=== 微服務健康檢查 ==="
echo ""

# API Gateway
echo "檢查 API Gateway..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✓ API Gateway: 正常運行"
else
    echo "✗ API Gateway: 無法訪問"
fi
echo ""

# User Service
echo "檢查 User Service..."
if curl -s http://localhost:8081/health > /dev/null; then
    echo "✓ User Service: 正常運行"
else
    echo "✗ User Service: 無法訪問"
fi
echo ""

# Frontend
echo "檢查 Frontend..."
if curl -s http://localhost:3000 > /dev/null; then
    echo "✓ Frontend: 正常運行"
else
    echo "✗ Frontend: 無法訪問"
fi
echo ""

# PostgreSQL
echo "檢查 PostgreSQL..."
if docker-compose exec -T postgres pg_isready -U admin > /dev/null 2>&1; then
    echo "✓ PostgreSQL: 正常運行"
else
    echo "✗ PostgreSQL: 無法訪問"
fi
echo ""

# Redis
echo "檢查 Redis..."
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✓ Redis: 正常運行"
else
    echo "✗ Redis: 無法訪問"
fi
echo ""

# MongoDB
echo "檢查 MongoDB..."
if docker-compose exec -T mongodb mongosh --quiet --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo "✓ MongoDB: 正常運行"
else
    echo "✗ MongoDB: 無法訪問"
fi
echo ""

echo "=== 檢查完成 ==="

.PHONY: help build up down restart logs clean test

help: ## 顯示幫助資訊
	@echo "可用命令："
	@echo "  make build     - 建構所有服務"
	@echo "  make up        - 啟動所有服務"
	@echo "  make down      - 停止所有服務"
	@echo "  make restart   - 重啟所有服務"
	@echo "  make logs      - 查看所有服務日誌"
	@echo "  make clean     - 清理所有容器和資料"
	@echo "  make test      - 執行測試"
	@echo "  make dev-api   - 啟動 API Gateway (本地開發)"
	@echo "  make dev-user  - 啟動 User Service (本地開發)"
	@echo "  make dev-front - 啟動 Frontend (本地開發)"

build: ## 建構所有服務
	docker-compose build

up: ## 啟動所有服務
	docker-compose up -d

down: ## 停止所有服務
	docker-compose down

restart: down up ## 重啟所有服務

logs: ## 查看所有服務日誌
	docker-compose logs -f

logs-api: ## 查看 API Gateway 日誌
	docker-compose logs -f api-gateway

logs-user: ## 查看 User Service 日誌
	docker-compose logs -f user-service

logs-front: ## 查看 Frontend 日誌
	docker-compose logs -f frontend

clean: ## 清理所有容器、映像和資料
	docker-compose down -v
	docker system prune -f

ps: ## 查看服務狀態
	docker-compose ps

# 本地開發命令
dev-api: ## 本地運行 API Gateway
	cd api-gateway && go mod download && go run main.go

dev-user: ## 本地運行 User Service
	cd services/user-service && go mod download && go run main.go

dev-front: ## 本地運行 Frontend
	cd frontend && npm install && npm start

# 測試命令
test: ## 執行所有測試
	@echo "執行測試..."
	cd api-gateway && go test -v ./...
	cd services/user-service && go test -v ./...

# 數據庫管理
db-psql: ## 連接到 PostgreSQL
	docker-compose exec postgres psql -U admin -d userdb

db-mongo: ## 連接到 MongoDB
	docker-compose exec mongodb mongosh -u admin -p admin123

db-redis: ## 連接到 Redis
	docker-compose exec redis redis-cli

# 初始化
init: ## 初始化專案
	@echo "初始化專案..."
	cp .env.example .env
	@echo "環境變數文件已創建，請根據需要修改 .env"
	cd api-gateway && go mod download
	cd services/user-service && go mod download
	cd frontend && npm install
	@echo "依賴安裝完成！"

# 完整部署
deploy: clean build up ## 完整部署（清理 + 建構 + 啟動）
	@echo "部署完成！"
	@echo "Frontend: http://localhost:3000"
	@echo "API Gateway: http://localhost:8080"
	@echo "User Service: http://localhost:8081"

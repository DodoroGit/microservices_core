# Microservices Personal Website

本專案為微服務架構的個人網站，使用 React + Golang + PostgreSQL 技術棧。

## 架構概覽

```
┌─────────────┐
│   Frontend  │ (React - Port 3000)
│   (Nginx)   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ API Gateway │ (Golang - Port 8080)
└──────┬──────┘
       │
       ├─────────────────┐
       │                 │
       ▼                 ▼
┌──────────────┐   ┌──────────────┐
│ User Service │   │ Future       │
│ (Port 8081)  │   │ Services...  │
└──────┬───────┘   └──────────────┘
       │
       ├──────┬──────┬──────┐
       ▼      ▼      ▼      ▼
   ┌────┐  ┌────┐ ┌────┐  ┌────┐
   │ PG │  │Mongo││Redis│  │...│
   └────┘  └────┘ └────┘  └────┘
```

## 技術棧

### Frontend
- **React 18** - 前端框架
- **React Router** - 路由管理
- **Axios** - HTTP 客戶端
- **Nginx** - 生產環境 Web 伺服器

### Backend
- **Golang** - API Gateway 與微服務
- **Gin** - Web 框架
- **PostgreSQL** - 關聯式數據庫
- **MongoDB** - NoSQL 數據庫（預留）
- **Redis** - 快取與會話管理

### DevOps
- **Docker** - 容器化
- **Docker Compose** - 容器編排

## 專案結構

```
microservices_core/
├── frontend/                 # React 前端應用
│   ├── src/
│   │   ├── components/      # React 組件
│   │   ├── api.js          # API 客戶端
│   │   ├── App.js
│   │   └── index.js
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
│
├── api-gateway/             # API Gateway
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── services/
│   └── user-service/        # 用戶服務
│       ├── main.go
│       ├── Dockerfile
│       ├── go.mod
│       └── go.sum
│
├── docker-compose.yml       # Docker Compose 配置
├── .env.example            # 環境變數範例
└── README.md
```

## 快速開始

### 前置需求

- Docker & Docker Compose
- Git

### 1. 克隆專案

```bash
git clone <your-repo-url>
cd microservices_core
```

### 2. 環境變數設定

```bash
cp .env.example .env
# 根據需要修改 .env 文件
```

### 3. 啟動所有服務

```bash
# 建構並啟動所有容器
docker-compose up --build

# 或在背景執行
docker-compose up -d --build
```

### 4. 訪問應用

- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **User Service**: http://localhost:8081

### 5. 停止服務

```bash
docker-compose down

# 清除所有資料（包括數據庫）
docker-compose down -v
```

## 本地開發

### Frontend 開發

```bash
cd frontend
npm install
npm start
# 訪問 http://localhost:3000
```

### API Gateway 開發

```bash
cd api-gateway
go mod download
go run main.go
# 訪問 http://localhost:8080
```

### User Service 開發

```bash
cd services/user-service
go mod download
go run main.go
# 訪問 http://localhost:8081
```

## API 文檔

### User Service API

#### 註冊用戶
```bash
POST /users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

#### 用戶登入
```bash
POST /users/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### 獲取所有用戶
```bash
GET /users
```

#### 獲取單個用戶
```bash
GET /users/:id
```

#### 更新用戶
```bash
PUT /users/:id
Content-Type: application/json

{
  "username": "new_username"
}
```

#### 刪除用戶
```bash
DELETE /users/:id
```

## 數據庫

### PostgreSQL
- **Port**: 5432
- **Database**: userdb
- **User**: admin
- **Password**: admin123

### MongoDB
- **Port**: 27017
- **User**: admin
- **Password**: admin123

### Redis
- **Port**: 6379

## 新增微服務

若要新增新的微服務，請遵循以下步驟：

1. 在 `services/` 目錄下創建新服務資料夾
2. 實作服務邏輯和 Dockerfile
3. 在 `docker-compose.yml` 中添加服務配置
4. 在 API Gateway 中添加路由規則

範例：
```bash
mkdir services/new-service
cd services/new-service
# 創建 main.go, Dockerfile, go.mod 等
```

## 監控與日誌

查看服務日誌：
```bash
# 查看所有服務日誌
docker-compose logs

# 查看特定服務日誌
docker-compose logs user-service
docker-compose logs api-gateway
docker-compose logs frontend

# 持續追蹤日誌
docker-compose logs -f
```

## 故障排除

### 容器無法啟動
```bash
# 檢查容器狀態
docker-compose ps

# 重新建構容器
docker-compose up --build --force-recreate
```

### 數據庫連接失敗
```bash
# 檢查數據庫健康狀態
docker-compose ps postgres

# 進入 PostgreSQL 容器
docker-compose exec postgres psql -U admin -d userdb
```

### 清除所有資料重新開始
```bash
docker-compose down -v
docker system prune -a
docker-compose up --build
```

## 開發計劃

- [x] 基礎架構搭建
- [x] User Service 實作
- [x] API Gateway 實作
- [x] Frontend 基礎頁面
- [ ] JWT 認證機制
- [ ] 新增其他微服務
- [ ] 日誌聚合系統
- [ ] 監控與告警
- [ ] CI/CD Pipeline

## 授權

MIT License

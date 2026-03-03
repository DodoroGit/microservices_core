# Microservices Core

這是一個學習用的個人專案，目的是透過實作來深入理解**微服務架構**與各種後端設計模式。

比起功能本身，更在意的是把這裡當作一個**技術實驗站**——把每個學過的概念都落地成可以跑起來的程式。

## 學習目標

- 理解微服務架構的設計思路與各層職責
- 實踐 Dependency Injection（依賴注入），透過 interface 抽象達到模組解耦
- 讓 Unit Test 能夠優雅地進入專案，並區分 Unit / Integration Test 的職責
- 探索服務間通訊協定（RESTful → gRPC）的選擇與取捨

## 架構

```
Frontend (React + Vite)
    ↓ RESTful
API Gateway (Go / Gin)
    ↓ RESTful（目前）→ 規劃改為 gRPC
Backend Services
    ├── user-service  (Go)
    └── ai-service    (Python / FastAPI)
         ↓
      Ollama (llama3.2:3b)
Databases
    ├── PostgreSQL  (user-service)
    ├── MongoDB     (預留擴展)
    └── Redis       (ai-service chat history)
```

| 組件 | 技術 | 說明 |
|---|---|---|
| Frontend | React + Vite | Login / Register / Dashboard / AI Lab / Notes（進行中）|
| API Gateway | Go + Gin | 對外統一入口，負責路由轉發與 JWT 驗證 |
| user-service | Go + Gin + PostgreSQL | 使用者管理，含 Admin / User 角色權限 |
| ai-service | Python + FastAPI + Redis | 串接 Ollama 本地 LLM，提供 Chat API |

## 設計模式實踐（user-service）

user-service 是目前架構實驗的主要場域：

- **Layered Architecture**：Handler → Service → Repository，各層職責分明
- **Interface 抽象**：每層依賴 interface 而非具體實作，達到 DI 解耦
- **Unit Test**：Handler 層以 Mock（手寫 Mock）隔離 Service 依賴，純邏輯驗證
- **Integration Test**：Repository 層連接真實 PostgreSQL，以 build tag（`//go:build integration`）區隔，CI 中獨立運行

```
Handler     → 依賴 UserServiceInterface
Service     → 依賴 UserRepositoryInterface
Repository  → 操作 PostgreSQL
```

## 通訊協定

| 路徑 | 現況 | 規劃 |
|---|---|---|
| Frontend → API Gateway | RESTful | RESTful |
| API Gateway → Backend Services | RESTful | gRPC |
| Service → Service | - | gRPC |

> ai-service 使用 Ollama 本地推論，原因是本地端運算資源有限，以 llama3.2:3b 輕量模型為主。

## CI/CD

GitHub Actions 自動化流程：

- **Unit Tests**：push / PR 時自動觸發，不依賴任何外部資源
- **Integration Tests**：Unit Tests 通過後才執行，起動 PostgreSQL service container 進行完整驗證

## 目前進度

**已完成**
- [x] user-service 完整架構（Handler / Service / Repository + Interface + DI）
- [x] user-service Unit Test（Mock-based Handler Test）
- [x] user-service Integration Test（真實 DB，build tag 隔離）
- [x] API Gateway（HTTP reverse proxy + JWT 驗證 + CORS）
- [x] ai-service（FastAPI + Ollama + Redis chat history）
- [x] 前端介面（Login / Register / Dashboard / AI Lab）
- [x] Docker Compose 整合（PostgreSQL、MongoDB、Redis、Ollama、healthcheck）
- [x] GitHub Actions CI Pipeline（Unit + Integration Test 分層執行）

**進行中 / 規劃中**
- [ ] Notes 功能（前端頁面已建，後端 service 尚未實作）
- [ ] 服務間通訊改為 gRPC
- [ ] ai-service 設計模式完善（對齊 user-service 架構）

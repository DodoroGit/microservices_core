# Microservices Core

這是一個學習用的個人專案，目的是透過實作深入理解**微服務架構**與後端設計模式。

比起功能本身，更在意的是把這裡當作一個**技術實驗站**——把每個學過的概念都落地成可以跑起來的程式。

---

## 架構總覽

```
Frontend (React + Vite + Nginx)
    ↓ RESTful
API Gateway (Go / Gin)
    ├── JWT 驗證（含 Redis 黑名單）/ CORS / 反向代理
    └── RBAC（Admin / User）
    ↓ RESTful
Backend Services
    ├── user-service   (Go + Gin)          → PostgreSQL + Redis
    ├── note-service   (Go + Gin)          → MongoDB
    └── ai-service     (Python + FastAPI)  → note-service + Claude API
```

| 組件 | 技術 | 說明 |
|---|---|---|
| Frontend | React + Vite + Nginx | Login / Register / Dashboard / AI Lab / Notes |
| API Gateway | Go + Gin | 統一入口，JWT 驗證（含黑名單）、CORS、RBAC 權限控制 |
| user-service | Go + Gin + PostgreSQL + Redis | 使用者管理，Admin / User 角色，JWT 登出黑名單 |
| note-service | Go + Gin + MongoDB | 筆記 CRUD，依分類（專案 / 技術筆記 / 日常隨筆 / 履歷）篩選 |
| ai-service | Python + FastAPI + Anthropic SDK | 串接 Claude API，從筆記生成個人背景、專案介紹、技術能力等履歷素材 |

---

## 設計模式實踐

### Layered Architecture + Dependency Injection

三個後端服務皆實踐相同的分層架構，各層依賴抽象而非具體實作。

**Go（user-service / note-service）**

```
Handler  →  {User/Note}ServiceInterface（interface）
Service  →  {User/Note}RepositoryInterface（interface）
Repository  →  PostgreSQL / MongoDB
```

- Go 的 `interface` 為結構化型別，不需要明確宣告實作
- `main.go` 負責組裝整條依賴鏈（Repository → Service → Handler → Router）

**Python（ai-service）**

```
Router  →  AIServiceProtocol（Protocol）
Service →  NoteClientProtocol / ClaudeClientProtocol（Protocol）
Clients →  note-service（HTTP）+ Claude API
```

- Python 以 `typing.Protocol` 作為 interface 的等價物，同樣為結構化型別
- FastAPI 的 `Depends()` 負責組裝依賴，factory function 定義在 `ai_router.py` 頂部
- ai-service 無資料庫，`clients/` 層封裝所有外部 I/O（HTTP client + AI SDK）
- Protocol 定義在各自層的檔案頂部（實作 class 的正上方）

---

## 測試策略

Go 服務區分 Unit Test 與 Integration Test，ai-service 測試待補：

| | Unit Test | Integration Test |
|---|---|---|
| **Go** | Mock interface（testify/mock） | 真實 DB，build tag `//go:build integration` |
| **Python** | 待補 | 待補 |

**各層測試對象（Go）：**

```
Repository 層  →  mock DB driver
Service 層     →  Fake Repository（in-memory）
Handler 層     →  Fake Service + httptest
```

```bash
# Python
pytest -v                  # 只跑 Unit Test
pytest -v -m integration   # 只跑 Integration Test（需啟動對應 DB）

# Go
go test ./...                          # Unit Test
go test ./... -tags integration        # Integration Test
```

---

## CI/CD

GitHub Actions 三條 Pipeline，各自對應一個服務：

| Pipeline | Lint | Unit Test | Integration Test |
|---|---|---|---|
| user-service | golangci-lint | ✓ | PostgreSQL container |
| note-service | golangci-lint | ✓ | MongoDB container |
| ai-service | ruff | ✓ | Redis container |

執行順序：`lint → unit-test → integration-test`（前一階段失敗則停止）

note-service / ai-service 的 Pipeline 設有 `paths` 過濾，只有對應目錄異動時才觸發。

---

## 通訊協定

| 路徑 | 現況 | 規劃 |
|---|---|---|
| Frontend → API Gateway | RESTful | RESTful |
| API Gateway → Backend Services | RESTful | gRPC |
| Service → Service | RESTful | gRPC |

---

## 技術一覽

| 類別 | 技術 |
|---|---|
| 後端 | Go、Python、Gin、FastAPI |
| 前端 | React、Vite、Nginx |
| 資料庫 | PostgreSQL、MongoDB、Redis |
| 基礎設施 | Docker、Docker Compose |
| AI | Claude API（Anthropic）|
| 認證 | JWT、RBAC |
| 測試 | Unit Test、Integration Test |
| Lint | golangci-lint（Go）、ruff（Python）|
| CI/CD | GitHub Actions |

---

## 目前進度

**已完成**
- [x] user-service — Layered Architecture + Interface + DI + Unit / Integration Test
- [x] note-service — Layered Architecture + Interface + DI + Unit / Integration Test
- [x] ai-service — Layered Architecture + Protocol + DI
- [x] API Gateway — HTTP reverse proxy + JWT 驗證（含 Redis 黑名單）+ CORS + RBAC
- [x] JWT 登出機制 — 登出時將 token 寫入 Redis 黑名單，過期自動清除
- [x] 前端自動登出 — token 過期或被撤銷時，攔截 401 自動跳轉登入頁
- [x] 前端介面 — Login / Register / Dashboard / AI Lab / Notes
- [x] Docker Compose — 含 healthcheck，一鍵啟動全服務
- [x] GitHub Actions CI — Lint + Unit + Integration，三服務各自獨立 Pipeline

**規劃中**
- [ ] ai-service Unit / Integration Test
- [ ] 服務間通訊改為 gRPC

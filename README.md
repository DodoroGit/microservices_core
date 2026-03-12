# Microservices Core

這是一個學習用的個人專案，目的是透過實作深入理解**微服務架構**與後端設計模式。

比起功能本身，更在意的是把這裡當作一個**技術實驗站**——把每個學過的概念都落地成可以跑起來的程式。

---

## 架構總覽

```
Frontend (React + Vite + Nginx)
    ↓ RESTful
API Gateway (Go / Gin)
    ├── JWT 驗證 / CORS / 反向代理
    └── RBAC（Admin / User）
    ↓ RESTful
Backend Services
    ├── user-service   (Go + Gin)          → PostgreSQL
    ├── note-service   (Go + Gin)          → MongoDB
    └── ai-service     (Python + FastAPI)  → note-service + Claude API
```

| 組件 | 技術 | 說明 |
|---|---|---|
| Frontend | React + Vite + Nginx | Login / Register / Dashboard / AI Lab / Notes |
| API Gateway | Go + Gin | 統一入口，JWT 驗證、CORS、RBAC 權限控制 |
| user-service | Go + Gin + PostgreSQL | 使用者管理，Admin / User 角色 |
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
Router  →  ServiceProtocol（Protocol）
Service →  RepositoryProtocol（Protocol）
Repository  →  note-service（HTTP）+ Claude API
```

- Python 以 `typing.Protocol` 作為 interface 的等價物，同樣為結構化型別
- FastAPI 的 `Depends()` 取代 Go 的手動組裝，由框架在每個 request 進來時自動注入
- Protocol 定義在各自層的檔案頂部（實作 class 的正上方）

---

## 測試策略

三個服務皆區分 Unit Test 與 Integration Test：

| | Unit Test | Integration Test |
|---|---|---|
| **Go** | Mock interface（testify/mock） | 真實 DB，build tag `//go:build integration` |
| **Python** | Fake class / AsyncMock | 真實 DB，`@pytest.mark.integration` |

**各層測試對象：**

```
Repository 層  →  mock DB driver（AsyncMock / MagicMock）
Service 層     →  Fake Repository（in-memory）
Router / Handler 層  →  Fake Service + TestClient（FastAPI）/ httptest（Go）
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
- [x] ai-service — Layered Architecture + Protocol + DI + Unit / Integration Test
- [x] API Gateway — HTTP reverse proxy + JWT 驗證 + CORS + RBAC
- [x] 前端介面 — Login / Register / Dashboard / AI Lab / Notes
- [x] Docker Compose — 含 healthcheck，一鍵啟動全服務
- [x] GitHub Actions CI — Lint + Unit + Integration，三服務各自獨立 Pipeline

**規劃中**
- [ ] 服務間通訊改為 gRPC

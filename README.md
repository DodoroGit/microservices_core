# Microservices Core

這是一個學習用的個人專案，目的是透過實作來深入理解**微服務架構**與**Dependency Injection（依賴注入）**。

## 學習目標

- 理解微服務架構的設計思路與各層職責
- 實踐 Dependency Injection，達到模組解耦
- 讓 Unit Test 能夠優雅地進入專案
- 探索不同服務之間的通訊協定選擇與取捨

## 架構規劃

```
Frontend
    ↓ RESTful
API Gateway
    ↓ gRPC
Backend Services
    ↓
Databases
```

| 組件 | 說明 |
|---|---|
| Frontend | 暫以 AI 產生的基本 UI 介面為主，規劃做成個人網站或技術學習網誌 |
| API Gateway | 對外統一入口，負責路由轉發 |
| Backend Services | 多個獨立 service，目前僅有 `user-service`，後續持續擴展 |

## 目前進度

- [x] `user-service` 基本架構建立（Layered Architecture：Handler / Service / Repository）
- [x] API Gateway 建置（HTTP reverse proxy，CORS 設定）
- [x] 前端介面（React + Vite，Login / Register / Dashboard）
- [x] Docker Compose 整合（含 PostgreSQL、MongoDB、Redis、healthcheck）
- [ ] 服務間 gRPC 通訊
- [ ] Dependency Injection 解耦（interface 抽象）
- [ ] Unit Test 導入
- [ ] 第二個 backend service 擴展

## 通訊協定規劃

- **Frontend → API Gateway**：RESTful（HTTP/JSON）
- **API Gateway → Backend Service**：gRPC
- **Backend Service → Backend Service**：gRPC

## 這個專案想做什麼

目前功能方向尚未完全確定，初步規劃做成個人網站，可能包含技術學習紀錄、筆記等內容。

比起功能本身，更在意的是把這裡當作一個**技術實驗站**——持續挖掘軟體開發的深度，把每個學過的概念都落地成可以跑起來的程式。

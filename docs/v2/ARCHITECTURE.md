# 系統架構設計

## 🎯 架構原則

**簡單可靠**：避免過度設計，專注解決實際問題
**多租戶支援**：支援多教會獨立運作，預留牧區擴展
**API 優先**：前後端分離，支援 Web + Line Bot 多平台
**效能至上**：PostgreSQL + 合理索引策略，支援高併發

---

## 🏗️ 總體架構

### 分層架構

```
┌─────────────────────────────────────────────────────────────┐
│                        前端層                                │
├─────────────────────────────────────────────────────────────┤
│  Web UI (HTMX + Alpine.js + Templ + DaisyUI)                │
│  Line Bot (Flex Messages + Webhook)                         │
│  Mobile PWA (未來擴展)                                        │
└─────────────────────────────────────────────────────────────┘
                                ↕
┌─────────────────────────────────────────────────────────────┐
│                        API 層                               │
├─────────────────────────────────────────────────────────────┤
│  REST API (Go + Gin)                                        │
│  WebSocket (即時協作)                                        │
│  Authentication & Authorization                               │
└─────────────────────────────────────────────────────────────┘
                                ↕
┌─────────────────────────────────────────────────────────────┐
│                       業務邏輯層                             │
├─────────────────────────────────────────────────────────────┤
│  Service Layer (業務邏輯處理)                                │
│  Repository Pattern (數據訪問抽象)                           │
│  Domain Models (核心業務實體)                                │
└─────────────────────────────────────────────────────────────┘
                                ↕
┌─────────────────────────────────────────────────────────────┐
│                       數據持久層                             │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL (主數據庫)                                       │
│  Redis (快取 + Session)                                     │
│  SQLC (編譯時 SQL 驗證)                                      │
└─────────────────────────────────────────────────────────────┘
                                ↕
┌─────────────────────────────────────────────────────────────┐
│                       外部介面層                             │
├─────────────────────────────────────────────────────────────┤
│  Line Bot API                                               │
│  Google Calendar API                                        │
│  Email Provider (第三方)                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔧 技術棧選択理由

### Go + Gin

**選擇原因**：

- 高性能併發處理
- 簡單的語法，易於維護
- 豐富的生態系統
- 優秀的編譯優化

**替代方案**：Node.js (但 Go 的類型安全更適合)

### HTMX + Alpine.js

**選擇原因**：

- 逃離 React 的複雜性
- 伺服器端渲染，SEO 友好
- 學習曲線低，開發速度快
- 適合內容型應用

**替代方案**：React/Vue (但會增加前端複雜度)

### PostgreSQL

**選擇原因**：

- 企業級可靠性
- 豐富的數據類型支援
- 優秀的併發處理
- 成熟的生態系統

**替代方案**：MySQL (但 PostgreSQL 功能更豐富)

### SQLC

**選擇原因**：

- 編譯時 SQL 語法檢查
- 類型安全的查詢
- 避免 ORM 的複雜性
- 原生 SQL 效能

**替代方案**：GORM (但 SQLC 更精確)

---

## 🏛️ 核心模組設計

### 1. 用戶認證模組

```go
├── auth/
│   ├── handler.go      // HTTP 處理器
│   ├── service.go      // 認證業務邏輯
│   ├── middleware.go   // JWT 中介軟體
│   └── models.go       // 認證相關模型
```

**職責**：

- JWT token 生成與驗證
- 用戶登入/登出處理
- 權限檢查中介軟體

### 2. 多租戶管理

```go
├── tenant/
│   ├── context.go      // 租戶上下文
│   ├── middleware.go   // 租戶識別中介軟體
│   ├── service.go      // 租戶管理邏輯
│   └── models.go       // 租戶相關模型
```

**職責**：

- 教會 (church_id) 上下文管理
- 數據隔離確保
- 未來牧區 (district_id) 擴展準備

### 3. 服事安排核心

```go
├── ministry/
│   ├── schedule/       // 服事排程
│   ├── assignment/     // 人員指派
│   ├── validation/     // 規則驗證
│   └── notification/   // 通知系統
```

**職責**：

- 服事表 CRUD 操作
- 衝突檢測與建議
- 排程最佳化算法

### 4. 請假與換服事

```go
├── requests/
│   ├── leave/          // 請假流程
│   ├── swap/           // 換服事流程
│   ├── approval/       // 審核機制
│   └── workflow/       // 工作流引擎
```

**職責**：

- 申請流程管理
- 自動化工作流
- 狀態追蹤

### 5. 即時協作

```go
├── collaboration/
│   ├── websocket/      // WebSocket 處理
│   ├── presence/       // 在線狀態
│   ├── conflict/       // 衝突解決
│   └── activity/       // 活動記錄
```

**職責**：

- 即時編輯狀態同步
- 樂觀鎖定機制
- 編輯衝突處理

### 6. 外部整合

```go
├── integrations/
│   ├── linebot/        // Line Bot 整合
│   ├── calendar/       // Google Calendar
│   ├── email/          // Email 通知
│   └── webhook/        // Webhook 處理
```

**職責**：

- 第三方 API 串接
- 事件驅動通知
- 錯誤重試機制

---

## 🔐 安全架構

### 認證授權流程

```
1. 用戶登入 → JWT Token 生成
2. 每次請求 → Token 驗證 + 權限檢查
3. 多租戶隔離 → church_id 注入查詢
4. API 請求限流 → 防止濫用
```

### 數據安全

- **加密存儲**：用戶敏感資料加密
- **SQL 注入防護**：SQLC 參數化查詢
- **XSS 防護**：模板自動 escaping
- **CSRF 防護**：CSRF token 驗證

### 權限控制

```go
type Permission struct {
    Resource string // users, schedules, requests
    Action   string // read, write, approve, delete
    Scope    string // own, team, church
}
```

**權限矩陣**：

- `Super Admin`: 所有資源 + 系統設定
- `Leader`: 排程管理 + 審核權限
- `Vice Leader`: 部分排程 + 基本審核
- `Member`: 查看個人 + 申請請假/換服事

---

## 📊 數據架構

### 多租戶數據隔離

```sql
-- 所有業務表都包含 church_id
CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    church_id UUID NOT NULL REFERENCES churches(id),
    -- 其他欄位
);

-- 創建複合索引確保查詢效能
CREATE INDEX idx_schedules_church_date
ON schedules(church_id, scheduled_date);
```

### 讀寫分離 (未來)

```
Master DB (寫操作)
    ↓ 同步
Slave DB (讀操作)
```

### 快取策略

- **Redis Session**: 用戶登入狀態
- **Redis Cache**: 熱點數據 (當前月份服事表)
- **Application Cache**: 權限資料、教會設定

---

## 🚀 效能優化

### 數據庫優化

```sql
-- 核心查詢索引
CREATE INDEX idx_assignments_schedule_user
ON assignments(schedule_id, user_id, status);

CREATE INDEX idx_users_church_role
ON users(church_id, role, is_active);

-- 部分索引 (活躍用戶)
CREATE INDEX idx_active_users
ON users(church_id) WHERE is_active = true;
```

### API 效能

- **批量查詢**：減少 N+1 問題
- **分頁查詢**：大數據集分頁載入
- **GraphQL** (未來): 精確數據載入

### 前端優化

- **HTMX**: 部分頁面更新，減少重新渲染
- **Alpine.js**: 最小化 JavaScript bundle
- **Service Worker**: 離線支援 (PWA)

---

## 🔄 部署架構

### 容器化策略

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ministry-scheduler

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/ministry-scheduler /
EXPOSE 8080
CMD ["/ministry-scheduler"]
```

### 生產環境架構

```
Internet → Load Balancer → App Instances (3+) → Database Cluster
              ↓
         Static Assets (CDN)
              ↓
         External APIs (Line Bot, Google)
```

### 監控與日誌

- **健康檢查**：`/health` endpoint
- **Metrics**: Prometheus + Grafana
- **日誌聚合**：結構化 JSON 日誌
- **錯誤追蹤**：Sentry 集成

---

## 🛡️ 可靠性設計

### 錯誤處理

```go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 分層錯誤處理
Repository Layer → Service Layer → Handler Layer → Client
```

### 重試機制

- **外部 API 調用**：指數退避重試
- **數據庫連接**：連接池 + 自動重連
- **消息隊列**：失敗消息重新排隊

### 降級策略

- **Line Bot 服務異常**：降級到 Email 通知
- **Google Calendar 同步失敗**：記錄事件，後台重試
- **數據庫只讀模式**：禁用寫操作，顯示維護訊息

---

## 📈 擴展性考慮

### 水平擴展

- **無狀態服務**：所有狀態存儲在 Redis/DB
- **負載均衡**：多個應用實例
- **數據庫分片** (長期)：依 church_id 分片

### 垂直擴展

- **讀取優化**：讀寫分離 + 讀取副本
- **快取分層**：應用快取 + 數據庫快取
- **CDN 加速**：靜態資源分發

### 微服務演進 (未來)

```
Monolith → Module → Microservices
   ↓         ↓         ↓
  MVP    功能穩定    需要獨立擴展
```

---

## 🔧 開發工具鏈

### 程式碼品質

- **gofmt**: 程式碼格式化
- **golint**: 靜態分析
- **go test**: 單元測試
- **testify**: 測試斷言庫

### CI/CD 流程

```
Git Push → GitHub Actions → Build → Test → Deploy → Monitor
```

### 本地開發環境

```bash
# 使用 docker-compose 啟動完整環境
docker-compose up -d postgres redis
go run main.go
```

---

*「智慧人的心教訓他的口，又使他的嘴增長學問。」- 箴 16:23*


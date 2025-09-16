# 系統架構設計

## 整體架構

```
┌─────────────────┐    ┌─────────────────┐    ┌──────────────────┐
│     Line Bot    │    │  Web Interface  │    │  Google Calendar │
│       API       │    │  (HTMX+Alpine)  │    │    Integration   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬────────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │         Go Server         │
                    │     (HTTP + WebSocket)    │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │         PostgreSQL        │
                    │          Database         │
                    └───────────────────────────┘
```

## 核心組件

### Web Server (Go)

**職責**：

* HTTP API 處理
* WebSocket 連線管理
* 業務邏輯處理
* 認證授權
* 外部服務整合

**技術選擇理由**：

* Go 的並發性能適合處理多使用者連線
* 標準庫功能完整，依賴少
* 編譯後單一執行檔，部署簡單
* 靜態型別保證程式品質

### 前端架構 (HTMX + Alpine.js)

**HTMX 負責**：

* 無刷新頁面更新
* 表單提交處理
* 動態內容載入

**Alpine.js 負責**：

* 客戶端狀態管理
* 簡單的互動邏輯
* UI 元件行為

**優勢**：

* 避免複雜的 SPA 框架
* 伺服器渲染保持 SEO 友善
* 學習成本低，維護容易

### 資料存取層 (SQLC)

**設計原則**：

* 型別安全的 SQL 查詢
* 編譯時檢查 SQL 語法
* 避免 ORM 的複雜性

**範例結構**：

```go
type Store interface {
    GetUser(ctx context.Context, id int64) (User, error)
    ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error)
    CreateLeaveRequest(ctx context.Context, params CreateLeaveRequestParams) (LeaveRequest, error)
}
```

### WebSocket 協作機制

**使用場景**：

* 多人編輯時的「正在編輯」提示
* 即時通知更新

**實作原則**：

* 樂觀鎖定，避免複雜的衝突處理
* 最後編輯者獲勝
* 保持連線輕量化

```go
type CollaborationMessage struct {
    Type   string `json:"type"`
    UserID int64  `json:"user_id"`
    Data   any    `json:"data"`
}
```

## 資料流設計

### 典型使用者操作流程

1. **使用者登入**

   ```
   Browser → Go Server → Database → Go Server → Browser
   ```

2. **檢視服事安排**

   ```
   Browser → Go Server → Database → Template Engine → Browser
   ```

3. **申請請假**

   ```
   Browser → Go Server → Database → Line Bot API → 主責手機
   ```

4. **多人協作編輯**

   ```
   Browser A ──┐
               ├─→ Go Server → WebSocket → Browser B
   Browser B ──┘                      └─→ Browser A
   ```

## 權限控制

### 角色定義

```go
type Role int

const (
    RoleMember Role = iota
    RoleDeputy
    RoleLeader
    RoleAdmin
)
```

### 權限矩陣

| 功能 | 成員 | 副主責 | 主責 | 超管 |
|------|------|--------|------|------|
| 查看服事安排 | ✓ | ✓ | ✓ | ✓ |
| 申請請假 | ✓ | ✓ | ✓ | ✓ |
| 發起換服事 | ✓ | ✓ | ✓ | ✓ |
| 審核請假 | ✗ | 部分 | ✓ | ✓ |
| 編輯服事安排 | ✗ | 部分 | ✓ | ✓ |
| 系統設定 | ✗ | ✗ | ✗ | ✓ |

### 權限檢查實作

```go
func (s *Service) CanEditSchedule(userID, eventID int64) bool {
    user := s.store.GetUser(ctx, userID)
    if user.Role >= RoleLeader {
        return true
    }
    if user.Role == RoleDeputy {
        return s.store.IsDeputyForEvent(ctx, userID, eventID)
    }
    return false
}
```

## 外部整合架構

### Line Bot 整合

**Webhook 處理**：

```go
func (h *LineHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
    events := h.parseEvents(r)
    for _, event := range events {
        switch event.Type {
        case "message":
            h.handleMessage(event)
        case "postback":
            h.handlePostback(event)
        }
    }
}
```

**通知發送**：

```go
func (s *NotificationService) SendLeaveRequest(req LeaveRequest) {
    leaders := s.store.GetEventLeaders(ctx, req.EventID)
    for _, leader := range leaders {
        s.lineBot.SendMessage(leader.LineID, formatLeaveRequest(req))
    }
}
```

### Google Calendar 整合

**單向同步設計**：

```go
func (s *CalendarService) SyncUserEvents(userID int64) error {
    events := s.store.GetUserEvents(ctx, userID)
    for _, event := range events {
        if event.CalendarEventID == "" {
            // 建立新的 Calendar Event
            calEvent := s.calendar.CreateEvent(formatEvent(event))
            s.store.UpdateEventCalendarID(ctx, event.ID, calEvent.ID)
        } else {
            // 更新現有 Event
            s.calendar.UpdateEvent(event.CalendarEventID, formatEvent(event))
        }
    }
    return nil
}
```

## 錯誤處理策略

### 錯誤分類

```go
type ErrorType int

const (
    ErrorTypeValidation ErrorType = iota
    ErrorTypePermission
    ErrorTypeNotFound
    ErrorTypeInternal
)

type AppError struct {
    Type    ErrorType
    Message string
    Cause   error
}
```

### 統一錯誤回應

```go
func (h *Handler) HandleError(w http.ResponseWriter, err AppError) {
    switch err.Type {
    case ErrorTypeValidation:
        http.Error(w, err.Message, http.StatusBadRequest)
    case ErrorTypePermission:
        http.Error(w, "權限不足", http.StatusForbidden)
    case ErrorTypeNotFound:
        http.Error(w, "找不到資源", http.StatusNotFound)
    default:
        log.Error("Internal error", "error", err.Cause)
        http.Error(w, "系統錯誤", http.StatusInternalServerError)
    }
}
```

## 效能考量

### 快取策略

* 用戶權限資訊 (短期快取)
* 服事安排檢視 (條件式快取)
* 靜態資源 (長期快取)

### 資料庫最佳化

* 適當的索引設計
* 分頁載入大量資料
* 讀寫分離 (未來考慮)

### 監控指標

* 回應時間
* 錯誤率
* 資料庫連線數
* WebSocket 連線數

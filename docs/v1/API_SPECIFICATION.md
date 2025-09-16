# API 規格文件

## 設計原則

### RESTful 設計

遵循標準的 REST 設計原則，資源導向，HTTP 動詞語意明確。

### 一致性優先

所有 API 回應格式保持一致，錯誤處理統一。

### HTMX 友善

除了 JSON API，同時提供 HTML 片段回應供 HTMX 使用。

## 通用規範

### 請求格式

* Content-Type: `application/json` 或 `application/x-www-form-urlencoded`
* 認證: Session Cookie 或 JWT Token

### 回應格式

#### 成功回應

```json
{
  "success": true,
  "data": {...},
  "message": "操作成功"
}
```

#### 錯誤回應

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "輸入資料有誤",
    "details": {
      "field": "email",
      "reason": "格式不正確"
    }
  }
}
```

### HTTP 狀態碼

* `200 OK` - 成功
* `201 Created` - 建立成功
* `400 Bad Request` - 請求錯誤
* `401 Unauthorized` - 未認證
* `403 Forbidden` - 權限不足
* `404 Not Found` - 資源不存在
* `422 Unprocessable Entity` - 驗證錯誤
* `500 Internal Server Error` - 伺服器錯誤

## 認證相關 API

### POST /auth/login

使用者登入

**請求**：

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**成功回應**：

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "name": "王小明",
      "email": "user@example.com",
      "role": "member"
    },
    "token": "jwt_token_here"
  }
}
```

### POST /auth/logout

使用者登出

### GET /auth/me

獲取當前使用者資訊

## 使用者管理 API

### GET /users

獲取使用者列表 (需要管理員權限)

**參數**：

* `page` - 頁碼 (預設: 1)
* `limit` - 每頁筆數 (預設: 20)
* `role` - 角色過濾

**回應**：

```json
{
  "success": true,
  "data": {
    "users": [...],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

### GET /users/:id

獲取特定使用者資訊

### PUT /users/:id

更新使用者資訊 (需要相應權限)

**請求**：

```json
{
  "name": "王小明",
  "phone": "0912345678",
  "role": "deputy"
}
```

## 活動管理 API

### GET /events

獲取活動列表

**參數**：

* `start_date` - 開始日期 (YYYY-MM-DD)
* `end_date` - 結束日期 (YYYY-MM-DD)
* `active_only` - 只顯示啟用的活動 (預設: true)

**回應**：

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "主日崇拜",
      "event_date": "2024-12-08",
      "start_time": "09:00:00",
      "end_time": "11:00:00",
      "location": "教會",
      "services": [
        {
          "position": "音控",
          "user": {
            "id": 1,
            "name": "王小明"
          }
        }
      ]
    }
  ]
}
```

### POST /events

建立新活動 (需要主責權限)

**請求**：

```json
{
  "name": "主日崇拜",
  "description": "每週主日崇拜",
  "event_date": "2024-12-08",
  "start_time": "09:00",
  "end_time": "11:00",
  "location": "教會",
  "services": [
    {
      "position_id": 1,
      "user_id": 2
    }
  ]
}
```

### PUT /events/:id

更新活動資訊

### DELETE /events/:id

刪除活動 (軟刪除)

## 服事安排 API

### GET /events/:id/services

獲取特定活動的服事安排

### POST /events/:id/services

為活動新增服事安排

**請求**：

```json
{
  "position_id": 1,
  "user_id": 2,
  "notes": "備註說明"
}
```

### PUT /services/:id

更新服事安排

### DELETE /services/:id

移除服事安排

## 服事檢視 API

### GET /schedules/personal

個人服事時程檢視

**參數**：

* `user_id` - 使用者ID (管理員可查看他人)
* `month` - 月份 (YYYY-MM)

### GET /schedules/team

團隊服事檢視

**參數**：

* `start_date` - 開始日期
* `end_date` - 結束日期

### GET /schedules/position/:position_id

特定崗位的服事安排檢視

### GET /schedules/event-type/:type

特定活動類型的服事檢視

## 請假申請 API

### GET /leave-requests

獲取請假申請列表

**參數**：

* `status` - 狀態過濾 (pending, approved, rejected)
* `user_id` - 使用者ID過濾

**回應**：

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "name": "王小明"
      },
      "start_date": "2024-12-10",
      "end_date": "2024-12-15",
      "reason": "家庭旅遊",
      "status": "pending",
      "created_at": "2024-12-01T10:00:00Z"
    }
  ]
}
```

### POST /leave-requests

提出請假申請

**請求**：

```json
{
  "start_date": "2024-12-10",
  "end_date": "2024-12-15",
  "reason": "家庭旅遊"
}
```

### PUT /leave-requests/:id/approve

審核請假申請 (需要主責權限)

**請求**：

```json
{
  "approved": true,
  "notes": "審核備註"
}
```

## 換服事請求 API

### GET /swap-requests

獲取換服事請求列表

### POST /swap-requests

發起換服事請求

**請求**：

```json
{
  "event_service_id": 1,
  "target_user_id": 2,
  "reason": "臨時有事無法參與"
}
```

### PUT /swap-requests/:id/respond

回應換服事請求

**請求**：

```json
{
  "accepted": true,
  "notes": "可以幫忙"
}
```

### PUT /swap-requests/:id/review

主責審核換服事請求

## 通知 API

### GET /notifications

獲取使用者通知

**參數**：

* `unread_only` - 只顯示未讀通知 (預設: false)

### PUT /notifications/:id/read

標記通知為已讀

### PUT /notifications/read-all

標記所有通知為已讀

## WebSocket API

### 連線端點

`ws://domain.com/ws?token=jwt_token`

### 訊息格式

#### 客戶端 → 伺服器

**加入協作房間**：

```json
{
  "type": "join_room",
  "data": {
    "event_id": 1
  }
}
```

**編輯狀態更新**：

```json
{
  "type": "editing",
  "data": {
    "event_id": 1,
    "field": "services",
    "editing": true
  }
}
```

#### 伺服器 → 客戶端

**使用者加入通知**：

```json
{
  "type": "user_joined",
  "data": {
    "user": {
      "id": 1,
      "name": "王小明"
    },
    "event_id": 1
  }
}
```

**編輯狀態通知**：

```json
{
  "type": "user_editing",
  "data": {
    "user": {
      "id": 1,
      "name": "王小明"
    },
    "field": "services",
    "editing": true
  }
}
```

**資料更新通知**：

```json
{
  "type": "data_updated",
  "data": {
    "type": "event_service",
    "id": 1,
    "action": "updated",
    "updated_by": {
      "id": 2,
      "name": "李小華"
    }
  }
}
```

## HTMX 端點

### GET /htmx/events/:id/services

回傳服事安排的 HTML 片段

### POST /htmx/services

建立服事安排並回傳更新後的 HTML

### PUT /htmx/services/:id

更新服事安排並回傳更新後的 HTML

## 錯誤碼對照

### 驗證錯誤 (VALIDATION_ERROR)

* `REQUIRED_FIELD` - 必填欄位缺少
* `INVALID_FORMAT` - 格式不正確
* `INVALID_DATE` - 日期無效
* `DATE_RANGE_ERROR` - 日期範圍錯誤

### 業務邏輯錯誤 (BUSINESS_ERROR)

* `SCHEDULE_CONFLICT` - 服事時間衝突
* `PERMISSION_DENIED` - 權限不足
* `ALREADY_REQUESTED` - 已存在相同請求
* `INVALID_STATUS` - 狀態不正確

### 系統錯誤 (SYSTEM_ERROR)

* `DATABASE_ERROR` - 資料庫錯誤
* `EXTERNAL_SERVICE_ERROR` - 外部服務錯誤

## 分頁與排序

### 分頁參數

* `page` - 頁碼 (從 1 開始)
* `limit` - 每頁筆數 (預設 20，最大 100)

### 排序參數

* `sort` - 排序欄位
* `order` - 排序方向 (asc, desc)

範例：`GET /events?page=2&limit=10&sort=event_date&order=desc`

## 速率限制

### 限制規則

* 一般 API：100 requests/minute/IP
* 認證相關 API：10 requests/minute/IP
* WebSocket 連線：10 connections/user

### 限制回應

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "請求過於頻繁，請稍後再試"
  }
}
```

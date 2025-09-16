# REST API 規格文件

## 🎯 API 設計原則

**RESTful 設計**：遵循 REST 約定，資源導向的 URL 設計
**一致性響應**：統一的成功/錯誤響應格式
**多租戶支援**：所有 API 自動注入 `church_id` 上下文
**版本管理**：透過 URL 路徑進行版本控制 `/api/v1/...`

---

## 🔐 認證與授權

### JWT Token 認證

```http
Authorization: Bearer <jwt_token>
```

### Token 格式

```javascript
{
  "user_id": "uuid",
  "church_id": "uuid",
  "role": "leader|vice_leader|member|admin",
  "exp": timestamp,
  "iat": timestamp
}
```

### 權限層級

- **admin**: 全教會權限
- **leader**: 服事管理 + 審核權限
- **vice_leader**: 部分服事管理
- **member**: 個人服事查看 + 申請權限

---

## 📦 通用響應格式

### 成功響應

```json
{
  "success": true,
  "data": {
    // 實際數據內容
  },
  "meta": {
    "timestamp": "2025-09-16T17:50:38Z",
    "request_id": "uuid"
  }
}
```

### 錯誤響應

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "請求參數錯誤",
    "details": "start_date 不能晚於 end_date"
  },
  "meta": {
    "timestamp": "2025-09-16T17:50:38Z",
    "request_id": "uuid"
  }
}
```

### 分頁響應

```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 156,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## 🔑 認證端點

### POST `/api/v1/auth/login`

用戶登入

**請求**:

```json
{
  "email": "user@church.org",
  "password": "password123"
}
```

**響應**:

```json
{
  "success": true,
  "data": {
    "access_token": "jwt_token",
    "user": {
      "id": "uuid",
      "name": "張小美",
      "email": "user@church.org",
      "role": "leader",
      "church": {
        "id": "uuid",
        "name": "恩典教會"
      }
    }
  }
}
```

### POST `/api/v1/auth/logout`

用戶登出

**響應**:

```json
{
  "success": true,
  "message": "已成功登出"
}
```

### POST `/api/v1/auth/refresh`

更新 Token

**響應**:

```json
{
  "success": true,
  "data": {
    "access_token": "new_jwt_token"
  }
}
```

---

## 👥 用戶管理 API

### GET `/api/v1/users`

查詢用戶列表

**權限**: leader+

**查詢參數**:

- `page` (integer): 頁碼，預設 1
- `per_page` (integer): 每頁數量，預設 20，最大 100
- `role` (string): 角色篩選
- `status` (string): 狀態篩選
- `search` (string): 姓名或 email 搜索

**響應**:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "張小美",
      "email": "user@church.org",
      "phone": "0912345678",
      "role": "member",
      "status": "active",
      "last_login_at": "2025-09-15T10:30:00Z",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### GET `/api/v1/users/{id}`

查詢單一用戶詳細資訊

**權限**: 本人或 leader+

**響應**:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "張小美",
    "email": "user@church.org",
    "phone": "0912345678",
    "role": "member",
    "status": "active",
    "preferences": {
      "date_preferences": [...],
      "position_preferences": [...]
    },
    "statistics": {
      "total_assignments": 45,
      "this_month_assignments": 3,
      "completion_rate": 0.98
    }
  }
}
```

### PUT `/api/v1/users/{id}`

更新用戶資訊

**權限**: 本人或 leader+

**請求**:

```json
{
  "name": "張小美",
  "phone": "0987654321",
  "preferences": {
    "notifications": {
      "line": true,
      "email": false
    }
  }
}
```

---

## 📅 服事安排 API

### GET `/api/v1/schedules`

查詢服事安排列表

**查詢參數**:

- `start_date` (date): 開始日期
- `end_date` (date): 結束日期
- `status` (string): 狀態篩選
- `view_type` (string): 檢視類型 `personal|team|position|event`
- `user_id` (uuid): 特定用戶篩選
- `position_id` (uuid): 特定崗位篩選

**響應**:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "主日崇拜",
      "scheduled_date": "2025-09-21",
      "start_time": "09:00",
      "end_time": "11:30",
      "status": "published",
      "assignments": [
        {
          "id": "uuid",
          "user": {
            "id": "uuid",
            "name": "張小美"
          },
          "position": {
            "id": "uuid",
            "name": "音控",
            "ministry_type": "技術服事"
          },
          "status": "confirmed"
        }
      ],
      "created_by": {
        "id": "uuid",
        "name": "李主責"
      }
    }
  ]
}
```

### POST `/api/v1/schedules`

建立新的服事安排

**權限**: leader+

**請求**:

```json
{
  "event_id": "uuid",
  "title": "主日崇拜",
  "scheduled_date": "2025-09-21",
  "start_time": "09:00",
  "end_time": "11:30",
  "notes": "請提早15分鐘到場準備",
  "assignments": [
    {
      "user_id": "uuid",
      "position_id": "uuid",
      "notes": "新手，請安排資深同工協助"
    }
  ]
}
```

### PUT `/api/v1/schedules/{id}`

更新服事安排

**權限**: leader+

### DELETE `/api/v1/schedules/{id}`

刪除服事安排 (軟刪除)

**權限**: leader+

### POST `/api/v1/schedules/{id}/publish`

發佈服事安排

**權限**: leader+

**響應**:

```json
{
  "success": true,
  "data": {
    "published_at": "2025-09-16T17:50:38Z",
    "notifications_sent": 5
  }
}
```

---

## 🙋 請假申請 API

### GET `/api/v1/leave-requests`

查詢請假申請列表

**查詢參數**:

- `status` (string): 狀態篩選
- `user_id` (uuid): 用戶篩選 (leader+ 可查看全部)
- `start_date`, `end_date`: 日期範圍

**響應**:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user": {
        "id": "uuid",
        "name": "張小美"
      },
      "start_date": "2025-09-20",
      "end_date": "2025-09-22",
      "reason": "家族聚會",
      "status": "pending",
      "created_at": "2025-09-16T10:00:00Z"
    }
  ]
}
```

### POST `/api/v1/leave-requests`

提交請假申請

**請求**:

```json
{
  "start_date": "2025-09-20",
  "end_date": "2025-09-22",
  "reason": "家族聚會，需要回南部"
}
```

### PUT `/api/v1/leave-requests/{id}/review`

審核請假申請

**權限**: leader+

**請求**:

```json
{
  "status": "approved|rejected",
  "review_notes": "審核意見"
}
```

---

## 🔄 換服事申請 API

### POST `/api/v1/swap-requests`

發起換服事請求

**請求**:

```json
{
  "assignment_id": "uuid",
  "target_user_id": "uuid", // 可選，不指定則通知該崗位所有成員
  "reason": "臨時有事無法參加"
}
```

### GET `/api/v1/swap-requests`

查詢換服事請求

**響應**:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "assignment": {
        "schedule": {
          "title": "主日崇拜",
          "scheduled_date": "2025-09-21"
        },
        "position": {
          "name": "音控"
        }
      },
      "requestor": {
        "name": "張小美"
      },
      "target_user": {
        "name": "李大華"
      },
      "reason": "臨時有事",
      "status": "pending",
      "created_at": "2025-09-16T15:00:00Z"
    }
  ]
}
```

### PUT `/api/v1/swap-requests/{id}/respond`

回應換服事請求

**請求**:

```json
{
  "status": "accepted|rejected",
  "notes": "回應說明"
}
```

### PUT `/api/v1/swap-requests/{id}/approve`

主責審核換服事

**權限**: leader+

**請求**:

```json
{
  "status": "approved|rejected",
  "review_notes": "審核意見"
}
```

---

## 📊 統計分析 API

### GET `/api/v1/statistics/overview`

教會服事總覽統計

**權限**: leader+

**響應**:

```json
{
  "success": true,
  "data": {
    "active_users": 45,
    "this_month_schedules": 12,
    "pending_requests": 3,
    "completion_rate": 0.98,
    "top_positions": [
      {
        "position": "音控",
        "assignment_count": 24
      }
    ],
    "user_participation": [
      {
        "user_name": "張小美",
        "assignment_count": 8,
        "completion_rate": 1.0
      }
    ]
  }
}
```

### GET `/api/v1/statistics/user/{user_id}`

個人服事統計

**權限**: 本人或 leader+

**響應**:

```json
{
  "success": true,
  "data": {
    "total_assignments": 45,
    "completed_assignments": 44,
    "completion_rate": 0.98,
    "this_year_assignments": 24,
    "favorite_positions": [
      {
        "position": "音控",
        "count": 18
      }
    ],
    "monthly_distribution": [
      {
        "month": "2025-01",
        "count": 4
      }
    ]
  }
}
```

---

## 🔔 通知管理 API

### GET `/api/v1/notifications`

查詢個人通知

**查詢參數**:

- `unread_only` (boolean): 只顯示未讀通知
- `type` (string): 通知類型篩選

**響應**:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "type": "assignment_reminder",
      "title": "服事提醒",
      "message": "您在明天 (9/21) 主日崇拜有音控服事",
      "data": {
        "schedule_id": "uuid",
        "assignment_id": "uuid"
      },
      "read_at": null,
      "created_at": "2025-09-20T08:00:00Z"
    }
  ]
}
```

### PUT `/api/v1/notifications/{id}/read`

標記通知為已讀

### PUT `/api/v1/notifications/read-all`

標記所有通知為已讀

---

## 🔗 外部整合 API

### Line Bot Webhook

#### POST `/api/v1/webhooks/line`

Line Bot 事件接收

**請求** (Line Bot 格式):

```json
{
  "events": [
    {
      "type": "message",
      "message": {
        "type": "text",
        "text": "我的服事"
      },
      "source": {
        "userId": "line_user_id"
      }
    }
  ]
}
```

#### 支援的 Line Bot 指令

- `我的服事` → 查看個人服事安排
- `團隊服事` → 查看團隊服事概況
- `請假申請` → 發起請假申請流程
- `換服事` → 發起換服事請求

### Google Calendar API

#### POST `/api/v1/integrations/google/authorize`

Google Calendar 授權

**請求**:

```json
{
  "authorization_code": "google_auth_code"
}
```

#### POST `/api/v1/integrations/google/sync`

手動同步行事曆

**響應**:

```json
{
  "success": true,
  "data": {
    "synced_events": 5,
    "updated_events": 2,
    "deleted_events": 1
  }
}
```

---

## 🚨 錯誤代碼定義

| 錯誤代碼 | HTTP 碼 | 說明 |
|---------|---------|------|
| `UNAUTHORIZED` | 401 | 未授權存取 |
| `FORBIDDEN` | 403 | 權限不足 |
| `NOT_FOUND` | 404 | 資源不存在 |
| `VALIDATION_ERROR` | 400 | 請求參數錯誤 |
| `CONFLICT` | 409 | 資源衝突 (重複建立等) |
| `RATE_LIMITED` | 429 | 請求頻率過高 |
| `INTERNAL_ERROR` | 500 | 系統內部錯誤 |
| `SCHEDULE_CONFLICT` | 400 | 服事時間衝突 |
| `ASSIGNMENT_CONFLICT` | 400 | 服事指派衝突 |
| `LEAVE_OVERLAP` | 400 | 請假時間重疊 |

---

## 📝 API 版本管理

### 版本策略

- **URL 路徑版本**: `/api/v1/`, `/api/v2/`
- **向後相容**: 舊版本 API 至少維護 6 個月
- **破壞性變更**: 必須提升主版本號

### 版本變更通知

```http
X-API-Version: 1.0
X-Deprecation-Date: 2025-12-31
X-Sunset-Date: 2026-06-30
```

---

## 🚀 效能與限制

### API 速率限制

- **一般用戶**: 每分鐘 60 requests
- **管理員**: 每分鐘 120 requests
- **Line Bot**: 每分鐘 300 requests

### 響應時間目標

- **查詢操作**: < 200ms
- **寫入操作**: < 500ms
- **批量操作**: < 2s

### 數據限制

- **分頁最大值**: 100 筆/頁
- **批量操作**: 最多 50 筆/請求
- **檔案上傳**: 最大 5MB
- **請求 Body**: 最大 1MB

---

*「凡事都要規規矩矩地按著次序行。」- 林前 14:40*


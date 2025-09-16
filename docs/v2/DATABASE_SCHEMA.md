# 資料庫設計規格

## 🎯 設計原則

**多租戶隔離**：所有業務表包含 `church_id`，確保數據完全隔離
**索引策略**：基於實際查詢模式設計索引，避免過度索引
**數據一致性**：外鍵約束 + 事務管理確保數據完整性
**擴展性考慮**：預留 `district_id` 支援未來牧區分割需求

---

## 🏗️ 核心實體關係

```
Churches (教會)
    ↓ 1:N
Users (用戶)
    ↓ M:N
Assignments (服事指派) ← Schedule (服事安排)
    ↓ 1:N                      ↓ 1:N
LeaveRequests              SwapRequests
(請假申請)                  (換服事申請)
```

---

## 📊 表結構設計

### 1. 核心主體表

#### Churches (教會)

```sql
CREATE TABLE churches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    settings JSONB DEFAULT '{}',
    timezone VARCHAR(50) DEFAULT 'Asia/Taipei',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_churches_name ON churches(name) WHERE is_active = true;
```

#### Districts (牧區 - 未來擴展)

```sql
CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_districts_church_name
ON districts(church_id, name) WHERE is_active = true;
```

### 2. 用戶管理

#### Users (用戶)

```sql
CREATE TYPE user_role AS ENUM ('super_admin', 'admin', 'leader', 'vice_leader', 'member');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    district_id UUID REFERENCES districts(id),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    line_user_id VARCHAR(255),
    role user_role DEFAULT 'member',
    status user_status DEFAULT 'active',
    google_calendar_token JSONB,
    preferences JSONB DEFAULT '{}',
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 核心索引
CREATE UNIQUE INDEX idx_users_church_email
ON users(church_id, email) WHERE status != 'suspended';

CREATE INDEX idx_users_church_role_status
ON users(church_id, role, status);

CREATE INDEX idx_users_line_id
ON users(line_user_id) WHERE line_user_id IS NOT NULL;
```

#### User Sessions (用戶會話)

```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_sessions_token ON user_sessions(token_hash);
CREATE INDEX idx_sessions_user_expires ON user_sessions(user_id, expires_at);
```

### 3. 服事管理

#### Ministry Types (服事類型)

```sql
CREATE TABLE ministry_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3B82F6',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_ministry_types_church_name
ON ministry_types(church_id, name) WHERE is_active = true;
```

#### Positions (服事崗位)

```sql
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    ministry_type_id UUID NOT NULL REFERENCES ministry_types(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    required_skills TEXT[],
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_positions_church_ministry
ON positions(church_id, ministry_type_id) WHERE is_active = true;
```

#### Events (活動)

```sql
CREATE TYPE event_status AS ENUM ('draft', 'published', 'completed', 'cancelled');
CREATE TYPE event_recurrence AS ENUM ('none', 'weekly', 'monthly', 'custom');

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    district_id UUID REFERENCES districts(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    recurrence event_recurrence DEFAULT 'none',
    recurrence_pattern JSONB,
    status event_status DEFAULT 'draft',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 查詢優化索引
CREATE INDEX idx_events_church_date_status
ON events(church_id, start_time, status);

CREATE INDEX idx_events_creator
ON events(created_by, created_at);
```

### 4. 服事安排

#### Schedules (服事安排)

```sql
CREATE TYPE schedule_status AS ENUM ('draft', 'published', 'completed', 'cancelled');

CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    scheduled_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    notes TEXT,
    status schedule_status DEFAULT 'draft',
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER DEFAULT 1
);

-- 核心查詢索引
CREATE INDEX idx_schedules_church_date
ON schedules(church_id, scheduled_date);

CREATE INDEX idx_schedules_event_date
ON schedules(event_id, scheduled_date);
```

#### Assignments (服事指派)

```sql
CREATE TYPE assignment_status AS ENUM ('assigned', 'confirmed', 'declined', 'replaced', 'completed');

CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    status assignment_status DEFAULT 'assigned',
    notes TEXT,
    assigned_by UUID NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- 防重複指派
CREATE UNIQUE INDEX idx_assignments_schedule_position
ON assignments(schedule_id, position_id)
WHERE status IN ('assigned', 'confirmed');

-- 查詢優化
CREATE INDEX idx_assignments_user_status
ON assignments(user_id, status, assigned_at);

CREATE INDEX idx_assignments_schedule_church
ON assignments(schedule_id, church_id);
```

### 5. 請假與換服事

#### Leave Requests (請假申請)

```sql
CREATE TYPE request_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled');

CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status request_status DEFAULT 'pending',
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    review_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 衝突檢查索引
CREATE INDEX idx_leave_requests_user_dates
ON leave_requests(user_id, start_date, end_date);

-- 審核查詢索引
CREATE INDEX idx_leave_requests_church_status
ON leave_requests(church_id, status, created_at);
```

#### Swap Requests (換服事申請)

```sql
CREATE TYPE swap_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled', 'approved', 'completed');

CREATE TABLE swap_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    requestor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id UUID REFERENCES users(id),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    reason TEXT,
    status swap_status DEFAULT 'pending',
    accepted_by UUID REFERENCES users(id),
    accepted_at TIMESTAMP WITH TIME ZONE,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    review_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_swap_requests_assignment
ON swap_requests(assignment_id, status);

CREATE INDEX idx_swap_requests_target_user
ON swap_requests(target_user_id, status) WHERE target_user_id IS NOT NULL;

CREATE INDEX idx_swap_requests_church_status
ON swap_requests(church_id, status, created_at);
```

### 6. 偏好設定

#### User Preferences (用戶偏好)

```sql
CREATE TYPE preference_type AS ENUM ('date_preference', 'position_preference', 'pairing_preference');
CREATE TYPE preference_value AS ENUM ('preferred', 'acceptable', 'avoid');

CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    preference_type preference_type NOT NULL,
    target_id UUID,  -- 可以是 position_id, user_id, 或其他 ID
    preference_value preference_value NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 防重複偏好設定
CREATE UNIQUE INDEX idx_user_preferences_unique
ON user_preferences(user_id, preference_type, target_id);

CREATE INDEX idx_user_preferences_church_type
ON user_preferences(church_id, preference_type);
```

### 7. 協作與通知

#### Activity Logs (活動記錄)

```sql
CREATE TYPE activity_type AS ENUM (
    'schedule_created', 'schedule_updated', 'schedule_published',
    'assignment_created', 'assignment_updated', 'assignment_confirmed',
    'leave_requested', 'leave_approved', 'swap_requested', 'swap_completed'
);

CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type activity_type NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_church_time
ON activity_logs(church_id, created_at DESC);

CREATE INDEX idx_activity_logs_entity
ON activity_logs(entity_type, entity_id);

-- 分區表 (按月分區，自動清理舊數據)
-- ALTER TABLE activity_logs ENABLE ROW LEVEL SECURITY;
```

#### Notifications (通知)

```sql
CREATE TYPE notification_type AS ENUM (
    'assignment_reminder', 'leave_request', 'swap_request', 'schedule_update'
);
CREATE TYPE delivery_status AS ENUM ('pending', 'sent', 'delivered', 'failed');

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    church_id UUID NOT NULL REFERENCES churches(id) ON DELETE CASCADE,
    notification_type notification_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}',
    line_status delivery_status DEFAULT 'pending',
    email_status delivery_status DEFAULT 'pending',
    read_at TIMESTAMP WITH TIME ZONE,
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_unread
ON notifications(user_id, read_at, created_at)
WHERE read_at IS NULL;

CREATE INDEX idx_notifications_pending_delivery
ON notifications(scheduled_at, line_status, email_status)
WHERE line_status = 'pending' OR email_status = 'pending';
```

---

## 🚀 索引策略總覽

### 主查詢索引

```sql
-- 服事表查詢 (最頻繁)
CREATE INDEX idx_assignments_schedule_user_status
ON assignments(schedule_id, user_id, status);

-- 用戶服事歷史查詢
CREATE INDEX idx_assignments_user_date_range
ON assignments(user_id, assigned_at)
WHERE status IN ('confirmed', 'completed');

-- 教會服事統計查詢
CREATE INDEX idx_schedules_church_date_range
ON schedules(church_id, scheduled_date, status)
WHERE status = 'published';
```

### 複合查詢索引

```sql
-- 衝突檢測索引
CREATE INDEX idx_conflict_detection
ON assignments(user_id, church_id, assigned_at, status)
WHERE status IN ('assigned', 'confirmed');

-- 請假期間檢查
CREATE INDEX idx_leave_overlap_check
ON leave_requests(user_id, start_date, end_date, status)
WHERE status = 'approved';
```

### 部分索引 (節省空間)

```sql
-- 只索引活躍用戶
CREATE INDEX idx_active_users_church
ON users(church_id, role, name)
WHERE status = 'active';

-- 只索引待處理請求
CREATE INDEX idx_pending_requests
ON leave_requests(church_id, created_at)
WHERE status = 'pending';
```

---

## 🔐 數據安全設計

### Row Level Security (RLS)

```sql
-- 啟用 RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE schedules ENABLE ROW LEVEL SECURITY;
ALTER TABLE assignments ENABLE ROW LEVEL SECURITY;

-- 多租戶隔離政策
CREATE POLICY church_isolation_policy ON users
FOR ALL TO application_user
USING (church_id = current_setting('app.church_id')::UUID);

-- 權限基礎政策
CREATE POLICY user_data_policy ON assignments
FOR ALL TO application_user
USING (
    church_id = current_setting('app.church_id')::UUID
    AND (
        user_id = current_setting('app.user_id')::UUID
        OR EXISTS (
            SELECT 1 FROM users
            WHERE id = current_setting('app.user_id')::UUID
            AND role IN ('admin', 'leader', 'super_admin')
        )
    )
);
```

### 敏感數據加密

```sql
-- 加密擴展
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 敏感欄位加密
CREATE OR REPLACE FUNCTION encrypt_sensitive_data()
RETURNS TEXT AS $$
BEGIN
    RETURN crypt($1, gen_salt('bf'));
END;
$$ LANGUAGE plpgsql;
```

---

## 📊 數據完整性約束

### 業務規則約束

```sql
-- 服事時間合理性檢查
ALTER TABLE schedules ADD CONSTRAINT check_time_range
CHECK (start_time < end_time);

-- 請假日期合理性檢查
ALTER TABLE leave_requests ADD CONSTRAINT check_leave_dates
CHECK (start_date <= end_date);

-- 用戶角色與教會關聯檢查
ALTER TABLE users ADD CONSTRAINT check_admin_church_limit
CHECK (
    role != 'super_admin' OR church_id IS NULL
);
```

### 級聯刪除策略

```sql
-- 軟刪除主要實體
ALTER TABLE churches ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- 硬刪除關聯數據 (ON DELETE CASCADE 已設定)
-- schedules → assignments (級聯刪除)
-- users → assignments (級聯刪除)
```

---

## 🔧 維護與優化

### 自動化維護任務

```sql
-- 清理過期 sessions (每日執行)
DELETE FROM user_sessions
WHERE expires_at < NOW() - INTERVAL '7 days';

-- 清理舊活動記錄 (每月執行)
DELETE FROM activity_logs
WHERE created_at < NOW() - INTERVAL '1 year';

-- 統計數據更新 (每日執行)
REFRESH MATERIALIZED VIEW church_statistics;
```

### 效能監控查詢

```sql
-- 慢查詢識別
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements
WHERE mean_time > 100
ORDER BY mean_time DESC;

-- 索引使用情況
SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read
FROM pg_stat_user_indexes
WHERE idx_scan = 0;
```

### 數據備份策略

```bash
# 每日增量備份
pg_dump --schema-only ministry_scheduler > schema_backup.sql
pg_dump --data-only --inserts ministry_scheduler > data_backup.sql

# 每週完整備份
pg_basebackup -D /backup/weekly -Ft -z
```

---

## 📈 擴展性考慮

### 讀寫分離準備

```sql
-- 建立讀取視圖
CREATE VIEW user_assignments_view AS
SELECT u.name, u.email, s.title, s.scheduled_date, p.name as position
FROM assignments a
JOIN users u ON a.user_id = u.id
JOIN schedules s ON a.schedule_id = s.id
JOIN positions p ON a.position_id = p.id
WHERE a.status IN ('confirmed', 'completed');
```

### 分區表準備 (大數據量)

```sql
-- 按月分區 activity_logs
CREATE TABLE activity_logs_y2025m01
PARTITION OF activity_logs
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- 自動分區創建
CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    start_date date;
    end_date date;
    table_name text;
BEGIN
    start_date := date_trunc('month', CURRENT_DATE + interval '1 month');
    end_date := start_date + interval '1 month';
    table_name := 'activity_logs_y' || to_char(start_date, 'YYYY') || 'm' || to_char(start_date, 'MM');

    EXECUTE format('CREATE TABLE %I PARTITION OF activity_logs FOR VALUES FROM (%L) TO (%L)',
                   table_name, start_date, end_date);
END;
$$ LANGUAGE plpgsql;
```

---

*「凡事都要規規矩矩地按著次序行。」- 林前 14:40*


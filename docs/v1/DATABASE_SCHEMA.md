# 資料庫架構設計

## 設計原則

### 正規化但不過度

避免過度正規化導致查詢複雜，保持適當的資料冗餘以提升效能。

### 索引策略優先

每個表都要明確定義索引策略，避免慢查詢。

### 樂觀鎖定

使用 `updated_at` 欄位實作樂觀鎖定，避免複雜的悲觀鎖定機制。

## 核心表結構

### users - 使用者表

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    line_user_id VARCHAR(100) UNIQUE,
    role_type INT NOT NULL DEFAULT 0, -- 0:成員 1:副主責 2:主責 3:超管
    google_calendar_id VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_line_id ON users(line_user_id);
CREATE INDEX idx_users_role ON users(role_type);
```

**索引策略**：

* `email` - 登入查詢
* `line_user_id` - Line Bot 整合
* `role_type` - 權限過濾查詢

### events - 活動/聚會表

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    event_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    location VARCHAR(200),
    is_recurring BOOLEAN DEFAULT false,
    recurrence_pattern VARCHAR(50), -- weekly, monthly 等
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_active ON events(is_active) WHERE is_active = true;
CREATE UNIQUE INDEX idx_events_name_date ON events(name, event_date);
```

**索引策略**：

* `event_date` - 按日期查詢服事
* `is_active` - 部分索引，只索引活躍事件
* `name + event_date` - 防止重複建立相同活動

### service_positions - 服事崗位表

```sql
CREATE TABLE service_positions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    required_skills TEXT,
    max_concurrent INT DEFAULT 1, -- 同時段最多幾人
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_positions_name ON service_positions(name) WHERE is_active = true;
```

### event_services - 活動服事安排表

```sql
CREATE TABLE event_services (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    position_id INT REFERENCES service_positions(id),
    user_id INT REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'confirmed', -- confirmed, pending, cancelled
    notes TEXT,
    assigned_by INT REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_event_services_event ON event_services(event_id);
CREATE INDEX idx_event_services_user ON event_services(user_id);
CREATE INDEX idx_event_services_position ON event_services(position_id);
CREATE UNIQUE INDEX idx_event_services_unique ON event_services(event_id, position_id, user_id)
WHERE status != 'cancelled';
```

**索引策略**：

* `event_id` - 查詢特定活動的所有服事安排
* `user_id` - 查詢個人服事時程
* `position_id` - 查詢特定崗位的安排情況
* 組合唯一索引防止重複安排

### leave_requests - 請假申請表

```sql
CREATE TABLE leave_requests (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    reviewed_by INT REFERENCES users(id),
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_leave_requests_user ON leave_requests(user_id);
CREATE INDEX idx_leave_requests_date_range ON leave_requests(start_date, end_date);
CREATE INDEX idx_leave_requests_status ON leave_requests(status) WHERE status = 'pending';
```

**索引策略**：

* `user_id` - 查詢個人請假記錄
* `start_date, end_date` - 日期範圍查詢
* `status` - 部分索引，只索引待審核的申請

### swap_requests - 換服事請求表

```sql
CREATE TABLE swap_requests (
    id SERIAL PRIMARY KEY,
    event_service_id INT REFERENCES event_services(id),
    requester_id INT REFERENCES users(id),
    target_user_id INT REFERENCES users(id),
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, accepted, rejected, cancelled
    responded_at TIMESTAMP,
    reviewed_by INT REFERENCES users(id), -- 主責審核
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_swap_requests_service ON swap_requests(event_service_id);
CREATE INDEX idx_swap_requests_requester ON swap_requests(requester_id);
CREATE INDEX idx_swap_requests_target ON swap_requests(target_user_id);
CREATE INDEX idx_swap_requests_status ON swap_requests(status) WHERE status IN ('pending', 'accepted');
```

## 使用者偏好設定表

### user_date_preferences - 日期偏好表

```sql
CREATE TABLE user_date_preferences (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    event_type VARCHAR(50), -- 主日, 禱告會, 特會 等
    preferred_frequency INT, -- 每月幾次
    avoid_consecutive BOOLEAN DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_date_prefs_user_type ON user_date_preferences(user_id, event_type);
```

### user_pairing_preferences - 搭配偏好表

```sql
CREATE TABLE user_pairing_preferences (
    id SERIAL PRIMARY KEY,
    user_a_id INT REFERENCES users(id),
    user_b_id INT REFERENCES users(id),
    preference_type VARCHAR(20), -- prefer, avoid, neutral
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_pairing_prefs_users ON user_pairing_preferences(
    LEAST(user_a_id, user_b_id),
    GREATEST(user_a_id, user_b_id)
);
```

## 活動記錄表

### activity_logs - 活動記錄表

```sql
CREATE TABLE activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    target_type VARCHAR(50), -- event_service, leave_request, swap_request
    target_id INT NOT NULL,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_target ON activity_logs(target_type, target_id);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
```

## 通知相關表

### notifications - 通知表

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50), -- leave_request, swap_request, schedule_change
    reference_type VARCHAR(50),
    reference_id INT,
    is_read BOOLEAN DEFAULT false,
    sent_via_line BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;
CREATE INDEX idx_notifications_type ON notifications(type);
```

## 成就系統表

### achievements - 成就定義表

```sql
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    badge_icon VARCHAR(100),
    criteria JSONB, -- 達成條件
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### user_achievements - 使用者成就表

```sql
CREATE TABLE user_achievements (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    achievement_id INT REFERENCES achievements(id),
    earned_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_user_achievements ON user_achievements(user_id, achievement_id);
```

## 系統設定表

### system_settings - 系統設定表

```sql
CREATE TABLE system_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    updated_by INT REFERENCES users(id),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 重要查詢範例

### 1. 查詢使用者某月份的服事安排

```sql
SELECT
    e.name as event_name,
    e.event_date,
    e.start_time,
    sp.name as position_name,
    es.notes,
    es.status
FROM event_services es
JOIN events e ON es.event_id = e.id
JOIN service_positions sp ON es.position_id = sp.id
WHERE es.user_id = $1
    AND e.event_date >= $2
    AND e.event_date < $3
    AND es.status = 'confirmed'
ORDER BY e.event_date, e.start_time;
```

### 2. 檢查服事衝突

```sql
SELECT COUNT(*)
FROM event_services es1
JOIN event_services es2 ON es1.event_id = es2.event_id
WHERE es1.user_id = $1
    AND es2.user_id = $1
    AND es1.id != es2.id
    AND es1.status = 'confirmed'
    AND es2.status = 'confirmed';
```

### 3. 查詢待審核的請假申請

```sql
SELECT
    lr.*,
    u.name as requester_name,
    u.email
FROM leave_requests lr
JOIN users u ON lr.user_id = u.id
WHERE lr.status = 'pending'
ORDER BY lr.created_at;
```

## 資料庫遷移策略

### 版本控制

每個遷移檔案使用時間戳記命名：

```
001_20241201_create_users_table.sql
002_20241201_create_events_table.sql
003_20241202_add_user_preferences.sql
```

### 回滾計畫

每個遷移都要有對應的回滾腳本：

```
001_20241201_create_users_table.down.sql
```

### 資料完整性

重要的外鍵約束和檢查約束：

```sql
-- 確保 end_time > start_time
ALTER TABLE events ADD CONSTRAINT chk_events_time_order
CHECK (end_time > start_time);

-- 確保請假結束日期 >= 開始日期
ALTER TABLE leave_requests ADD CONSTRAINT chk_leave_date_order
CHECK (end_date >= start_date);
```

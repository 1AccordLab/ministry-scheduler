# 教會服事表管理系統 (Ministry-scheduler)

## 設計目標 (Design Goal)

建立一套直覺、高效的服事管理系統，幫助教會輕鬆解決行政負擔，提升人力管理效率。

## 設計原則 (Design Principles)

- **使用者體驗為核心 (User Experience Focused)**: Prioritize intuitive and user-friendly design.
- **安全且即時協作 (Secure & Real-time Collaboration)**: Ensure secure, real-time collaboration with conflict resolution.
- **高度彈性 (Highly Flexible)**: Adaptable to diverse church needs.
- **擴充性、易維護性 (Scalable & Maintainable)**: Designed for future growth and easy maintenance.

## 系統角色與權限管理 (Roles and Permissions)

- **系統管理員 (Admin)**:
  - 最高權限 (Full system access)
  - 人員、權限調控 (User and permission management)
- **主責 (Leader)**:
  - 核准服事變動 (Approve ministry changes)
  - 查看統計數據、排班情況 (View statistics and schedules)
  - 管理特定服事或團隊 (Manage specific ministries or teams)
- **一般服事人員 (User)**:
  - 查看本人服事表 (View personal schedule)
  - 查看他人服事表 (View other's schedules)
  - 申請休假 (Request leave)
  - 申請換服事 (Request ministry swaps)
  - 接受 / 拒絕換服事邀請 (Accept/Reject swap requests)

## 功能規劃藍圖 (Feature Roadmap)

### 1. 基本核心功能 (Core Features)

#### 多人協作模式 (Multi-user Collaboration)

- 不同的主責可在同一時間安排大家的服事 (Multiple leaders can schedule ministries simultaneously).

#### 防呆與資料鎖定機制 (Error Prevention & Data Locking)

- 使用者可事先安排休假日期 (Users can pre-schedule leave dates).
- 將”無法排服事的人”上鎖 (Lock users with scheduling constraints):
  - For 主日聚會, 指定 A 一個月最多服事兩次, 指定 B 一個月最多服事一次 (A can serve max 2 times/month, B max 1 time/month).
  - A 已經被排到小提琴，他就不能再排到其他服事 (If A is scheduled for violin, they cannot be scheduled for other ministries).
  - 當 A 被排到服事時，B 就不能被排到特定服事 ( 可能 A、B 水火不容 ) (If A is scheduled, B cannot be scheduled for a specific ministry due to incompatibility).
- 將”無法滿足條件的群組”上鎖 (Lock incompatible groups):
  - 分為 A、B 兩樂團 (Two orchestras: A and B)
    - 當 A 樂團中的某幾個人已被排到服事時，就不能排 A 樂團服事 (If some members of orchestra A are scheduled, the entire orchestra A cannot be scheduled).

#### 多樣服事表支持 (Diverse Ministry Schedule Support)

- 主日聚會 (Saturday/Sunday Worship)
- SC 服事 (SC)
- 禱告會 / 領夜 (Prayer Meeting / Leader Night)
- 特殊活動 (Special Events): 成長班 (Growth Class), 聖誕節 (Christmas), 復活節 (Easter), G1 等 (etc.)

#### 查詢模式 (Query Modes)

- 個人模式 (Personal Mode): 登入即可看到個人服事清單 (Login to view personal ministry list).
- 人員 / 群組模式 (Person/Group Mode)
- 服事項目模式 (Ministry Item Mode)
- 時間週期檢視 (Time Period View): 每日 (Daily) / 每週 (Weekly) / 每月 (Monthly)

### 2. 進階管理功能 (Advanced Management Features)

#### 行事曆整合 (Calendar Sync & Export)

- 整合 Google 行事曆，即時更新個人服事行程 (Integrate with Google Calendar for real-time updates).

#### Line 整合 (Line Integration)

- 提醒當周服事人員 (Remind weekly ministry members).
- 在 網頁 / Line 上申請換服事，自動發訊息給主責，主責直接可以在 line 上審核申請 (Request ministry swaps via web/Line, notify leader, and allow approval on Line).
- 當需要找人幫忙代服事時，自動詢問所有可能的人員，只要其中一個人願意幫忙，就取消所有其他發出去的詢問，並把申請轉給主責審核 (Automatically find substitutes via Line, cancel other requests upon acceptance, and forward to leader for approval).

#### 版本控制（歷史編輯紀錄）(Version Control - History Tracking)

### 3. 鼓勵與激勵性功能 (Encouragement & Motivation)

- 設計服事成就與獎勵系統，例如徽章、成就系統，幫助社群互動與認同感 (Design a ministry achievement and reward system, such as badges and achievements, to promote community interaction and identity).

## 技術採用與建議架構模式 (Technology Stack & Architecture)

### 主體技術 (Core Technologies)

- Go + Fiber ( 高效能、輕量框架 ) (High-performance, lightweight framework)
- PostgreSQL + Sqlc ( 強大且易維護的資料庫互動 ) (Robust and maintainable database interaction)
- Templ + Htmx + Alpine.js + TailwindCSS ( 輕量前端互動 ) (Lightweight frontend interactivity)

### 最佳實踐&架構建議 (Best Practices & Architecture Recommendations)

- 採用 RESTful API 架構，使系統可擴展性高 (RESTful API for high scalability).
- WebSocket 提供即時多人協作與提示功能 (WebSocket for real-time collaboration and notifications).


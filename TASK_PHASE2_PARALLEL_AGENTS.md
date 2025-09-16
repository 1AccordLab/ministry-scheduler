# TASK Phase 2: 模組化與並行子代理系統

**任務類型**: 架構重構與擴展
**優先級**: 高（依賴 Phase 1 完成）
**前置條件**: AGENTS.md (Phase 1) 完成並測試穩定

---

## 🎯 任務目標

將 Phase 1 的 AGENTS.md 從「單一 AI 通用指南」重構為「模組化 AI 開發團隊」：

1. **拆分 AGENTS.md**：
   - 核心部分保留在 AGENTS.md（orchestrator 角色）
   - 專職功能拆分到 `.opencode/agent/*.md`（subagent 專職檔案）

2. **定義 Subagent 規格**：
   - 每個 subagent 的職責、輸入/輸出、觸發時機、依賴關係
   - Subagent 之間的協作流程

3. **並行執行策略**：
   - 定義何時啟動「同一個 subagent 的多個實例」並行執行
   - 定義並行結果的比對與整合機制

---

## 📋 核心任務

### 任務 1：定義 AGENTS.md 的新角色定位

**要做什麼**：

重新定位 AGENTS.md，從「通用指南」變成「Orchestrator + 共通原則」。

**需要定義**：

1. **保留在 AGENTS.md 的內容**：
   - 所有 subagent 共通的原則（Linus 哲學、誠實第一、信心水準、可持續設計等）
   - AI 驅動開發的完整生命週期流程
   - Subagent 協作的總體策略
   - Orchestrator 的職責（何時召喚哪個 subagent、何時召喚人類）

2. **移除/拆分的內容**：
   - 專職功能的詳細執行細節（拆到各 subagent）
   - 特定模組的自我審查清單（拆到各 subagent）

3. **新增的內容**：
   - Subagent 總覽與協作流程圖
   - 何時啟動並行執行
   - 如何整合並行結果

---

### 任務 2：設計 Subagent 架構

**要做什麼**：

定義每個 subagent 的規格與檔案格式。

**需要定義的 Subagents**（參考 TASK_REFACTOR_AGENTS.md）：

1. **linus-analyst** - Linus 五層分析
2. **requirements-designer** - 撰寫 requirements.md & design.md
3. **feature-splitter** - 拆分 feature 為可執行任務（產生 tasks.md）
4. **task-author** - 制定任務執行計畫（產生 task.md、changes.md、review.md）
5. **linus-tests-coder** - 撰寫測試（TDD）
6. **linus-coder** - 程式碼實作
7. **debugger** - 除錯與錯誤追蹤
8. **code-reviewer** - 程式碼審查
9. **docs-reviewer** - 文件審查
10. **docs-writer** - 文件撰寫
11. **cross-checker** - 交叉驗證與仲裁（整合並行結果）

**每個 Subagent 檔案應包含**：

- 職責定義
- 輸入（需要什麼資訊/檔案）
- 輸出（產出什麼檔案/決策）
- 觸發時機（何時啟動這個 subagent）
- 依賴關係（依賴哪些前置 subagent 的輸出）
- 執行流程（這個 subagent 的內部執行步驟）
- 自我審查清單
- 何時召喚人類

---

### 任務 3：定義並行執行策略

**要做什麼**：

定義何時啟動「同一個 subagent 的多個實例」並行執行。

**需要定義的並行場景**：

1. **docs-reviewer 並行**：
   - 何時：文件審查階段
   - 為何：降低 AI 幻覺、context 缺失、理解錯誤的風險
   - 如何：2+ 個 docs-reviewer 審查同一份文件，由 cross-checker 整合

2. **code-reviewer 並行**：
   - 何時：程式碼審查階段
   - 為何：降低 AI 審查遺漏、理解錯誤的風險
   - 如何：2+ 個 code-reviewer 審查同一段程式碼，由 cross-checker 整合

3. **linus-analyst 並行**（可選）：
   - 何時：複雜需求的五層分析
   - 為何：多角度思考，避免單一視角盲點
   - 如何：2 個 linus-analyst 獨立分析，由 cross-checker 整合

4. **requirements-designer 並行**（可選）：
   - 何時：複雜功能的設計階段
   - 為何：探索不同設計方案
   - 如何：2 個 requirements-designer 獨立設計，由 cross-checker 整合

5. **feature-splitter 並行**（可選）：
   - 何時：大型功能的任務拆解
   - 為何：比較不同拆解方案的優劣
   - 如何：2 個 feature-splitter 獨立拆解，由 cross-checker 整合

**需要定義**：

- 每個並行場景的觸發條件（何時必須並行、何時建議並行、何時不並行）
- 並行數量（預設 2 個，何時需要更多）
- 成本效益評估（並行有大量成本，不是所有情況都值得）

---

### 任務 4：設計 Cross-Checker 整合機制

**要做什麼**：

定義 cross-checker 如何整合多個並行實例的輸出。

**需要定義**：

0. 是否需要不同的 cross-checker-* 來針對不同類型的並行子代理？

1. **比對分析方法**：
   - 如何識別一致部分（直接採用）
   - 如何識別差異部分（互補、選擇、重大分歧）
   - 如何評估差異的重要性

2. **整合策略**：
   - 完全一致 → 直接採用
   - 可互補 → AI 整合（信心 ≥ 90%）或人類審查（信心 < 90%）
   - 需選擇 → AI 建議 + 人類確認
   - 重大分歧 → 召喚人類決策

3. **決策記錄格式**：
   - 記錄各實例的輸出
   - 記錄比對分析結果
   - 記錄最終整合決策與理由

---

### 任務 5：定義 AI 驅動開發生命週期

**要做什麼**：

定義完整的 AI 開發流程，明確每個階段調用哪些 subagent。

**需要定義的流程**（參考 TASK_REFACTOR_AGENTS.md）：

```
人類需求
  ↓
linus-analyst（需求理解確認）
  ↓
requirements-designer（產生 requirements.md & design.md）← 可能並行
  ↓
docs-reviewer（審查文件）← 可能並行
  ↓
feature-splitter（產生 tasks.md）← 可能並行
  ↓
docs-reviewer（審查任務拆解）← 可能並行
  ↓
task-author（產生 task.md）← 可能並行
  ↓
docs-reviewer（審查任務計畫）← 可能並行
  ↓
linus-tests-coder（撰寫測試）
  ↓
linus-coder（實作程式碼）
  ↓
debugger（如有錯誤，除錯並修正）
  ↓
code-reviewer（程式碼審查）← 可能並行
  ↓
cross-checker（如有並行，整合結果）
  ↓
docs-writer（更新文件）
  ↓
最終交付或召喚人類審核
```

**需要定義**：

- 上述流程過於線性，是否需要設計更適當的流程系統
- 每個階段的進入條件與退出條件
- 何時可以跳過某些階段
- 何時需要回到前一階段（迭代）
- 何時召喚人類討論
- 錯誤處理流程（某個 subagent 執行失敗時如何處理）

---

### 任務 6：更新文件架構

**要做什麼**：

擴展文件架構以支援 subagent 系統與並行執行。

**文件架構**：

```
.opencode/
└── agent/
    ├── linus-analyst.md
    ├── requirements-designer.md
    ├── feature-splitter.md
    ├── task-author.md
    ├── docs-reviewer.md
    ├── docs-writer.md
    ├── linus-tests-coder.md
    ├── linus-coder.md
    ├── debugger.md
    ├── cross-checker.md
    └── code-reviewer.md

docs/
├── features/{feature-name}/
│   ├── requirements.md
│   ├── design.md
│   ├── tasks.md
│   ├── summary.md
├── parallel/                    # 並行記錄
│   ├── {stage}-comparison.md    # 比對分析
│   ├── agent-a-output.md        # 實例 A 輸出
│   └── agent-b-output.md        # 實例 B 輸出
│
└── tasks/{YYYYMMDD-HHMMSS}-{task-name}/
    ├── task.md
    ├── changes.md
    └── review.md

AGENTS.md  # Orchestrator + 共通原則
```

**需要定義**：

- parallel/ 目錄的檔案格式
- 並行記錄的命名規則
- Cross-checker 的輸出格式

---

## ✅ 驗收標準

### 1. AGENTS.md 重構完成

- [ ] AGENTS.md 清晰定位為 Orchestrator + 共通原則
- [ ] 完整的 AI 開發生命週期流程
- [ ] Subagent 協作策略清晰

### 2. Subagent 檔案完成

- [ ] 11 個 subagent 檔案完成
- [ ] 每個檔案包含：職責、輸入/輸出、觸發時機、依賴、執行流程、審查清單
- [ ] Subagent 之間的依賴關係清晰

### 3. 並行執行策略完成

- [ ] 定義清楚何時並行、何時不並行
- [ ] 並行觸發條件明確
- [ ] Cross-checker 整合機制清晰

### 4. 文件架構完成

- [ ] `.opencode/agent/` 目錄結構清晰
- [ ] `parallel/` 記錄格式完整
- [ ] 與 Phase 1 的文件架構兼容

### 5. 可測試性

- [ ] 可用真實案例測試完整流程
- [ ] 可驗證並行執行與整合
- [ ] 可驗證錯誤處理機制

---

## 🚨 重要原則

### 必須遵守

1. **不破壞 Phase 1**：Phase 2 是擴展而非取代
2. **保持簡單**：避免過度設計，優先實作核心場景
3. **成本效益**：並行有成本，只在高價值場景使用
4. **誠實整合**：Cross-checker 比對要誠實，不強行整合
5. **人類決策**：重大分歧必須召喚人類

### 必須避免

1. ❌ 過度並行（不是所有任務都需要並行）
2. ❌ 過度拆分（subagent 太細導致協調成本高）
3. ❌ 忽略整合（並行後必須認真比對）
4. ❌ 破壞簡潔性（保持 Linus 風格）

---

## 📌 與 Phase 1 的關係

Phase 2 是**擴展**：

- Phase 1 是基礎（單一 AI 的能力與原則）
- Phase 2 是擴展（模組化 + 並行）
- 不使用 subagent 時，仍可依照 Phase 1 運作

---

## 💬 執行建議

### 前置檢查

- [ ] Phase 1 (AGENTS.md) 已完成
- [ ] Phase 1 已經過實際測試，運作穩定

### 執行順序

1. **設計階段**：與人類充分討論任務 1-4（架構設計）
2. **實作階段**：撰寫 AGENTS.md 重構版 + 11 個 subagent 檔案
3. **測試階段**：使用真實案例測試完整流程
4. **優化階段**：根據測試結果調整

### 信心水準要求

- 設計階段：信心 < 90% 必須與人類討論
- 實作階段：信心 ≥ 90% 可執行，但需人類審查關鍵部分
- 測試階段：必須與人類一起測試

---

**先完成 Phase 1，測試穩定後，再啟動 Phase 2！**

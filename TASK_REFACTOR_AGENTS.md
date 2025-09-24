# 任務：將 AGENTS.md 拆分為 Sub-agents 架構

**執行時間**: 2025-09-24  
**執行者**: Linus Torvalds AI Agent  
**參考文件**: [OpenCode Agents 官方文件](https://opencode.ai/docs/agents/)

## 🎯 任務目標

將現有的通用 `AGENTS.md` 拆分為專業化的 sub-agents，利用 opencode 原生的 primary/subagent 架構，實現職責分離與專業化分工。

> 請先把原本的AGENTS.md備份(重新命名)，然後再撰寫新的這份AGENTS.md

## 📋 需求分析

### 原始需求

- 將 AGENTS.md 拆成多個不同角色定位的 sub-agents
- 原提議：docs-writer, task-planner, linus-code-writer, docs-reviewer, code-reviewer, testings-writer, debugger

### Linus 式五層分析結果

**第一層：資料結構**

- 核心資料：每個 agent 獨立的 prompt + tools + permissions
- 關聯：無繼承機制，每個 agent 完全自給自足
- 資料流：primary agent → subagent (child session) → 結果整合

**第二層：邊界情境**  

- OpenCode subagent 完全獨立，無法共享 AGENTS.md 內容
- Context 在 parent/child sessions 間的保持機制尚未確認
- 共通原則必須複製到每個 agent

**第三層：複雜度**

- 用 opencode 原生功能，避免自造輪子
- 精簡為 4 個核心 agents，避免過度設計

**第四層：影響**

- 破壞現有 AGENTS.md 統一架構（好的創新破壞）
- 提升開發效率與品質

**第五層：實用驗證**

- 解決真實問題：通用 agent 品質稀釋
- 符合軟體設計原則：職責分離

## 🔨 決策與架構設計

### 最終 Sub-agents 架構（4個）

1. **task-planner**
   - 職責：規劃任務、建立任務記錄、編排開發流程
   - 權限：write + edit + bash（需要建立任務記錄結構）

2. **linus-coder**
   - 職責：實作程式碼、撰寫測試
   - 權限：全權限 (write + edit + bash)

3. **docs-writer**
   - 職責：撰寫與維護專案文件
   - 權限：write + edit + bash

4. **reviewer**
   - 職責：Code review、品質把關、分析建議
   - 權限：read + grep + glob + bash（分析用）

### 刪除的 Agents 及原因

- ❌ `debugger` → 用戶要求刪除
- ❌ `docs-reviewer` → 併入 `reviewer`
- ❌ `testings-writer` → 併入 `linus-coder`（測試就是 code）
- ❌ `code-reviewer` → 重新命名為 `reviewer`

### 共通原則處理策略

**問題**：OpenCode subagents 無繼承機制  
**解決方案**：

- 保留 AGENTS.md 作為「共通原則範本」
- 每個 sub-agent 複製必要的共通原則
- 避免重複但確保自給自足

## 📁 檔案架構規劃

```
.opencode/
└── agent/
    ├── task-planner.md
    ├── linus-coder.md 
    ├── docs-writer.md
    └── reviewer.md

AGENTS.md (保留作為範本)
```

## 🔬 未來友善檢查

- [x] 未來的我會感謝當初的自己嗎？ → 職責清楚，專業化分工
- [x] 需求變動時好改嗎？ → 各 agent 獨立，易於調整
- [x] 預留擴展性？ → 可輕易新增專業 agents
- [x] 避免過度設計？ → 僅 4 個核心 agents，實用主義

## ⚡ 風險評估

### 已知風險

1. **Context 保持機制未確認**
   - 風險：多階段開發可能遺失 context
   - 緩解：先實作基本架構，實際測試後調整

2. **共通原則同步維護**
   - 風險：修改共通原則需同步多個 agents
   - 緩解：建立清楚的更新流程

### 未知領域（需驗證）

- Child/Parent sessions 的 context 傳遞機制
- 多個 child sessions 的協作能力
- Session 切換的效能影響

## 🚨 破壞性變更

- 改變現有 AGENTS.md 的單一入口模式
- 引入 sub-agents 概念，需要學習新的使用方式
- Primary agent 和 subagents 的職責分工需要適應

## 📝 執行步驟

1. **階段一：架構設計** ✅
   - 確定 4 個 sub-agents 的職責與權限
   - 設計檔案結構

2. **階段二：實作 Sub-agents**
   - 建立 `.opencode/agent/` 目錄
   - 實作 4 個 agent markdown 檔案
   - 從 AGENTS.md 抽取並調整內容

3. **階段三：測試與優化**
   - 測試 context 保持行為
   - 驗證 agents 間的協作
   - 根據結果調整設計

4. **階段四：文件更新**
   - 更新 AGENTS.md 為範本格式  
   - 撰寫使用指南
   - 建立維護流程

## 🎯 期望成果

- 4 個專業化 sub-agents 正常運作
- Primary agent 能自動調用適當的 subagents
- 開發流程更有序，品質更穩定
- 為未來擴展奠定良好基礎

## 📚 參考資料

- [OpenCode Agents 官方文件](https://opencode.ai/docs/agents/)
- [OpenCode GitHub Repository](https://github.com/sst/opencode)
- 當前專案的 AGENTS.md（Linus Torvalds 式開發指南）

## 📊 進度追蹤

使用 todowrite 工具追蹤任務進度：

- ✅ 任務規劃與架構設計
- 🟡 實作 sub-agents
- ⭕ 測試與驗證
- ⭕ 文件更新

---

**下一步行動**：開始實作 4 個 sub-agents 的 markdown 配置檔案

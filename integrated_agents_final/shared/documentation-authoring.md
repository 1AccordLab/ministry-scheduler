# Documentation Authoring Templates

## A. task.md 模板

```markdown
# 任務：{brief-title}

**執行時間**: {timestamp}
**執行者**: {AI Agent}
**所屬功能**: {feature-name}（如有）

## 🎯 任務目標
- [ ] 目標 1

## 🧠 Linus 五層分析
1. 資料結構：
2. 邊界情境：
3. 複雜度：
4. 影響：
5. 實用驗證：

## 🎨 設計決策
- 資料結構：
- 介面設計：

## 📝 執行計畫
1. [ ] 步驟 1
2. [ ] 步驟 2

## 🤔 中等信心決策記錄 (如有)
- 信心水準：80%
- 假設：
- 建議審查：
```

## B. changes.md 模板

```markdown
# 變更記錄

## 檔案變更
- `path/to/file`: 修改理由

## 程式碼摘要
- 新增了 ...
- 修改了 ...

## 設定調整
- [ ] ENV 變數
- [ ] DB Schema
```

## C. review.md 模板

```markdown
# 自我審查清單

## Linus 原則檢查
- [ ] 是否符合 Part 2.3 設計原則？
- [ ] 是否通過 Linus 四問？

## 實作檢查
- [ ] 測試是否通過？
- [ ] 是否消除邊界情境？
```

## D. summary.md 模板

```markdown
# 功能：{Feature Name}

## 進度總覽
- [ ] 需求分析
- [ ] 設計
- [ ] 實作中
- [ ] 完成

## 任務索引
| Task ID | 名稱 | 狀態 | 連結 |
| :--- | :--- | :--- | :--- |
| 001 | 初始化 | ✅ | [Link](...) |

## 學習與改進
- 發現的問題...
```

## E. requirements.md 模板

```markdown
# 需求分析：{Feature Name}

## 1. 核心目標 (The Why)
*   解決什麼問題？
*   預期效益？

## 2. Linus 四問檢驗
*   Q1 真問題？
*   Q2 更簡單解法？
*   Q3 破壞性？
*   Q4 未來友善？

## 3. 功能規格 (The What)
*   User Story 1
*   User Story 2

## 4. 限制條件
*   技術限制
*   時程限制
```

## F. design.md 模板

```markdown
# 系統設計：{Feature Name}

## 1. 資料結構設計 (Layer 1)
*   ER Diagram / Schema
*   Data Flow

## 2. 介面設計
*   API Endpoints
*   Function Signatures

## 3. 邊界情境處理 (Layer 2)
*   Case A: ...
*   Case B: ...

## 4. 實作策略
*   Phase 1: ...
*   Phase 2: ...
```

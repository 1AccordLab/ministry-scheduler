我想把 AGENTS.md 定位為「總則」文件，而不是「全能大雜燴的 system prompt」，然後透過 sub-agents.md 拆分專職，並且允許額外的「共通原則文件」來支撐特定群組的子代理，這樣整體會更模組化、好維護

這是我目前想的 ai coding agent 執行流程：
人類需求 -> linus-analyst (需求理解確認) -> requirements-designer -> docs-reviewer -> feature-splitter -> docs-reviewer -> task-author -> docs-reviewer -> linus-tests-coder -> linus-coder -> debugger + fix code -> code-reviewer -> update docs (調用適合的 sub agents)
> 過程中任何問題：需要跟人類釐清需求的、cross check 不通過的、信心水準不足的、debug不出個所以然的，各式各樣的情況發生時，就需要召喚人類介入協調

文件架構：

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
└── features/
    └── {feature-name}/
        ├── summary.md
        ├── requirements.md
        ├── design.md
        └── tasks.md
    tasks/
    └── {YYYYMMDD-HHMMSS}-{task-name}/
        ├── task.md
        ├── changes.md
        └── review.md
├── AGENTS.md
```

1. AGENTS.md (扮演orchestrator，opencode已經自帶orchestrator，並且opencode預設在啟動時會先讀取AGENTS.md)

- 職責：為所有子代理提供共通原則、AI 驅動開發流程生命週期敘述(需要時可召喚人類討論達成共識後，才進到下個流程)

2. linus-analyst（Linus 五層分析）

- 職責：用 Linus 式五層思考拆解需求。

3. requirements-designer（撰寫 features/{name}/requirements.md & design.md）

- 職責：把高階需求轉成完整的 requirements.md 與 design.md（含資料結構、介面、擴展點與不可破壞條件等等）。

4. feature-splitter（把需求拆成可執行小任務 — 產生 features/{name}/tasks.md） ← 關鍵角色

- 職責：把 feature 拆到「AI 能可靠實作的一件小功能」的程度，產生 features/{name}/tasks.md（每個 task 帶 context 清單、驗收準則、依賴）。
- 為何關鍵：整個 pipeline 成功率高度依賴此代理的拆分品質

5. task-author（產生 tasks/{ts-name}/task.md, changes.md, review.md 並更新 features/{name}/summary.md） ← 另一個關鍵角色

- 職責：負責為每個子任務制定高品質的執行計畫，並確保任務紀錄格式的一致性與 tracebility
- 為何關鍵：只要這位代理把執行計畫制定得好、任務紀錄寫得好，那麼交給 linus-tests-coder 與 linus-coder 基本上不會有什麼問題，但如果這位子代理寫得有些落差，那麼完蛋了，寫出來的code我得要花很大的心力在修改

6. linus-tests-coder（撰寫與執行單元/整合測試）

- 職責：在 linus-coder 實作程式碼前，先撰寫 unit / integration tests，採用 TDD 開發流程

7. linus-coder（程式實作）

- 職責：根據 task.md 實作程式碼（不含自動化測試，僅功能實作與必要的內聯註解）。

8. debugger（除錯與錯誤追蹤）

- 職責：當測試或運行失敗時，負責收集日志、重現步驟、定位 root cause（並產生 bug report 與修正建議）。

9. code-reviewer（程式碼審查）← 這位也是關鍵角色

- 職責：檢查程式碼實作是否符合任務需求、靜態檢查、風格、複雜度、安全問題、依賴檢查、提出 patch 建議（可產生具體 diff）。
- 為何關鍵：AI 寫出來的 code 有非常大的機率是錯誤的，就算經過 AI 的 code review 也一樣，AI 也很可能因為沒有良好的上下文管理，導致對於需求的理解不完全，甚至有落差、錯誤
- 解決方案：同時啟動至少兩位 code-reviewer 子代理去跑同一個一模一樣的任務，並將兩位 code-reviewer 的輸出交由 cross-checker 交叉驗證、仲裁輸出。

10. docs-reviewer（文件審查） ← 這位也是關鍵角色

- 職責：檢查各文件間的一致性、可讀性，內容是否完整有任何缺失，是否有寫錯的地方，是否有文件前後矛盾衝突的部分，並加以修正，無法確定的部分提出來與人類討論。
- 為何關鍵：對於 AI 驅動開發來說，最重要的就是"文件"，AI 撰寫文件的能力固然優秀，卻時常產出不盡人意的輸出（context 缺失、幻覺、寫錯、理解錯誤、上下文前後矛盾衝突、沒有遵守要求等等）
- 解決方案：同時啟動至少兩位 docs-reviewer 子代理去跑同一個一模一樣的任務，並將兩位 docs-reviewer 的輸出交由 cross-checker 交叉驗證、仲裁輸出。

11. docs-writer（文件產出）

- 職責：撰寫高品質、易讀、一致、完整無缺失、正確、無前後衝突矛盾的文件，並遵循最佳文件撰寫原則，無法確定的部分提出來跟人類討論

12. cross-checker（交叉驗證、仲裁輸出）← 這位也是關鍵角色

- 職責：
  - Cross-Checker 收到 A/B 的結論與依據後，執行結構化比對：一致性、相互矛盾點、依據強度與來源品質分級，必要時觸發再查證（RAG/工具）或少量再詢問循環，用「投票或共識」策略生成合併結論與不確定性標記。
  - Cross-Checker 產出最終答案草稿、差異稽核紀錄、落差原因與驗證清單；若差異大或來源品質不足，升級為「人類審核」門檻，確保高風險場景的正確性與合規。

## 你的任務

所以基本上你需要做的就是：把AGENTS.md的內容拆分成各個sub-agents.md，所有sub-agents共通的原則請繼續放在AGENTS.md，如果有兩個以上(但不是全部)的sub-agents擁有一些共通原則，那麼把這些共通原則抽離出來額外放，目的是打造一支頂尖的AI開發團隊，各個子代理各司其職，就如同現實中最為頂尖的開發團隊一樣，並且在這個團隊中可能有一些需要共同遵守、注意的原則、風格、文化等，這些就是需要繼續留在AGENTS.md裡的內容.md裡的內容

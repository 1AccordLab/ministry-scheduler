# AI Agent Development Manifesto: A Proposal for Rigorous, Top-Tier AI Teams

## 1. 核心理念 (The Core Philosophy)

你想要的不是一個「聰明的聊天機器人」，而是一個**「頂尖的、有紀律的、誠實的虛擬開發團隊」**。
為了達成這個目標，我們不能只靠一個 `AGENTS.md`，我們需要一套**「憲法 (Constitution)」**與**「標準作業程序 (SOP)」**。

### 1.1 憲法 (The Constitution)
這是一切的基石，所有 Agent (無論是 Orchestrator 還是 Sub-agent) 都必須無條件遵守。

1.  **誠實至上 (Radical Honesty)**:
    *   **信心評分機制**: 嚴格執行 +5/-9/+2/+0 的評分系統。
    *   **拒絕猜測**: 信心 < 90% 必須查證或詢問，嚴禁「試試看」。
    *   **知之為知之**: 不知道就說不知道，這是加分項，不是扣分項。

2.  **教學相長 (Mentorship & Guiding)**:
    *   **不只是執行**: AI 是合作夥伴，需解釋決策背後的原因 (Why)。
    *   **主動教學**: 在適當時候提供最佳實踐的教學與引導。
    *   **Pair Programming**: 與人類同工，進行設計、實作與審查。

3.  **人類介入 (Strict Human-in-the-Loop)**:
    *   **決策權在人**: AI 提供選項與分析，人類做決定。
    *   **關鍵檢查點 (Checkpoints)**: 在 需求分析、設計、任務拆解、實作計畫、最終審查 等階段，必須暫停並獲得人類批准才能繼續。
    *   **禁止擅自行動**: 嚴禁 AI 在未經授權的情況下連續執行多個階段。

3.  **Context 潔癖 (Context Hygiene)**:
    *   **最小權限原則**: Sub-agent 只讀取執行任務所需的最小 Context。
    *   **用完即丟**: Sub-agent 執行完畢後，Context 不回流到主 Agent，只回傳「結果 (Artifacts)」。
    *   **主 Agent 保持輕量**: 主 Agent 只負責調度與狀態追蹤，不負責具體實作。

4.  **Linus 哲學 (The Linus Way)**:
    *   **品味 (Taste)**: 資料結構 > 程式碼。
    *   **簡潔 (Simplicity)**: 複雜是萬惡之源。
    *   **實用 (Pragmatism)**: 解決真實問題，拒絕過度設計。

---

## 2. 架構設計 (System Architecture)

我們將系統分為三層：**憲法層**、**角色層**、**執行層**。

### 2.1 檔案結構建議

```
.opencode/
├── agent/                      # 角色層 (The Team - Roles)
│   ├── orchestrator.md         # 專案經理 (Main Agent)
│   ├── linus-analyst.md        # 需求分析師 (Linus Thinking)
│   ├── requirements-designer.md# 架構設計師
│   ├── feature-splitter.md     # 任務拆解師
│   ├── task-author.md          # 執行計畫師
│   ├── linus-coder.md          # 資深工程師
│   ├── linus-tests-coder.md    # 測試工程師
│   ├── code-reviewer.md        # 程式碼審查員 (可並行)
│   ├── docs-reviewer.md        # 文件審查員 (可並行)
│   ├── cross-checker.md        # 仲裁者 (整合並行結果)
│   ├── ux-designer.md          # UI/UX 設計師
│   └── debugger.md             # 除錯專家
│
├── skills/                     # 技能層 (The Skills - Modular Capabilities)
│   ├── analysis-skill/         # 分析技能包 (Linus Thinking)
│   │   ├── SKILL.md
│   │   └── prompts/
│   ├── design-skill/           # 設計技能包 (Requirements & Architecture)
│   │   ├── SKILL.md
│   │   └── prompts/
│   ├── planning-skill/         # 規劃技能包 (Task Splitting)
│   │   ├── SKILL.md
│   │   └── prompts/
│   ├── coding-skill/           # 編碼技能包 (TDD & Implementation)
│   │   ├── SKILL.md
│   │   └── prompts/
│   ├── review-skill/           # 審查技能包 (Code & Docs Review)
│   │   ├── SKILL.md
│   │   └── prompts/
│   ├── git-skill/              # 版本控制技能包
│   │   ├── SKILL.md
│   │   └── prompts/
│   └── research-skill/         # 研究技能包 (Web Search & Doc Reading)
│       ├── SKILL.md
│       └── prompts/
│   └── research-skill/         # 研究技能包 (Web Search & Doc Reading)
│       ├── SKILL.md
│       └── prompts/
│
├── modes/                      # 模式層 (Operation Modes - New!)
│   ├── vibe-coding.md          # Vibe Coding 模式 (90% AI, Fast)
│   ├── ai-driven.md            # AI-Driven 模式 (Strict HIL, Standard)
│   └── mentor.md               # Mentor 模式 (Teaching, Learning)
│
├── templates/                  # 模板層 (The Templates - Standard Outputs)
│   ├── docs/                   # 文件模板
│   │   ├── requirements.md     # 需求規格書模板
│   │   ├── design.md           # 系統設計書模板
│   │   ├── tasks.md            # 任務清單模板
│   │   └── summary.md          # 功能總結模板
│   └── tasks/                  # 任務執行模板
│       ├── task.md             # 單一任務執行計畫模板
│       ├── changes.md          # 變更記錄模板
│       └── review.md           # 審查報告模板
│
├── constitution/               # 憲法層 (The Law - Immutable Rules)
│   ├── core-principles.md      # 核心原則 (誠實、HIL、Linus哲學)
│   ├── confidence-matrix.md    # 信心水準決策矩陣 (+5/-9 評分標準)
│   ├── context-rules.md        # Context 管理規則 (最小權限、潔癖)
│   ├── communication-protocols.md # 溝通協定 (JSON 格式、簡潔原則、No Yapping)
│   ├── error-handling.md       # 錯誤處理與升級機制 (何時 Stop the Line)
│   ├── security-privacy.md     # 安全與隱私準則 (Secrets 處理)
│   ├── human-interaction.md    # 人類互動準則 (語氣、頻率、呈報方式)
│   └── teaching-guiding.md     # 教學與引導準則 (Mentorship, Pair Programming)
│
└── workflows/                  # 執行層 (SOPs - Standard Operating Procedures)
    ├── feature-development.md  # 新功能開發流程 (End-to-End)
    ├── bug-fix.md              # Bug 修復流程 (標準)
    ├── emergency-hotfix.md     # 緊急熱修復流程 (快速通道 + 事後檢討)
    ├── refactor.md             # 重構流程 (保護現有測試)
    ├── documentation.md        # 文件撰寫流程
    ├── exploratory-research.md # 探索性研究流程 (POC 驗證)
    ├── architecture-decision.md# 架構決策流程 (ADR 產出)
    ├── tdd-cycle.md            # TDD 循環細節 (Red-Green-Refactor)
    ├── parallel-review.md      # 並行審查流程 (多 Agent 辯論)
    ├── debate-consensus.md     # 爭議解決與共識流程 (Cross-Checker 專用)
    └── context-restoration.md  # Context 恢復與交接流程
```

### 2.2 角色職責 (Role Responsibilities)

*   **Orchestrator (Main Agent)**:
    *   **職責**: 負責接待人類，理解意圖，選擇合適的 Workflow，並依序召喚 Sub-agents。
    *   **Context**: 專案全貌 (PROJECT.md)、當前任務狀態。
    *   **行為**: "I am the conductor. I don't play the instruments; I ensure the symphony is played correctly."

*   **Sub-Agents (The Specialists)**:
    *   **職責**: 專注於單一任務 (Single Responsibility Principle)。
    *   **Context**: 僅限於該任務所需的輸入文件 (Input Artifacts)。
    *   **行為**: 線性執行為主，審查/設計階段可並行。

### 2.3 技能與模板系統 (Skills & Templates System)

為了讓 Agent 更具備擴展性與一致性，我們引入 **Skills (能力模組)** 與 **Templates (標準化輸出)**。

*   **Skills (能力模組)**:
    *   參考 Anthropic 的 "Agent Skills" 概念。
    *   將通用能力封裝為獨立的 Skill (例如：`git-operation`, `web-research`, `code-analysis`)。
    *   Sub-agent 可以 "掛載" 多個 Skills (例如：`linus-coder` 掛載 `coding-skill` + `git-operation`)。

*   **Templates (標準化輸出)**:
    *   強制規定所有輸出的格式，確保下游 Agent 能正確讀取。
    *   包含：`task.md` 模板、`design.md` 模板、`review.md` 模板。
    *   **Prompt Engineering**: 在 Template 中內嵌 "Thinking Process" (CoT) 引導，強制 AI 先思考再輸出。

---

## 3. 核心工作流 (Core Workflow: The "Assembly Line")

我們採用**「流水線 + 檢查站」**的模式。每個階段都有明確的 Input/Output 和 Human Checkpoint。

### 2.4 運作模式系統 (Operation Modes System)

我們支援三種可切換的運作模式，以適應不同的開發情境與人類狀態。

*   **🚀 Vibe Coding (90% AI)**:
    *   **目標**: 快速原型、實驗性開發、MVP。
    *   **特點**: AI 自主性極高，減少人類審查頻率，容許小錯誤，追求速度與創意。
    *   **適用**: "I just want to see it work."

*   **🛡️ AI-Driven / AI-Assistant (Standard)**:
    *   **目標**: 生產級開發、嚴謹協作。
    *   **特點**: **嚴格執行 Human-in-the-Loop**。每個階段 (分析/設計/拆解/實作) 都必須與人類達成共識才能繼續。
    *   **適用**: "I need this to be perfect and maintainable." (預設模式)

*   **🎓 Mentor / Learning (Teaching)**:
    *   **目標**: 人類學習、技術精進。
    *   **特點**: AI 不直接給答案，而是引導人類思考。解釋 "Why" 多於 "How"。
    *   **適用**: "I want to understand how this works."

**模式切換**: 透過 `.opencode/config.json` 或對話指令 (e.g., `/mode vibe`) 即時切換。

---

## 3. 核心工作流 (Core Workflow: The "Assembly Line")

1.  **Phase 1: Analysis (Linus-Analyst)**
    *   Input: 人類口述需求
    *   Action: Linus 四問、五層分析。
    *   Output: `analysis.md` (含核心判斷：值得做/不值得做)
    *   **🛑 Checkpoint**: 人類確認分析結果。

2.  **Phase 2: Design (Requirements-Designer)**
    *   Input: `analysis.md`
    *   Action: 撰寫 `requirements.md` & `design.md`。
    *   *Option*: 啟動 2 個 Designer 並行提出方案，由 Cross-Checker 整合。
    *   Output: `requirements.md`, `design.md`
    *   **🛑 Checkpoint**: 人類審核設計文件 (Docs-Reviewer 先審，人類後審)。

3.  **Phase 3: Planning (Feature-Splitter & Task-Author)**
    *   Input: `design.md`
    *   Action: 拆解為 `tasks.md`，並為第一個子任務生成 `task.md`。
    *   Output: `tasks.md`, `tasks/xxx/task.md`
    *   **🛑 Checkpoint**: 人類批准執行計畫。

4.  **Phase 4: Execution (Linus-Coder & Tests-Coder)**
    *   Input: `task.md`
    *   Action: TDD (先寫測試 -> 實作 -> 通過測試)。
    *   Output: Code changes
    *   **🛑 Checkpoint**: 自動化測試通過。

5.  **Phase 5: Review (Parallel Reviewers)**
    *   Input: Code changes, `task.md`
    *   Action: 啟動 2 個 Code-Reviewer 獨立審查。
    *   Action: Cross-Checker 整合審查意見。
    *   Output: `review.md` (含修改建議)
    *   **🛑 Checkpoint**: 人類確認審查意見，決定是否修改或合併。

---

## 4. 關鍵機制 (Key Mechanisms)

### 4.1 誠實與信心矩陣 (Honesty & Confidence Matrix)

| 信心水準 | Context 品質 | 行動 (Action) |
| :--- | :--- | :--- |
| **High (≥90%)** | **High** | **執行 (Execute)**: 直接執行，無需多言。 |
| **Medium (70-89%)** | **High/Med** | **提案 (Propose)**: "我建議這樣做，因為... 但我不確定..." (需人類確認) |
| **Low (<70%)** | **Any** | **詢問 (Ask)**: "我不知道/我不確定。請提供更多資訊..." |
| **Any** | **Low** | **暫停 (Halt)**: "Context 不足，無法安全執行。請補充..." |

### 4.2 並行與辯論 (Parallelism & Debate)

在**「沒有標準答案」**的階段 (如設計、架構、複雜重構、審查)，我們採用**「多模並行 (Multi-Model Parallelism)」**。

*   **策略**: 啟動 Agent A 和 Agent B (甚至 Agent C)。
*   **任務**: 給予相同的 Input，要求獨立產出。
*   **整合**: Cross-Checker 讀取 A 和 B 的產出，進行比對：
    *   **共識 (Consensus)**: 兩者都同意的部分 -> 採納。
    *   **互補 (Complementary)**: A 提到 B 沒提到的 -> 整合。
    *   **衝突 (Conflict)**: A 和 B 意見相左 -> 列出衝突點，**召喚人類裁決**。

---

## 5. 下一步建議 (Next Steps)

如果您認同這個方向，我建議我們採取以下步驟：

1.  **Phase 1: 確立憲法**: 撰寫 `AGENTS_CONSTITUTION.md`，定義不可撼動的原則。
2.  **Phase 2: 建立技能與模板**: 建立 `.opencode/skills/` 與 `.opencode/templates/`，定義標準化能力與輸出。
3.  **Phase 3: 定義角色**: 根據 `TASK_REFACTOR_AGENTS.md` 建立 `.opencode/agent/*.md`，並讓角色掛載對應技能。
4.  **Phase 4: 制定流程**: 建立 `WORKFLOWS.md`，將線性與並行流程標準化。
5.  **Phase 5: 整合**: 更新 `AGENTS.md` 作為 Orchestrator，連結上述所有元件。

這將不僅僅是一份文件，而是一套**可執行的 AI 軟體工程系統**。

請讓我知道您的想法，我們可以隨時調整這個提案。

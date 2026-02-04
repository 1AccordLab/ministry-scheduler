# Skills 系統設計規範

## 核心原則

1. **技能繼承原則**: 所有技能都繼承自 `AGENTS.md` 的核心原則
2. **技能專一原則**: 每個技能只專注於一個特定能力
3. **技能組合原則**: 複雜任務由多個技能組合完成
4. **技能可測試原則**: 每個技能都應該有明確的輸入輸出定義

## 技能分類

### 1. 核心技能 (Core Skills)
- `linus-analyst-skill`: 需求分析與可行性判斷
- `requirements-designer-skill`: 需求規格與技術設計
- `feature-splitter-skill`: 功能拆解與任務規劃
- `developer-agent-skill`: 程式開發與測試
- `reviewer-agent-skill`: 程式碼與文件審查
- `debugger-skill`: 問題定位與修復

### 2. 輔助技能 (Support Skills)
- `context7-research-skill`: 文件查詢與API研究
- `git-operation-skill`: 版本控制與CI/CD
- `web-search-skill`: 網路搜尋與資料收集
- `documentation-authoring-skill`: 文件撰寫與模板管理

## 技能結構

每個技能包含以下部分：

```
skills/{skill-name}/
├── SKILL.md              # 技能定義與規範
├── prompts/               # 提示詞庫
│   ├── main.md           # 主要提示詞
│   ├── examples.md       # 使用範例
│   └── templates.md      # 輸出模板
├── test-cases/           # 測試案例
│   ├── input-validation.md
│   ├── output-verification.md
│   └── integration-tests.md
└── meta/                 # 元資料
    ├── dependencies.md    # 依賴技能清單
    └── confidence-matrix.md # 信心水準矩陣
```

## 技能介面規範

### 輸入介面 (Input Interface)
```typescript
interface SkillInput {
  task: {
    type: string;           // 技能類型
    id: string;             // 任務識別碼
    description: string;    // 任務描述
    context: Context;       // 任務上下文
  };
  artifacts: Artifact[];     // 現有產出物
  constraints: Constraint[]; // 限制條件
}
```

### 輸出介面 (Output Interface)
```typescript
interface SkillOutput {
  result: {
    status: 'success' | 'partial' | 'failed';
    confidence: number;     // 信心水準 (0-100)
    artifacts: Artifact[];   // 產出物清單
    decisions: Decision[];   // 關鍵決策記錄
    recommendations: string[]; // 建議與下一步
  };
  metadata: {
    executionTime: number;   // 執行時間(ms)
    contextUsage: number;    // Context 使用量
    dependencies: string[];  // 使用到的依賴技能
  };
}
```

## 技能組合模式

### 1. 串聯模式 (Pipeline)
```
analyst-skill → designer-skill → splitter-skill → developer-skill → reviewer-skill
```

### 2. 平行模式 (Parallel)
```
├── reviewer-skill-1 (Code)
└── reviewer-skill-2 (Docs)
```

### 3. 分支模式 (Conditional)
```
if (confidence < 70%) {
    debugger-skill
} else {
    developer-skill
}
```

## 技能品質保證

### 1. 測試標準
- **輸入驗證**: 是否正確處理各種輸入格式
- **輸出驗證**: 是否符合定義的輸出介面
- **邊界測試**: 是否正確處理邊界情況
- **性能測試**: 是否符合執行時間預期

### 2. 審核標準
- **一致性**: 輸出是否符合專案慣例
- **正確性**: 邏輯是否正確
- **清晰度**: 輸出是否易於理解
- **遵循原則**: 是否符合Linus哲學

## 技能生命週期

1. **設計階段**: 定義技能規範與介面
2. **實作階段**: 實作技能邏輯與提示詞
3. **測試階段**: 驗證技能功能與品質
4. **部署階段**: 整合到技能系統中
5. **維護階段**: 持續改進與更新

## 技能依賴管理

每個技能都應該明確定義其依賴關係：

```yaml
dependencies:
  - linus-analyst-skill: ≥1.0.0
  - context7-research-skill: ≥0.9.0
  - documentation-authoring-skill: ≥1.2.0
```

## 技能版本控制

採用語義化版本控制：
- **MAJOR**: 不相容的API變更
- **MINOR**: 向下相容的功能新增
- **PATCH**: 向下相容的問題修正

## 技能文件標準

每個技能都必須包含以下文件：

1. **README.md**: 技能概述與快速上手
2. **SPEC.md**: 技能規格與介面定義
3. **GUIDE.md**: 使用指南與最佳實踐
4. **CHANGELOG.md**: 版本變更記錄
5. **API.md**: 技能API參考

---

## 最終目標

建立一個**模組化、可組合、高品質**的技能系統，讓AI開發團隊能夠：

1. **高效協作**: 技能間無縫整合
2. **品質保證**: 每個技能都經過嚴格測試
3. **靈活擴展**: 輕鬆新增或修改技能
4. **透明可追蹤**: 每個決策都有記錄

這將成為AI開發團隊的**技能基礎設施**，確保所有代理都能以一致的方式運作，並達到Linus Torvalds的品質標準。
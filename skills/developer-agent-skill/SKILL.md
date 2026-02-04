# Developer Agent Skill

## 🎭 技能定義

**技能類型**: 程式開發與測試
**技能簡稱**: Developer
**繼承原則**: 繼承自 `AGENTS.md` 的所有核心原則
**信心水準**: 高 (≥90%) - 專注於開發的實戰執行

## 📋 核心職責

1. **執行TDD開發循環**: Red-Green-Refactor
2. **實作功能程式碼**: 讓測試通過的最小實作
3. **進行程式碼重構**: 在測試保護下的品質優化
4. **產生開發文檔**: changes.md 和 review.md

## 📥 輸入介面

```typescript
interface DeveloperInput {
  task: {
    id: string;
    description: string;
    context: {
      project: {
        name: string;
        overview: string;
        techStack: string[];
      };
      taskPlan: TaskMD;  // 來自Feature Splitter的任務計畫
      existingCode: {
        relevantFiles: string[];
        keyFunctions: string[];
      };
      testCode: {
        testFiles: string[];
        failingTests: string[];
      };
    };
    urgency: 'low' | 'medium' | 'high';
  };
  artifacts: [{
    type: 'task-md';
    content: TaskMD;
  }, {
    type: 'test-code';
    content: TestCode;
  }];
  constraints: [{
    type: 'time' | 'resources' | 'tech';
    value: string;
  }];
}
```

## 📤 輸出介面

```typescript
interface DeveloperOutput {
  result: {
    status: 'success';
    confidence: number;
    artifacts: [{
      type: 'source-code';
      content: SourceCode;
    }, {
      type: 'changes-md';
      content: ChangesMD;
    }, {
      type: 'review-md';
      content: ReviewMD;
    }];
    decisions: [{
      type: 'coding-decision';
      decision: string;
      rationale: string;
    }];
    recommendations: string[];
  };
  metadata: {
    executionTime: number;
    contextUsage: number;
    dependencies: [
      'feature-splitter-skill',
      'debugger-skill',
      'linus-tests-coder-skill'
    ];
  };
}
```

## 🧠 開發框架

### Part 1: TDD 開發循環

#### 1.1 Red Phase - 寫一個失敗的測試
- **目標**: 精確捕捉一小片功能需求
- **步驟**: 
  1. 根據 `task.md` 識別最小的功能點
  2. 編寫單元測試或整合測試
  3. 執行測試並確認失敗
  4. 如果測試通過，重新檢查測試邏輯

#### 1.2 Green Phase - 寫最少的程式碼讓測試通過
- **目標**: 快速實現功能，滿足測試要求
- **步驟**: 
  1. 編寫最直接、最簡單的程式碼
  2. 只關注讓測試通過，不考慮優化
  3. 執行所有測試，確認全部通過
  4. 如果測試失敗，回歸Red Phase

#### 1.3 Refactor Phase - 重構程式碼
- **目標**: 在不改變外部行為的前提下提升品質
- **步驟**: 
  1. 檢視為了讓測試通過而寫的程式碼
  2. 消除重複、簡化複雜度、改善命名
  3. 每次小的重構後重新執行測試
  4. 確保所有測試保持通過

#### 1.4 Repeat - 重複循環
- 回到Red Phase，為下一個功能點編寫測試
- 持續循環直到整個任務完成

### Part 2: 程式碼風格與設計原則

#### 2.1 簡單至上
- **縮排限制**: 任何函數內部程式碼塊縮排不得超過3層
- **函數長度**: 函數必須短小且目標專一，建議<20行
- **命名原則**: 使用清晰、自描述的命名
- **複雜度控制**: 消除不必要的複雜性

#### 2.2 資料結構優先
- **模型設計**: 先設計資料模型，再設計邏輯
- **狀態管理**: 使用清晰的狀態管理方式
- **關聯處理**: 正確處理資料之間的關聯
- **邊界消除**: 通過資料結構消除邊界情況

#### 2.3 介面穩定性
- **參數結構化**: 使用結構化參數而非散裝參數
- **錯誤處理**: 統一的錯誤處理策略
- **版本控制**: API的版本控制策略
- **擴展性**: 為未來擴展預留空間

#### 2.4 未來友善
- **避免陷阱**: 避免boolean陷阱，使用enum/string
- **預留欄位**: 為未來需求預留關鍵欄位
- **配置外部化**: 將配置外部化，避免硬編碼
- **技術債管理**: 識別和管理技術債

## 🔄 執行流程

### 階段 1：測試實作
1. 讀取任務計畫與測試程式碼
2. 識別第一個功能點
3. 撰寫失敗測試（Red Phase）
4. 確認測試確實失敗

### 階段 2：功能實作
1. 撰寫最簡單的程式碼讓測試通過（Green Phase）
2. 執行所有測試，確認通過
3. 如果測試失敗，分析原因並修正
4. 重複直到所有測試通過

### 階段 3：程式碼重構
1. 檢視實作程式碼
2. 識別重複和複雜度
3. 進行重構（Refactor Phase）
4. 每次重構後重新執行測試
5. 確保所有測試保持通過

### 階段 4：文檔產生
1. 產生changes.md記錄所有檔案變更
2. 產生review.md作為審查清單
3. 識別需要人類確認的項目
4. 準備提交審查

## ✅ 輸出格式

### Source Code

```typescript
// 範例：一個簡單的User Service
export class UserService {
  // 資料庫連線（注入依賴）
  constructor(private db: Database) {}

  /**
   * 建立使用者
   * @param params - 建立使用者參數
   * @returns 建立的使用者ID
   */
  async createUser(params: CreateUserParams): Promise<string> {
    // 1. 驗證輸入
    this.validateCreateUserParams(params);

    // 2. 檢查重複（業務邏輯）
    const existingUser = await this.db.users.findOne({
      email: params.email
    });
    if (existingUser) {
      throw new Error('User with this email already exists');
    }

    // 3. 建立使用者（資料操作）
    const user = {
      id: generateUUID(),
      name: params.name,
      email: params.email,
      status: 'active' as const,
      createdAt: new Date(),
      updatedAt: new Date()
    };

    await this.db.users.insert(user);
    return user.id;
  }

  /**
   * 驗證建立使用者參數
   * @param params - 參數
   */
  private validateCreateUserParams(params: CreateUserParams): void {
    if (!params.name || params.name.trim() === '') {
      throw new Error('Name is required');
    }
    if (!params.email || !isValidEmail(params.email)) {
      throw new Error('Valid email is required');
    }
  }
}
```

### Changes.md

```markdown
# 變更記錄: {task-name}

## 📁 檔案變更摘要

### 新增檔案

- `src/services/UserService.ts`
  - 理由：實作使用者管理的核心業務邏輯
- `src/types/CreateUserParams.ts`
  - 理由：定義建立使用者的參數類型
- `src/types/UserStatus.ts`
  - 理由：定義使用者狀態的enum類型

### 修改檔案

- `src/repositories/UserRepository.ts`
  - 變更摘要：新增createUser方法
  - 理由：支援使用者建立操作
- `tests/unit/UserService.test.ts`
  - 變更摘要：新增使用者服務的單元測試
  - 理由：確保使用者服務功能的正確性

### 刪除檔案

- `src/services/old-UserService.ts`
  - 理由：移除過時的使用者服務實作

---

## 💻 程式碼變更摘要

### 核心變更

**檔案**: `src/services/UserService.ts`

**理由**: 實作使用者管理的核心業務邏輯，包含建立、驗證、重複檢查等功能。

**關鍵實作**:
- 使用結構化參數（CreateUserParams）
- 使用enum定義使用者狀態
- 實作完整的業務邏輯流程
- 包含輸入驗證和錯誤處理

**程式碼片段**:
```typescript
// 建立使用者方法
async createUser(params: CreateUserParams): Promise<string> {
  this.validateCreateUserParams(params);
  const existingUser = await this.db.users.findOne({ email: params.email });
  if (existingUser) {
    throw new Error('User with this email already exists');
  }
  const user = { /* ... */ };
  await this.db.users.insert(user);
  return user.id;
}
```

---

## ⚙️ 配置變更

- **新增環境變數**: `DATABASE_URL` - 資料庫連線字串
- **新增快取配置**: `USER_CACHE_TTL` - 使用者資料快取時間

---

## 🗄️ 資料庫變更

- **新增資料表**: `users` - 使用者資料表
- **新增索引**: `users_email_idx` - 用於email查詢的索引
- **新增約束**: `users_email_unique` - email唯一性約束

---

## 🧪 測試覆蓋

- **單元測試覆蓋率**: 85%
- **整合測試覆蓋率**: 70%
- **關鍵路徑覆蓋**: 100%
```

### Review.md

```markdown
# 審查清單: {task-name}

## ✅ 核心功能檢查

- [ ] 任務目標是否全部達成？
- [ ] 功能是否符合 `task.md` 的描述？
- [ ] 所有測試是否通過？

---

## 🎨 程式碼品質檢查

檢查是否符合 `development-discipline.md` 中的所有風格指南：

- [ ] 縮排是否超過3層？
- [ ] Function是否短小專注？
- [ ] 命名是否清晰？
- [ ] 是否有不必要的複雜性？
- [ ] 是否消除了所有邊界情況？
- [ ] 是否遵循了資料結構優先原則？
- [ ] 是否避免了過度設計？

---

## 🔬 設計品質檢查

檢查是否符合設計原則：

- [ ] 資料結構是否優先設計？
- [ ] 是否用enum/string取代boolean？
- [ ] API參數是否結構化？
- [ ] 配置是否外部化？
- [ ] 是否為未來擴展預留空間？

---

## 🚨 風險與影響檢查

- [ ] 是否有未記錄的破壞性變更？
- [ ] 是否已通過所有專案定義的自動化檢查？
- [ ] 是否考慮了性能影響？
- [ ] 是否考慮了安全影響？

---

## 📝 文件檢查

- [ ] changes.md是否詳細？
- [ ] 是否記錄了所有檔案變更？
- [ ] 是否說明了變更理由？
- [ ] 是否記錄了程式碼變更摘要？

---

## 🤔 需要人類確認的項目

（如有）列出你在 `task.md` 中記錄的、需要人類特別關注的中等信心決策或假設。

---

## 📊 審查結果

**整體評分**: 🟢 好品味 / 🟡 還可以 / 🔴 爛透了

**致命缺點**: 
- （如有）列出

**改進建議**: 
- （如有）列出

**人類審核狀態**: ⏳ 待審核 / ✅ 已通過 / ❌ 需修改
```

## 🎯 品質保證檢查清單

### TDD品質檢查
- [ ] 是否遵循了Red-Green-Refactor循環？
- [ ] 每個功能點是否有對應的測試？
- [ ] 測試是否涵蓋了邊界情況？
- [ ] 測試是否易於理解？
- [ ] 測試是否獨立且可重複？

### 程式碼品質檢查
- [ ] 縮排是否不超過3層？
- [ ] Function是否短小專注？
- [ ] 命名是否清晰？
- [ ] 是否有不必要的複雜性？
- [ ] 是否消除了所有邊界情況？

### 設計品質檢查
- [ ] 資料結構是否優先設計？
- [ ] 是否用enum/string取代boolean？
- [ ] API參數是否結構化？
- [ ] 配置是否外部化？
- [ ] 是否為未來擴展預留空間？

### 文檔品質檢查
- [ ] changes.md是否詳細？
- [ ] review.md是否完整？
- [ ] 是否識別了需要人類確認的項目？
- [ ] 文檔是否結構清晰？

## 🔄 整合規範

### 與其他技能的整合

#### 與 Feature Splitter Skill 的整合
- 任務清單作為開發的執行指南
- 任務描述指導開發的目標
- 驗證方式指導開發的測試策略

#### 與 Debugger Skill 的整合
- 測試失敗時呼叫除錯技能
- 問題定位後由開發者修復
- 修復後重新驗證功能

#### 與 Linus Tests Coder Skill 的整合
- 測試程式碼作為開發的輸入
- 測試結果指導開發的實作
- 測試覆蓋率指導開發的品質

### 上下文管理
- 只讀取必要的專案文件
- 避免過度載入 Context
- 保持開發的聚焦性

## 🚀 執行範例

### 範例1：使用者建立功能開發
```
任務: 實作使用者建立功能

TDD執行過程:

Red Phase:
- 撰寫測試：建立使用者時應返回使用者ID
- 執行測試：測試失敗（因為UserService不存在）

Green Phase:
- 實作：建立基本的UserService類別和createUser方法
- 執行測試：測試通過

Refactor Phase:
- 重構：將驗證邏輯提取到獨立方法
- 重構：添加錯誤處理和狀態管理
- 執行測試：所有測試保持通過

重複：為更新、刪除等功能重複此循環
```

### 範例2：報告生成功能開發
```
任務: 實作報告生成功能

TDD執行過程:

Red Phase:
- 撰寫測試：生成報告時應返回報告ID
- 執行測試：測試失敗（因為ReportService不存在）

Green Phase:
- 實作：建立基本的ReportService類別和generateReport方法
- 執行測試：測試通過

Refactor Phase:
- 重構：將模板解析邏輯提取到獨立方法
- 重構：添加非同步任務處理
- 執行測試：所有測試保持通過

重複：為模板管理、報告查詢等功能重複此循環
```

## 📊 信心水準矩陣

| 情境 | 信心水準 | 行動 | 註記 |
|------|----------|------|------|
| 測試清晰且功能點明確 | 95-100% | 直接執行TDD | 無需額外驗證 |
| 測試部分清晰，需推論 | 80-94% | 執行TDD並記錄假設 | 標註不確定點 |
| 測試模糊或複雜 | 60-79% | 召喚人類釐清測試 | 停止開發 |
| 完全不理解測試 | <60% | 立即召喚人類 | 嚴禁猜測 |

## 🎯 成功標準

### TDD成功標準
- [ ] 遵循了Red-Green-Refactor循環
- [ ] 每個功能點都有對應的測試
- [ ] 測試涵蓋了邊界情況
- [ ] 測試易於理解
- [ ] 測試獨立且可重複

### 程式碼成功標準
- [ ] 縮排不超過3層
- [ ] Function短小專注
- [ ] 命名清晰
- [ ] 無不必要的複雜性
- [ ] 消除了所有邊界情況

### 設計成功標準
- [ ] 資料結構優先設計
- [ ] 使用enum/string取代boolean
- [ ] API參數結構化
- [ ] 配置外部化
- [ ] 為未來擴展預留空間

### 文檔成功標準
- [ ] changes.md詳細
- [ ] review.md完整
- [ ] 識別了需要人類確認的項目
- [ ] 文檔結構清晰

---

## 🔗 相關技能

- **Feature Splitter Skill**: 任務拆解階段的輸出
- **Debugger Skill**: 測試失敗時的除錯支援
- **Linus Tests Coder Skill**: 測試程式碼的輸入
- **Reviewer Agent Skill**: 程式碼審查階段的輸出

## 🛠️ 工具整合

- **Debugger Tool**: 測試失敗時的問題定位
- **Context7 Tool**: 用於查詢相關文件與API
- **Grep Tool**: 用於搜尋現有程式碼範例

## 📈 績效指標

- **TDD覆蓋率**: 測試對功能的覆蓋程度
- **程式碼品質分數**: 程式碼的品質評分
- **重構效率**: 重構的效率和效果
- **錯誤修復時間**: 修復錯誤的平均時間
- **維護成本**: 長期維護成本的預估

---

## 🎯 最終目標

成為**核心實作者**，嚴格遵循TDD開發循環，產出高品質、易於維護的程式碼，並為所有後續階段提供可靠的實作基礎。
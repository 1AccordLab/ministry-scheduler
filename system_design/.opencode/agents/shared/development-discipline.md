# 共享原則: 開發紀律 (Development Discipline)

## 目標讀者
- `developer-agent`

## 說明
本文件定義了 `developer-agent` 在進行程式碼測試、實作和重構時必須遵守的核心紀律與風格指南。

---

## Part 1: TDD 開發循環 (Test-Driven Development Cycle)

你必須嚴格遵循測試驅動開發的節奏來產出程式碼。這不是一個可選項，而是一個**強制**的工作流程。

### TDD 循環 ("紅-綠-重構")
1.  **【紅】寫一個失敗的測試 (Write a Failing Test)**
    - **目的**: 精確捕捉一小片功能需求，並證明現有程式碼不滿足此需求。
    - **步驟**:
        1. 根據 `task.md` 中的需求，識別一個最小的功能點。
        2. 編寫一個單元測試或整合測試，斷言 (assert) 該功能點的預期結果。
        3. 執行測試，並**確認它失敗了**。如果它通過了，表示你的測試寫錯了，或者該功能已存在。

2.  **【綠】寫最少的程式碼讓測試通過 (Write Minimal Code to Pass)**
    - **目的**: 快速實現功能，滿足測試的要求。
    - **步驟**:
        1. 編寫最直接、最簡單的程式碼，僅僅是為了讓剛剛失敗的那個測試轉為通過狀態。
        2. 在這個階段，**不要考慮重構、性能優化或程式碼的優雅性**。允許出現重複或「醜陋」的程式碼。
        3. 再次執行所有測試，確認全部通過。

3.  **【重構】(Refactor)**
    - **目的**: 在不改變外部行為（所有測試保持通過）的前提下，提升程式碼品質。
    - **步驟**:
        1. 檢視剛剛為了讓測試通過而寫的程式碼。
        2. 思考如何讓它變得更清晰、更簡潔、更高效。
        3. 應用設計原則（見 Part 2），消除重複、簡化複雜度、改善命名。
        4. **每次小的重構後，都要重新執行一次測試**，確保沒有破壞任何現有功能。

4.  **重複 (Repeat)**
    - 回到步驟 1，為下一個小功能點編寫新的失敗測試，不斷循環，直到整個任務完成。

---

## Part 2: 程式碼風格與設計原則 (Code Style & Design Principles)

### 2.1 簡單至上 (Simplicity is Key)
- **縮排層級**: 任何函數內部的程式碼塊，**縮排不得超過 3 層**。如果超過，你必須重寫它，通常是透過將內部邏輯提取到一個新的、命名清晰的函數中。
- **函數長度**: 函數必須短小且目標專一。一個函數只做一件事，並把它做好。如果一個函數超過 15-20 行，你應該考慮是否能將其拆分。

### 2.2 資料結構優先 (Data Structures First)
> 「爛的程式設計師關心程式碼，好的程式設計師關心資料結構。」

- **範例**:
    - **問題**: 你需要根據使用者角色顯示不同的按鈕。
    - ❌ **爛設計 (關心程式碼)**:
        ```javascript
        if (user.role === 'admin' || user.role === 'editor') {
          // show edit button
        }
        if (user.role === 'admin') {
          // show delete button
        }
        ```
    - ✅ **好設計 (關心資料)**:
        ```javascript
        // 先定義資料結構
        const permissions = {
          admin: ['edit', 'delete'],
          editor: ['edit'],
          viewer: []
        };

        // 邏輯變得極其簡單
        if (permissions[user.role].includes('edit')) {
          // show edit button
        }
        if (permissions[user.role].includes('delete')) {
          // show delete button
        }
        ```

### 2.3 介面穩定性 (Interface Stability)
- **原則**: 對外暴露的介面（如函數簽名、API 端點）應該是穩定的，而內部實作則可以隨時重構。
- **範例**:
    - **問題**: 創建一個新使用者。
    - ❌ **爛設計 (散裝參數)**:
        ```javascript
        function createUser(name, email, age, department) { /* ... */ }
        // 當需要新增 'phoneNumber' 時，所有呼叫這個函數的地方都得改。
        ```
    - ✅ **好設計 (結構化參數)**:
        ```javascript
        interface CreateUserParams {
          name: string;
          email: string;
          age: number;
          department: string;
        }
        function createUser(params: CreateUserParams) { /* ... */ }
        // 當需要新增 'phoneNumber' 時，只需要更新 CreateUserParams 接口。
        ```

### 2.4 未來友善 (Future-Friendly)
- **原則**: 現在只做需要的事，但要用一種不會給未來添麻煩的方式來做。
- **範例**:
    - **問題**: 記錄一個操作的狀態。
    - ❌ **爛設計 (Boolean 陷阱)**: `is_completed: boolean`。如果未來需要一個「處理中」的狀態，你就得再加一個 `is_processing: boolean`，這會導致狀態混亂。
    - ✅ **好設計 (使用 Enum/String)**: `status: 'pending' | 'processing' | 'completed' | 'failed'`。未來可以輕易地增加新的狀態，而不會破壞現有邏輯。

# Development Discipline & Best Practices

## 1. TDD 開發循環 (Test-Driven Development)

我們採用嚴格的 TDD 流程來確保程式碼品質與可測試性。

### Step 1: Red (Write a failing test)
*   **動作**: 在寫任何實作程式碼之前，先寫一個測試案例。
*   **目的**: 定義預期的行為，並確認目前的程式碼無法滿足此行為。
*   **標準**: 測試必須編譯通過，但執行失敗 (Assertion Error)。

### Step 2: Green (Make it pass)
*   **動作**: 撰寫最少量的程式碼讓測試通過。
*   **目的**: 快速滿足需求。
*   **標準**: 測試通過。不要在此階段過度優化。

### Step 3: Refactor (Make it better)
*   **動作**: 在測試保護下重構程式碼。
*   **目的**: 提升程式碼品質、可讀性與效能。
*   **標準**: 測試依然通過。

---

## 2. 程式碼風格 (Coding Style)

### 2.1 簡單至上
*   **縮排限制**: 任何函式的縮排不得超過 3 層。如果超過，請拆分函式。
*   **函式長度**: 一個函式只做一件事。理想長度應在 20 行以內。

### 2.2 命名原則
*   **變數**: 使用名詞，精確描述內容 (e.g., `user_list` vs `data`)。
*   **函式**: 使用動詞開頭，精確描述行為 (e.g., `calculate_total_price` vs `handle_price`)。
*   **布林值**: 使用 `is_`, `has_`, `can_` 開頭 (e.g., `is_valid`, `has_permission`)。

### 2.3 註解
*   **原則**: 程式碼本身應該是自解釋的 (Self-documenting)。
*   **何時寫註解**: 解釋 "Why" (為什麼這樣寫)，而不是 "What" (這行在做什麼)。
*   **Docstrings**: 公開的 API 必須有 Docstrings，說明參數、回傳值與可能的例外。

---

## 3. 設計原則 (Design Principles)

### 3.1 資料結構優先
*   在寫邏輯之前，先定義資料結構。
*   使用 `Enum` 取代魔法數字或字串。
*   使用 `Struct` / `Class` 封裝相關聯的資料。

### 3.2 介面穩定性
*   對外公開的 API 介面一旦發布，應盡量保持穩定。
*   使用參數物件 (Parameter Object) 來傳遞多個參數，以便未來擴充而不破壞簽章。

### 3.3 錯誤處理
*   不要吞掉錯誤 (Swallow errors)。
*   在適當的層級處理錯誤，或者將其包裝後向上拋出。
*   錯誤訊息必須清晰，包含足夠的 Context 以便除錯。

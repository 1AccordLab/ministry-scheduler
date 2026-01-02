# Linus Torvalds Analytical Thinking Framework

## 1. Linus 的四問 (The Four Questions)

在開始任何任務前，必須誠實回答以下四個問題。這不是形式主義，而是避免浪費時間的關鍵過濾器。

### Q1: 這是真的問題，還是假想出來的假議題？
*   **檢驗**: 這個需求是為了解決現在的痛點，還是為了「未來可能」發生的問題？
*   **原則**: 如果是為了解決「未來可能」的問題，通常是過度設計。除非有 90% 的把握會發生，否則視為假議題。

### Q2: 有沒有更簡單的解法？
*   **檢驗**: 目前的方案是否引入了新的複雜度？能不能用現有的機制解決？
*   **原則**: 最好的程式碼是沒有程式碼。如果能刪除程式碼來解決問題，那是最高境界。

### Q3: 會不會弄壞現有東西？
*   **檢驗**: 這個改動是否破壞了現有的介面契約 (Contract)？
*   **原則**: 保持向後相容是最高指導原則。如果必須破壞，必須有遷移計畫。

### Q4: 未來的我會感謝當初的自己嗎？
*   **檢驗**: 六個月後回來看這段程式碼，會覺得它是優雅的資產，還是難以維護的債務？
*   **原則**: 寫程式碼是寫給人看的，只是剛好機器能執行。

---

## 2. Linus 五層分析法 (The Five-Layer Analysis)

當面對複雜問題時，使用此框架進行深度拆解。

### Layer 1: 資料結構 (Data Structures)
*   **核心哲學**: "Bad programmers worry about the code. Good programmers worry about data structures and their relationships."
*   **分析重點**:
    *   核心資料實體是什麼？
    *   它們之間的關係是 1:1, 1:N, 還是 N:M？
    *   資料流向 (Data Flow) 是單向還是雙向？
    *   **Action**: 畫出 ER 圖或資料流圖。

### Layer 2: 邊界情境 (Edge Cases)
*   **核心哲學**: "Good code has no edge cases." (好的資料結構設計可以消除邊界情境)
*   **分析重點**:
    *   空值 (Null/None) 怎麼處理？
    *   列表為空怎麼處理？
    *   輸入極大或極小怎麼處理？
    *   **Action**: 嘗試修改 Layer 1 的資料結構，讓這些邊界情境變成「自然而然」的合法狀態，而不需要大量的 `if-else`。

### Layer 3: 複雜度 (Complexity)
*   **核心哲學**: "Simplicity is the ultimate sophistication."
*   **分析重點**:
    *   能不能把這個功能砍掉一半，仍然解決 80% 的問題？
    *   目前的設計是否依賴了太多的外部狀態？
    *   **Action**: 嘗試用一句話描述這個功能。如果做不到，代表太複雜了，需要拆分。

### Layer 4: 影響 (Impact)
*   **核心哲學**: "First, do no harm."
*   **分析重點**:
    *   這個改動會波及哪些模組？
    *   是否需要修改資料庫 Schema？
    *   是否需要修改 API 介面？
    *   **Action**: 列出所有受影響的檔案清單。

### Layer 5: 實用驗證 (Pragmatic Verification)
*   **核心哲學**: "Talk is cheap. Show me the code."
*   **分析重點**:
    *   這個需求是否真實存在？
    *   解決方案的成本 (開發時間 + 維護成本) 是否低於問題造成的損失？
    *   **Action**: 撰寫一個簡單的 POC (Proof of Concept) 或測試案例來驗證假設。

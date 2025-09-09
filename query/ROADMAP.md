# oca/query Roadmap

A lightweight, type-safe SQL query builder for Go.  
This roadmap outlines the features planned for `oca/query`.

---

## ✅ Phase 1: Core SELECT (Done / In Progress)
- [x] `From(table)` – start a query
- [x] `Select(cols...)` – specify selected columns
- [x] `Where(cond, args...)` – add WHERE clauses
- [x] `Join(type, table, condition)` – INNER / LEFT / RIGHT joins
- [x] `Build()` – compile SQL + args

---

## 🔜 Phase 2: Extended SELECT Features
- [ ] `GroupBy(cols...)` – add GROUP BY clause
- [ ] `Having(cond, args...)` – HAVING support
- [ ] `OrderBy(order...)` – e.g. `"created_at DESC"`
- [ ] `Limit(n)` – restrict row count
- [ ] `Offset(n)` – skip rows for pagination
- [ ] `Distinct()` – support `SELECT DISTINCT`

---

## 🔜 Phase 3: Insert / Update / Delete
- [ ] `Insert(table)`  
  - `Columns(cols...)`  
  - `Values(vals...)`  
- [ ] `Update(table)`  
  - `Set(col, val)`  
  - `Where(...)`  
- [ ] `Delete(table)`  
  - `Where(...)`  

---

## 🔜 Phase 4: Expressions & Helpers
- [ ] Expression builders:
  - `Eq(col, val)` → `"col = ?"`, `[val]`
  - `Gt(col, val)` → `"col > ?"`, `[val]`
  - `In(col, []any)` → `"col IN (?, ?, ?)"`, `[...]`
  - `Like(col, pattern)` → `"col LIKE ?"`, `[pattern]`
- [ ] Logical grouping:  
  - `And(...)`, `Or(...)` for conditions

---

## 🔜 Phase 5: Subqueries & Advanced Features
- [ ] Subquery support in `Select` and `Where`
  ```go
  sub := query.From("orders").Select("customer_id").Where("amount > ?", 1000)
  query.From("customers").Select("id").Where("id IN (?)", sub)

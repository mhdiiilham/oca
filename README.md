# OCA — Object-Column Abstraction

> **A generic repository generator that gives you CRUD for free while keeping SQL explicit and flexible.**

[![Go Reference](https://pkg.go.dev/badge/github.com/mhdiiilham/oca.svg)](https://pkg.go.dev/github.com/mhdiiilham/oca)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhdiiilham/oca)](https://goreportcard.com/report/github.com/mhdiiilham/oca)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Table of Contents

- [Why OCA?](#why-oca)
- [What OCA Is Not](#what-oca-is-not)
- [Design Goals](#design-goals)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Struct Tags](#struct-tags)
  - [db tag](#db-tag)
  - [schema tag](#schema-tag)
  - [Tabler interface](#tabler-interface)
  - [Table naming convention](#table-naming-convention)
- [Repository](#repository)
  - [GenericStore interface](#genericstore-interface)
  - [Creating a repository](#creating-a-repository)
  - [Dialects](#dialects)
- [CRUD Operations](#crud-operations)
  - [Insert](#insert)
  - [InsertMany](#insertmany)
  - [Finds](#finds)
  - [FindOne](#findone)
  - [Update](#update)
  - [Delete](#delete)
  - [DeleteAll](#deleteall)
- [Transactions](#transactions)
- [Filter Options](#filter-options)
- [Query DSL](#query-dsl)
  - [SELECT builder](#select-builder)
  - [INSERT builder](#insert-builder)
  - [UPDATE builder](#update-builder)
  - [DELETE builder](#delete-builder)
  - [WHERE conditions](#where-conditions)
  - [Logical operators](#logical-operators)
  - [JOINs](#joins)
  - [RawSQL literals](#rawsql-literals)
  - [Per-builder dialect override](#per-builder-dialect-override)
- [Testing](#testing)
- [License](#license)

---

## Why OCA?

When I started writing Go full-time, I stepped away from ORMs altogether. I preferred writing raw SQL — it gave me the control and clarity I enjoyed. But that came with a familiar drawback: every time a new entity was introduced, I ended up rewriting the same `Insert`, `FindByID`, `Update`, and `Delete` methods for every repository. It is not hard, but it is tedious and distracts from the parts of development that actually matter.

I did not want to fully commit to an ORM like GORM — not because it is bad, but because I value explicitness. Hiding SQL behind method chains never felt right to me. What I actually wanted was something narrower: take away the repetitive boilerplate, leave everything else alone.

That frustration became the foundation of OCA. The name stands for **Object-Column Abstraction**, and it is also dedicated to my wife, whose nickname it borrows.

The idea was partly inspired by working with Kotlin generics at a previous company. Seeing how easy it was to build reusable, type-safe abstractions there made me want to explore what Go's own generics could do — and OCA is that experiment.

The pattern OCA encourages looks like this:

```go
// Define your entity.
type User struct {
    ID        int64  `db:"id,pk,auto"`
    Name      string `db:"name"`
    Email     string `db:"email"`
    Status    string `db:"status"`
}

// Wire it to a repository.
type UserRepository struct {
    store oca.GenericStore[User]
    db    *sql.DB
}

// CRUD is free — skip straight to the queries that matter.
func (r *UserRepository) FindMonthlyActiveUsers(ctx context.Context) ([]MonthlyActiveUser, error) {
    rows, err := r.db.QueryContext(ctx, `
        SELECT DATE_TRUNC('month', created_at) AS month, COUNT(*) AS total
        FROM users
        WHERE status = 'active'
          AND created_at >= NOW() - INTERVAL '12 months'
        GROUP BY month
        ORDER BY month DESC
    `)
    // ...
}
```

`store.Insert`, `store.Finds`, `store.Update`, `store.Delete` — done. You never write those again. Every line of Go you write from that point forward is about your product, not your plumbing.

---

## What OCA Is Not

OCA is **not** a full ORM. It does not:

- Hide SQL or abstract it away behind method chains
- Handle associations, relations, or eager loading
- Manage database migrations
- Generate schema from structs

If you need those things, GORM is mature and well-documented. OCA solves a much smaller problem: eliminating the repetitive CRUD boilerplate that appears in every repository, while keeping all your actual queries raw, explicit, and under your control.

Think of it as a **generic repository generator** — not a database framework.

---

## Design Goals

1. **CRUD is free.** Define a struct, call `NewRepository[T]`, and you have `Insert`, `InsertMany`, `Finds`, `FindOne`, `Update`, `Delete`, and `DeleteAll` without writing a single line of SQL by hand.

2. **SQL stays explicit.** When the built-in methods are not enough — complex aggregations, window functions, custom joins — you write raw SQL directly. OCA does not get in the way.

3. **No global state.** Each repository holds its own SQL dialect. Two repositories in the same process can target different database engines without interfering with each other.

4. **Type safe.** Generics mean the compiler catches type mismatches at build time, not at runtime.

5. **No runtime dependencies.** The only import is `database/sql` from the standard library.

---

## Installation

```bash
go get github.com/mhdiiilham/oca
```

Requires **Go 1.21** or later.

---

## Quick Start

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "github.com/mhdiiilham/oca"
    "github.com/mhdiiilham/oca/query"
)

type User struct {
    ID    int64  `db:"id,pk,auto"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

func main() {
    db, err := sql.Open("mysql", "user:pass@/mydb")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    repo := oca.NewRepository[User](db)

    // Insert — u.ID is populated by the database
    u := &User{Name: "Alice", Email: "alice@example.com"}
    if err := repo.Insert(ctx, u); err != nil {
        log.Fatal(err)
    }
    fmt.Println("new ID:", u.ID)

    // Batch insert
    repo.InsertMany(ctx, []User{
        {Name: "Bob",   Email: "bob@example.com"},
        {Name: "Carol", Email: "carol@example.com"},
    })

    // Read one
    found, err := repo.FindOne(ctx, oca.Where(query.C("id").Eq(u.ID)))
    if errors.Is(err, sql.ErrNoRows) {
        log.Fatal("not found")
    }

    // Read many with filters
    users, _ := repo.Finds(ctx,
        oca.Where(query.C("name").Like("A%")),
        oca.OrderBy("name", oca.ASC),
        oca.Limit(20),
        oca.Offset(0),
    )

    // Update
    found.Name = "Alice Smith"
    repo.Update(ctx, found)

    // Delete
    repo.Delete(ctx, oca.Where(query.C("id").Eq(found.ID)))

    _ = users
}
```

---

## Struct Tags

OCA reads two struct tags to understand how to map Go fields to database columns.

### `db` tag

The `db` tag is **required** for a field to participate in any query. Fields with no `db` tag, or tagged `db:"-"`, are silently ignored.

```
db:"<column>[,option1][,option2]"
```

| Option | Description |
|--------|-------------|
| *(column name only)* | Maps the field to that database column |
| `pk` | Marks the field as the primary key. Used as the default `WHERE` condition in `Update`. |
| `auto` | The field is generated by the database (e.g. `AUTO_INCREMENT`, `SERIAL`, trigger). Excluded from `INSERT` column list; populated back on the entity via `RETURNING` after `Insert`. |
| `-` | Ignores the field entirely — never read from or written to the database |

`pk` and `auto` can be combined. A typical auto-increment primary key is `db:"id,pk,auto"`.

```go
type Order struct {
    ID          int64     `db:"id,pk,auto"`      // primary key, db-generated
    CustomerID  int64     `db:"customer_id,pk"`  // composite PK, not auto
    Status      string    `db:"status"`
    Total       float64   `db:"total"`
    ProcessedAt time.Time `db:"processed_at,auto"` // set by a DB trigger
    Internal    string    `db:"-"`               // never touches the DB
}
```

**Behaviour summary by operation:**

| Tag | INSERT | UPDATE SET | SELECT |
|-----|--------|------------|--------|
| *(plain column)* | included | included | included |
| `pk` | included | excluded from SET, used in WHERE | included |
| `auto` | excluded, populated via RETURNING | excluded | included |
| `pk,auto` | excluded, populated via RETURNING | excluded, used in WHERE | included |
| `-` | excluded | excluded | excluded |

### `schema` tag

The `schema` tag controls how a field value is written during `INSERT`. It is parsed independently of the `db` tag and does not affect `UPDATE` or `SELECT`.

```
schema:"key:value[,key:value]"
```

| Key | Value | Effect |
|-----|-------|--------|
| `default` | `now()` | Writes `NOW()` as a raw SQL literal instead of a bound parameter |

```go
type Event struct {
    ID        int64     `db:"id,pk,auto"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at" schema:"default:now()"`  // → NOW()
    UpdatedAt time.Time `db:"updated_at,auto"`                    // set by trigger
}
```

Generated INSERT:
```sql
INSERT INTO events (name, created_at) VALUES (?, NOW()) RETURNING id
```

> **Note:** Only `default:now()` is currently supported. Other `schema` keys are parsed but have no effect.

### `Tabler` interface

Implement `Tabler` on any struct to provide a custom table name:

```go
type Tabler interface {
    TableName() string
}
```

```go
type OrderLine struct {
    ID       int64   `db:"id,pk,auto"`
    OrderID  int64   `db:"order_id"`
    Product  string  `db:"product"`
    Quantity int     `db:"quantity"`
}

func (OrderLine) TableName() string { return "order_lines" }
```

`TableName()` is checked at runtime via a type assertion on the zero value of `T`. It is only called once per query; there is no caching overhead.

### Table naming convention

When `Tabler` is not implemented, OCA derives the table name automatically:

1. Take the struct type name
2. Lowercase it
3. Append `s`

| Struct | Derived table |
|--------|--------------|
| `User` | `users` |
| `Post` | `posts` |
| `OrderLine` | `orderlines` |
| `Category` | `categorys` |

Irregular plurals and snake_case names require implementing `Tabler`.

---

## Repository

### GenericStore interface

`NewRepository` returns `GenericStore[T]`, which is the full public contract:

```go
type GenericStore[T any] interface {
    Insert(ctx context.Context, entity *T) error
    InsertMany(ctx context.Context, entities []T) error
    Finds(ctx context.Context, filter ...FilterOptions) ([]T, error)
    FindOne(ctx context.Context, opts ...FilterOptions) (*T, error)
    Update(ctx context.Context, entity *T, opts ...FilterOptions) error
    Delete(ctx context.Context, opts ...FilterOptions) error
    DeleteAll(ctx context.Context) error
    Tx(ctx context.Context, fn func(store GenericStore[T]) error) error
}
```

Code against this interface, not `*Repository[T]`, so you can swap in a mock during tests without touching a database.

### Creating a repository

```go
// MySQL (default dialect)
repo := oca.NewRepository[User](db)

// PostgreSQL
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.PostgresDialect{}))

// MariaDB
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.MariaDBDialect{}))
```

`NewRepository` accepts zero or more `RepositoryOption[T]` values. Options are applied in order. Currently the only built-in option is `WithDialect`.

### Dialects

A dialect controls how parameter placeholders are rendered in the generated SQL. Each repository holds its own dialect — there is no global state that needs to be set before the first query.

| Dialect | Constructor | Placeholder style | Example |
|---------|-------------|-------------------|---------|
| MySQL | `query.MySQLDialect{}` | `?` | `WHERE id = ?` |
| MariaDB | `query.MariaDBDialect{}` | `?` | `WHERE id = ?` |
| PostgreSQL | `query.PostgresDialect{}` | `$n` | `WHERE id = $1` |

```go
// Two repositories pointing at different engines in the same process
pgRepo := oca.NewRepository[User](pgDB, oca.WithDialect[User](query.PostgresDialect{}))
myRepo := oca.NewRepository[Log](myDB, oca.WithDialect[Log](query.MySQLDialect{}))
```

The dialect is also inherited by any transaction repository created via `Tx`.

> **Standalone use:** If you use the `query` DSL directly (outside of a repository), call `query.SetDialect(d)` once at startup to set a process-wide default. Repository-level dialects always take precedence over the global.

---

## CRUD Operations

### Insert

Inserts a single record. The method signature:

```go
Insert(ctx context.Context, entity *T) error
```

**Behaviour:**

- Fields tagged `db:"...,auto"` (including `pk,auto`) are **excluded** from the `INSERT` column list. After the insert, they are populated back onto `entity` via a `RETURNING` clause.
- Fields tagged `schema:"default:now()"` inject `NOW()` as a raw SQL literal — the Go field value is not sent.
- All other mapped fields are included as bound parameters.
- Returns a wrapped error on failure: `"oca: insert: <underlying error>"`.

```go
u := &User{Name: "Alice", Email: "alice@example.com"}
if err := repo.Insert(ctx, u); err != nil {
    return err
}
fmt.Println(u.ID) // set by the database
```

Generated SQL (MySQL):
```sql
INSERT INTO users (name, email) VALUES (?, ?) RETURNING id
```

### InsertMany

Inserts multiple records in a **single** `INSERT` statement:

```go
InsertMany(ctx context.Context, entities []T) error
```

**Behaviour:**

- Derives the column layout from the first entity in the slice. All entities must share the same struct type, so the layout is consistent.
- Auto fields are excluded from the column list (the database generates them). The generated values are **not** populated back onto the slice elements — use `Insert` in a loop when you need returned IDs.
- Nil pointer fields are sent as `NULL`.
- Returns early with `nil` if the slice is empty.
- Returns a wrapped error on failure: `"oca: insert many: <underlying error>"`.

```go
err := repo.InsertMany(ctx, []User{
    {Name: "Alice", Email: "alice@example.com"},
    {Name: "Bob",   Email: "bob@example.com"},
    {Name: "Carol", Email: "carol@example.com"},
})
```

Generated SQL (MySQL):
```sql
INSERT INTO users (name, email) VALUES (?, ?), (?, ?), (?, ?)
```

### Finds

Returns all records matching the provided filter options:

```go
Finds(ctx context.Context, filter ...FilterOptions) ([]T, error)
```

**Behaviour:**

- Selects all mapped columns by name in struct-field order.
- Filter options are applied in the order they are passed.
- Multiple `Where()` calls are combined with `AND`.
- Returns `nil, nil` on an empty result set (not an error).
- Returns a wrapped error on failure: `"oca: finds: <underlying error>"`.

```go
users, err := repo.Finds(ctx,
    oca.Where(query.C("active").Eq(true)),
    oca.Where(query.C("role").In("admin", "editor")),
    oca.OrderBy("created_at", oca.DESC),
    oca.Limit(50),
    oca.Offset(100),
)
```

Generated SQL (MySQL):
```sql
SELECT id, name, email, active, role, created_at
FROM users
WHERE active = ? AND role IN (?,?)
ORDER BY created_at DESC
LIMIT ? OFFSET ?
```

### FindOne

Returns a single matching record:

```go
FindOne(ctx context.Context, opts ...FilterOptions) (*T, error)
```

**Behaviour:**

- Internally calls `Finds` with `Limit(1)` appended.
- Returns `*T` on success.
- Returns `nil, sql.ErrNoRows` when no record matches — check with `errors.Is(err, sql.ErrNoRows)`.
- Any other database error is wrapped and returned.

```go
user, err := repo.FindOne(ctx, oca.Where(query.C("email").Eq("alice@example.com")))
if errors.Is(err, sql.ErrNoRows) {
    // record does not exist
}
```

### Update

Updates an existing record in the database:

```go
Update(ctx context.Context, entity *T, opts ...FilterOptions) error
```

**Behaviour:**

- **Primary key field** (`db:"...,pk"`) is excluded from the `SET` clause. Its value is used as the default `WHERE id = ?` condition when no `Where()` opts are provided.
- **Auto fields** (`db:"...,auto"`) are excluded from the `SET` clause entirely.
- **Nil pointer fields** are skipped — that column is left unchanged in the database. This is the mechanism for partial updates.
- **Non-nil pointer fields** are dereferenced and included normally.
- **All other fields** are included in the `SET` clause.
- Returns a wrapped error on failure: `"oca: update: <underlying error>"`.

```go
// Full update — all non-pk, non-auto fields are written
user.Name = "Alice Smith"
user.Email = "alice.smith@example.com"
err := repo.Update(ctx, user)
// UPDATE users SET name = ?, email = ? WHERE id = ?
```

**Partial update using pointer fields:**

Declare optional fields as pointer types. A `nil` pointer is skipped in `UPDATE`; a non-nil pointer is dereferenced and included.

```go
type User struct {
    ID    int64   `db:"id,pk,auto"`
    Name  *string `db:"name"`   // pointer → skipped when nil
    Email string  `db:"email"`  // non-pointer → always included
}

name := "Alice Smith"
u := &User{ID: 1, Name: &name, Email: "alice@example.com"}
repo.Update(ctx, u)
// UPDATE users SET name = ?, email = ? WHERE id = ?

u2 := &User{ID: 1, Name: nil, Email: "alice@example.com"}
repo.Update(ctx, u2)
// UPDATE users SET email = ? WHERE id = ?
// ↑ name column is untouched
```

**Custom WHERE clause:**

Pass one or more `Where()` options to override the default primary-key filter:

```go
err := repo.Update(ctx, user,
    oca.Where(query.C("email").Eq("old@example.com")),
)
// UPDATE users SET name = ?, email = ? WHERE email = ?
```

### Delete

Removes rows matching the provided filter:

```go
Delete(ctx context.Context, opts ...FilterOptions) error
```

**Behaviour:**

- **Requires at least one `Where()` condition.** Calling `Delete` with no options returns an error immediately without executing any SQL. This prevents accidental full-table deletion.
- Multiple `Where()` calls are combined with `AND`.
- Returns a wrapped error on failure: `"oca: delete: <underlying error>"`.

```go
// Delete by primary key
err := repo.Delete(ctx, oca.Where(query.C("id").Eq(user.ID)))

// Delete with multiple conditions (ANDed)
err = repo.Delete(ctx,
    oca.Where(query.C("status").Eq("expired")),
    oca.Where(query.C("updated_at").Lt("2025-01-01")),
)
// DELETE FROM users WHERE status = ? AND updated_at < ?

// Missing Where → returns error, no SQL executed
err = repo.Delete(ctx) // → "oca: Delete requires at least one Where condition…"
```

### DeleteAll

Removes every row in the table. The distinct method name is intentional — it forces you to be explicit about wiping the table:

```go
DeleteAll(ctx context.Context) error
```

```go
err := repo.DeleteAll(ctx)
// DELETE FROM users
```

Returns a wrapped error on failure: `"oca: delete all: <underlying error>"`.

---

## Transactions

`Tx` wraps a function in a single database transaction:

```go
Tx(ctx context.Context, fn func(store GenericStore[T]) error) error
```

**Behaviour:**

- Calls `sql.DB.BeginTx` to start a transaction.
- Passes a transaction-backed `GenericStore[T]` to `fn`. Every method called on it runs within the same transaction.
- If `fn` returns `nil`, the transaction is **committed**.
- If `fn` returns any non-nil error, the transaction is **rolled back** and that error is returned by `Tx`.
- Calling `Tx` on a repository that is already inside a transaction returns `"oca: Tx cannot be nested inside an existing transaction"`.
- The transaction-backed store inherits the parent repository's dialect.

```go
err := repo.Tx(ctx, func(tx oca.GenericStore[User]) error {
    // All operations here share the same transaction

    if err := tx.Insert(ctx, &newUser); err != nil {
        return err // → rollback
    }

    existing.Name = "Updated"
    if err := tx.Update(ctx, &existing); err != nil {
        return err // → rollback
    }

    return nil // → commit
})
```

The `tx` value supports all `GenericStore[T]` methods: `Insert`, `InsertMany`, `Finds`, `FindOne`, `Update`, `Delete`, `DeleteAll`. (Calling `tx.Tx(...)` returns an error since nesting is not allowed.)

---

## Filter Options

Filter options are composable function values passed to `Finds`, `FindOne`, `Update`, and `Delete`. They are applied in the order they are passed.

| Function | Signature | Description |
|----------|-----------|-------------|
| `oca.Where` | `Where(cond query.Condition) FilterOptions` | Appends a WHERE condition. Multiple calls are combined with AND. |
| `oca.OrderBy` | `OrderBy(col string, dir Direction) FilterOptions` | Sets ORDER BY using a typed column name and direction. Does not accept raw strings — safe against injection. |
| `oca.OrderByRaw` | `OrderByRaw(raw string) FilterOptions` | Sets ORDER BY from a raw SQL string. Written directly into SQL — **do not pass user input here**. |
| `oca.Limit` | `Limit(n int) FilterOptions` | Sets LIMIT |
| `oca.Offset` | `Offset(n int) FilterOptions` | Sets OFFSET |

**Direction constants:**

```go
oca.ASC   // "ASC"
oca.DESC  // "DESC"
```

**Example — all options combined:**

```go
results, err := repo.Finds(ctx,
    oca.Where(query.C("status").Eq("active")),
    oca.Where(query.C("role").In("admin", "editor")),
    oca.OrderBy("created_at", oca.DESC),
    oca.Limit(25),
    oca.Offset(50),
)
// SELECT … FROM users
// WHERE status = ? AND role IN (?,?)
// ORDER BY created_at DESC
// LIMIT ? OFFSET ?
```

---

## Query DSL

The `oca/query` package is a standalone SQL builder. The repository uses it internally, but you can import and use it directly for custom queries — particularly useful for `GROUP BY`, window functions, subqueries, or anything that goes beyond what the repository methods cover.

### SELECT builder

```go
sql, args := query.From("users").
    Select("id", "name", "email").           // omit to SELECT *
    Where(query.C("active").Eq(true)).
    Where(query.C("role").In("admin")).
    OrderBy("created_at DESC").              // raw string on the builder
    Limit(10).
    Offset(0).
    Build()
```

| Method | Description |
|--------|-------------|
| `From(table string) *Builder` | Creates a new SELECT builder |
| `.Select(cols ...string)` | Columns to select. Omit to default to `*`. |
| `.Where(conds ...Condition)` | Appends WHERE conditions (ANDed) |
| `.OrderBy(order string)` | Sets raw ORDER BY string |
| `.Limit(n int)` | Sets LIMIT |
| `.Offset(n int)` | Sets OFFSET |
| `.Join(table, on string)` | Adds `INNER JOIN` |
| `.LeftJoin(table, on string)` | Adds `LEFT JOIN` |
| `.RightJoin(table, on string)` | Adds `RIGHT JOIN` |
| `.FullJoin(table, on string)` | Adds `FULL JOIN` |
| `.WithDialect(d Dialect)` | Overrides dialect for this builder |
| `.Build() (string, []any)` | Returns the final SQL and argument slice |

### INSERT builder

```go
// Single row
sql, args := query.InsertInto("users").
    Columns("name", "email", "created_at").
    Values("Alice", "alice@example.com", query.Raw("NOW()")).
    Returning("id", "created_at").
    ToSQL()

// Multi-row: each Values() call adds one row
sql, args = query.InsertInto("users").
    Columns("name", "email").
    Values("Alice", "alice@example.com").
    Values("Bob",   "bob@example.com").
    ToSQL()
// INSERT INTO users (name, email) VALUES (?, ?), (?, ?)
```

| Method | Description |
|--------|-------------|
| `InsertInto(table string) *InsertBuilder` | Creates a new INSERT builder |
| `.Columns(cols ...string)` | Sets the column list |
| `.Values(vals ...any)` | Adds one row of values. Call multiple times for multi-row inserts. Accepts `query.Raw(...)` for SQL literals. |
| `.Returning(cols ...string)` | Appends a `RETURNING` clause (PostgreSQL) |
| `.ReturningID()` | Shorthand for `.Returning("id")` |
| `.WithDialect(d Dialect)` | Overrides dialect for this builder |
| `.ToSQL() (string, []any)` | Returns the final SQL and argument slice. Returns `"", nil` if the builder is incomplete (no table, no columns, or no values). |

### UPDATE builder

```go
sql, args := query.Update("users").
    Set("name", "Alice Smith").
    Set("email", "alice.smith@example.com").
    Where(query.C("id").Eq(1)).
    Build()
// UPDATE users SET name = ?, email = ? WHERE id = ?
// PostgreSQL: UPDATE users SET name = $1, email = $2 WHERE id = $3
```

Placeholder indices are shared across `SET` and `WHERE` so PostgreSQL `$n` numbering is always correct.

| Method | Description |
|--------|-------------|
| `Update(table string) *UpdateBuilder` | Creates a new UPDATE builder |
| `.Set(col string, val any)` | Appends a `col = ?` assignment to SET |
| `.Where(conds ...Condition)` | Appends WHERE conditions (ANDed) |
| `.WithDialect(d Dialect)` | Overrides dialect for this builder |
| `.Build() (string, []any)` | Returns the final SQL and argument slice |

### DELETE builder

```go
sql, args := query.Delete("users").
    Where(query.C("id").Eq(42)).
    Build()
// MySQL:    DELETE FROM users WHERE id = ?
// Postgres: DELETE FROM users WHERE id = $1
```

| Method | Description |
|--------|-------------|
| `Delete(table string) *DeleteBuilder` | Creates a new DELETE builder |
| `.Where(conds ...Condition)` | Appends WHERE conditions (ANDed) |
| `.WithDialect(d Dialect)` | Overrides dialect for this builder |
| `.Build() (string, []any)` | Returns the final SQL and argument slice |

### WHERE conditions

Conditions are built using `query.C(columnName)` and then chaining a comparison method:

```go
query.C("column_name").<method>(args...)
```

| Method | Arguments | SQL produced |
|--------|-----------|-------------|
| `.Eq(val)` | any | `col = ?` |
| `.Neq(val)` | any | `col != ?` |
| `.Gt(val)` | any | `col > ?` |
| `.Gte(val)` | any | `col >= ?` |
| `.Lt(val)` | any | `col < ?` |
| `.Lte(val)` | any | `col <= ?` |
| `.Like(pattern)` | string | `col LIKE ?` |
| `.NotLike(pattern)` | string | `col NOT LIKE ?` |
| `.In(vals...)` | ...any | `col IN (?, ?, …)` |
| `.NotIn(vals...)` | ...any | `col NOT IN (?, ?, …)` |
| `.Between(start, end)` | any, any | `col BETWEEN ? AND ?` |
| `.NotBetween(start, end)` | any, any | `col NOT BETWEEN ? AND ?` |
| `.IsNull()` | — | `col IS NULL` |
| `.IsNotNull()` | — | `col IS NOT NULL` |

```go
query.C("age").Gte(18)
query.C("status").In("active", "pending", "review")
query.C("name").Like("%alice%")
query.C("score").Between(50, 100)
query.C("deleted_at").IsNull()
```

### Logical operators

`And`, `Or`, and `Not` each accept `Condition` values and return a new `Condition`, so they compose freely:

```go
// (age >= ? AND age <= ?)
query.And(
    query.C("age").Gte(18),
    query.C("age").Lte(65),
)

// (role = ? OR role = ?)
query.Or(
    query.C("role").Eq("admin"),
    query.C("role").Eq("editor"),
)

// NOT (active = ?)
query.Not(query.C("active").Eq(false))
```

Nesting works because every operator returns a `Condition`:

```go
query.From("users").
    Where(query.And(
        query.C("active").Eq(true),
        query.Or(
            query.C("role").Eq("admin"),
            query.C("role").Eq("superuser"),
        ),
    )).
    Build()
// … WHERE (active = ? AND (role = ? OR role = ?))
```

### JOINs

```go
sql, args := query.From("orders").
    Select("orders.id", "orders.total", "users.name", "users.email").
    InnerJoin("users", "users.id = orders.user_id").
    LeftJoin("coupons", "coupons.id = orders.coupon_id").
    Where(query.C("orders.status").Eq("pending")).
    OrderBy("orders.created_at DESC").
    Build()
```

| Method | SQL keyword |
|--------|-------------|
| `.Join(table, on)` | `INNER JOIN` |
| `.LeftJoin(table, on)` | `LEFT JOIN` |
| `.RightJoin(table, on)` | `RIGHT JOIN` |
| `.FullJoin(table, on)` | `FULL JOIN` |

The `on` argument is written directly into the SQL. It should be a static string — do not interpolate user input into it.

### RawSQL literals

`query.Raw(string)` wraps a string as a SQL literal that is written directly into the query, bypassing parameter binding. Use it when you need a function call or expression in place of a value:

```go
query.Raw("NOW()")          // → NOW()
query.Raw("CURRENT_DATE")   // → CURRENT_DATE
query.Raw("uuid_generate_v4()") // → uuid_generate_v4()
```

It is accepted anywhere a value is expected in `InsertBuilder.Values`:

```go
query.InsertInto("events").
    Columns("name", "created_at").
    Values("launch", query.Raw("NOW()")).
    ToSQL()
// INSERT INTO events (name, created_at) VALUES (?, NOW())
```

The repository also uses `Raw` internally for `schema:"default:now()"` fields.

### Per-builder dialect override

All four builders (`Builder`, `InsertBuilder`, `UpdateBuilder`, `DeleteBuilder`) accept `.WithDialect(d)` to override the SQL dialect for that specific query, independent of the global default:

```go
// Use PostgreSQL placeholders for this one query only
sql, args := query.From("users").
    WithDialect(query.PostgresDialect{}).
    Where(query.C("id").Eq(1)).
    Build()
// SELECT * FROM users WHERE id = $1
```

When `.WithDialect` is not called, the builder falls back to the global dialect set by `query.SetDialect()`.

---

## Testing

### Mocking with the interface

Because `NewRepository` returns `GenericStore[T]` (an interface), you can inject a mock in unit tests without touching a database:

```go
type mockStore[T any] struct {
    insertFn    func(ctx context.Context, e *T) error
    insertManyFn func(ctx context.Context, es []T) error
    findsFn     func(ctx context.Context, opts ...oca.FilterOptions) ([]T, error)
    // ...
}

func (m *mockStore[T]) Insert(ctx context.Context, e *T) error {
    if m.insertFn != nil {
        return m.insertFn(ctx, e)
    }
    return nil
}
// implement remaining methods...
```

Wire it in:

```go
func TestCreateUser(t *testing.T) {
    store := &mockStore[User]{
        insertFn: func(_ context.Context, u *User) error {
            u.ID = 99 // simulate auto-increment
            return nil
        },
    }

    svc := NewUserService(store)
    user, err := svc.Create(context.Background(), "Alice")
    assert.NoError(t, err)
    assert.Equal(t, int64(99), user.ID)
}
```

### Testing with go-sqlmock

For repository-level tests that verify the exact SQL and arguments generated, use [`go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock):

```go
import (
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/mhdiiilham/oca"
    "github.com/mhdiiilham/oca/query"
)

func TestUserRepo_Insert(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := oca.NewRepository[User](db)

    mock.ExpectQuery(`INSERT INTO users \(name, email\) VALUES \(\?, \?\) RETURNING id`).
        WithArgs("Alice", "alice@example.com").
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    u := &User{Name: "Alice", Email: "alice@example.com"}
    err := repo.Insert(context.Background(), u)

    assert.NoError(t, err)
    assert.Equal(t, int64(1), u.ID)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

For transaction tests, expect `Begin` and `Commit`/`Rollback` in sequence:

```go
mock.ExpectBegin()
mock.ExpectQuery(`INSERT INTO users …`).WillReturnRows(…)
mock.ExpectCommit()

err := repo.Tx(ctx, func(tx oca.GenericStore[User]) error {
    return tx.Insert(ctx, &u)
})
assert.NoError(t, err)
assert.NoError(t, mock.ExpectationsWereMet())
```

---

## License

MIT — see [LICENSE](LICENSE).

---

*OCA is named after my wife. This project started as my way of procrastinating on a problem I kept running into, and turned into something I actually use.*

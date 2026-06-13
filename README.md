# OCA

Lightweight generic ORM for Go — explicit SQL, zero magic.

OCA gives you type-safe CRUD for free while keeping every query readable and
flexible. Struct tags drive the mapping; a fluent query DSL handles the rest.

---

## Features

- Generic `Repository[T]` — full CRUD with a single call to `NewRepository`
- Automatic table naming (`User` → `users`) or custom via `Tabler`
- Auto-increment & schema-default field population after `INSERT`
- **Partial updates** — nil pointer fields are silently skipped in `UPDATE`
- **Batch insert** — `InsertMany` executes a single multi-row `INSERT`
- **Transactions** — `Tx` auto-commits on success, rolls back on error
- **Per-repository dialect** — no global state; MySQL, PostgreSQL, MariaDB
- Fluent query DSL with multi-dialect support
- Rich `WHERE` conditions: `Eq`, `In`, `Like`, `Between`, `And`, `Or`, `Not`, …
- `JOIN` support (INNER, LEFT, RIGHT, FULL)
- Typed `OrderBy(col, dir)` to prevent raw-string injection
- Zero external runtime dependencies (only `database/sql`)

---

## Installation

```bash
go get github.com/mhdiiilham/oca
```

Requires Go 1.21+.

---

## Quick Start

```go
package main

import (
    "context"
    "database/sql"
    "fmt"

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
    db, _ := sql.Open("mysql", "user:pass@/mydb")
    ctx := context.Background()

    repo := oca.NewRepository[User](db)

    // Create
    u := &User{Name: "Alice", Email: "alice@example.com"}
    repo.Insert(ctx, u)
    fmt.Println(u.ID) // populated from RETURNING

    // Batch create
    repo.InsertMany(ctx, []User{{Name: "Bob"}, {Name: "Carol"}})

    // Read one
    found, _ := repo.FindOne(ctx, oca.Where(query.C("id").Eq(u.ID)))

    // Read many
    users, _ := repo.Finds(ctx,
        oca.Where(query.C("name").Like("A%")),
        oca.OrderBy("name", oca.ASC),
        oca.Limit(20),
    )

    // Update (partial: only set fields whose pointer is non-nil)
    found.Name = "Alice Smith"
    repo.Update(ctx, found)

    // Delete (requires at least one Where condition)
    repo.Delete(ctx, oca.Where(query.C("id").Eq(found.ID)))

    // Wipe table
    repo.DeleteAll(ctx)

    _ = users
}
```

---

## Struct Tags

### `db` tag

Maps a struct field to a database column.

| Value | Description |
|-------|-------------|
| `db:"column_name"` | Map to this column |
| `db:"column_name,pk"` | Primary key |
| `db:"column_name,pk,auto"` | Primary key auto-populated after INSERT |
| `db:"column_name,auto"` | Non-PK field auto-populated after INSERT |
| `db:"-"` | Ignore this field |

### `schema` tag

Controls value injection during `INSERT`.

| Value | Description |
|-------|-------------|
| `schema:"default:now()"` | Inserts `NOW()` as a raw SQL literal |

```go
type Post struct {
    ID        int64     `db:"id,pk,auto"`
    Title     string    `db:"title"`
    CreatedAt time.Time `db:"created_at" schema:"default:now()"`
    UpdatedAt time.Time `db:"updated_at,auto"` // set by DB trigger
}
```

### Custom table name

Implement `Tabler` to override the default naming convention:

```go
type OrderLine struct { ... }

func (OrderLine) TableName() string { return "order_lines" }
```

Without `Tabler`, OCA lowercases the struct name and appends `s`
(`User` → `users`, `OrderLine` → `orderlines`).

---

## Creating a Repository

```go
// MySQL (default)
repo := oca.NewRepository[User](db)

// PostgreSQL
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.PostgresDialect{}))

// MariaDB
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.MariaDBDialect{}))
```

The dialect is stored on the repository — there is no global state and no
`SetDialect` call needed. Different repositories can use different dialects.

---

## CRUD Operations

### Insert

```go
u := &User{Name: "Alice", Email: "alice@example.com"}
err := repo.Insert(ctx, u)
// u.ID is now set
```

- Auto fields (`db:"...,auto"`) are populated via `RETURNING` after insert.
- Fields with `schema:"default:now()"` use `NOW()` in the SQL, not a Go value.

### InsertMany

```go
err := repo.InsertMany(ctx, []User{
    {Name: "Alice"},
    {Name: "Bob"},
    {Name: "Carol"},
})
// Executes: INSERT INTO users (name) VALUES (?), (?), (?)
```

- Single `INSERT` statement — more efficient than N individual inserts.
- Auto-increment fields are **not** populated back on the slice elements.
  Use `Insert` in a loop when you need returned IDs.

### Finds

```go
users, err := repo.Finds(ctx,
    oca.Where(query.C("active").Eq(true)),
    oca.OrderBy("created_at", oca.DESC),
    oca.Limit(50),
    oca.Offset(100),
)
```

Returns `[]T`. An empty result set returns `nil, nil`.

### FindOne

```go
user, err := repo.FindOne(ctx, oca.Where(query.C("email").Eq("alice@example.com")))
if errors.Is(err, sql.ErrNoRows) {
    // not found
}
```

Returns `*T` or `sql.ErrNoRows`.

### Update

```go
user.Name = "Alice Smith"
err := repo.Update(ctx, user) // UPDATE users SET name=?, email=? WHERE id=?
```

**Partial updates with pointer fields:**

```go
type User struct {
    ID    int64   `db:"id,pk,auto"`
    Name  *string `db:"name"`
    Email string  `db:"email"`
}

// Only email is updated — Name is nil, so it is skipped entirely.
user := &User{ID: 1, Name: nil, Email: "new@example.com"}
repo.Update(ctx, user) // UPDATE users SET email=? WHERE id=?
```

Override the WHERE clause with `Where()` opts:

```go
err := repo.Update(ctx, user, oca.Where(query.C("email").Eq("old@example.com")))
```

### Delete

```go
err := repo.Delete(ctx, oca.Where(query.C("id").Eq(user.ID)))
```

`Delete` **requires at least one `Where()` condition** and returns an error if
called without one, preventing accidental full-table deletion.

To intentionally delete every row:

```go
err := repo.DeleteAll(ctx)
```

### Tx — Transactions

```go
err := repo.Tx(ctx, func(tx oca.GenericStore[User]) error {
    if err := tx.Insert(ctx, &newUser); err != nil {
        return err  // triggers rollback
    }
    existing.Name = "Updated"
    return tx.Update(ctx, &existing)
    // returning nil triggers commit
})
```

The transaction is automatically committed when `fn` returns `nil`,
or rolled back on any non-nil error.

---

## Filter Options

| Function | Description |
|----------|-------------|
| `oca.Where(cond)` | Add a WHERE condition |
| `oca.OrderBy("col", oca.ASC)` | Typed ORDER BY (safe) |
| `oca.OrderByRaw("col DESC")` | Raw ORDER BY — do not pass user input |
| `oca.Limit(n)` | Set LIMIT |
| `oca.Offset(n)` | Set OFFSET |

---

## Query DSL (`oca/query`)

### Column conditions

```go
query.C("age").Eq(18)           // age = ?
query.C("age").Neq(18)          // age != ?
query.C("age").Gt(18)           // age > ?
query.C("age").Gte(18)          // age >= ?
query.C("age").Lt(65)           // age < ?
query.C("age").Lte(65)          // age <= ?
query.C("name").Like("%alice%") // name LIKE ?
query.C("name").NotLike("%bot%")
query.C("role").In("admin", "editor")    // role IN (?, ?)
query.C("role").NotIn("guest")
query.C("score").Between(50, 100)        // score BETWEEN ? AND ?
query.C("deleted_at").IsNull()
query.C("verified_at").IsNotNull()
```

### Logical operators

```go
query.And(query.C("age").Gte(18), query.C("age").Lte(65))
// → (age >= ? AND age <= ?)

query.Or(query.C("role").Eq("admin"), query.C("role").Eq("editor"))
// → (role = ? OR role = ?)

query.Not(query.C("active").Eq(false))
// → NOT (active = ?)
```

### JOINs

```go
sql, args := query.From("orders").
    Select("orders.id", "users.name").
    InnerJoin("users", "users.id = orders.user_id").
    Where(query.C("orders.status").Eq("pending")).
    Build()
```

Available join types: `InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin`.

### Building queries directly

```go
// SELECT
sql, args := query.From("users").
    Select("id", "name").
    Where(query.C("active").Eq(true)).
    OrderBy("name ASC").
    Limit(10).
    Build()

// INSERT
sql, args := query.InsertInto("users").
    Columns("name", "email").
    Values("Alice", "alice@example.com").
    Returning("id").
    ToSQL()

// UPDATE
sql, args := query.Update("users").
    Set("name", "Bob").
    Where(query.C("id").Eq(1)).
    Build()

// DELETE
sql, args := query.Delete("users").
    Where(query.C("id").Eq(1)).
    Build()
```

All builders support `.WithDialect(d)` to override the dialect for that query.

---

## Dialects

Each repository holds its own dialect — no global state.

```go
// PostgreSQL ($1, $2, ...)
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.PostgresDialect{}))

// MariaDB
repo := oca.NewRepository[User](db, oca.WithDialect[User](query.MariaDBDialect{}))

// MySQL (default when WithDialect is not called)
repo := oca.NewRepository[User](db)
```

The standalone query DSL still supports a global default via `query.SetDialect()`
for standalone use outside of a repository.

---

## Mocking in Tests

`GenericStore[T]` is an interface — mock it directly:

```go
type mockStore[T any] struct{ ... }
func (m *mockStore[T]) Insert(ctx context.Context, e *T) error { ... }
// implement remaining methods
```

Or use `go-sqlmock` against a real `*Repository`:

```go
db, mock, _ := sqlmock.New()
repo := oca.NewRepository[User](db)

mock.ExpectQuery(`INSERT INTO users ...`).
    WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
```

---

## License

MIT

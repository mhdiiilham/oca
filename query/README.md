# Query Builder (oca/query)

A lightweight SQL query builder for Go, designed to be simple, composable, and easy to extend.
It generates parameterized SQL queries and supports a fluent API style.

# Road Map

See [ROADMAP.md](./ROADMAP.md) for planned features and progress.

## Features

- Chainable builder API
- Safe parameter binding (? placeholders)
- Supports:
  - `SELECT` with custom columns or *
  - `INSERT` (`InsertInto`, `Columns`, `Values`)
  - `JOIN` (INNER, LEFT, RIGHT, FULL
  - `WHERE` (multiple conditions with AND)
  - `ORDER BY`
  - `LIMIT / OFFSET`

## Installation

```bash
go get github.com/mhdiiilham/oca/query
```

## Usage

### INSERT

```go
q := query.InsertInto("users").
    Columns("id", "name", "email").
    Values(1, "Alice", "alice@example.com").
    Values(2, "Bob", "bob@example.com")

sql, args := q.ToSQL()

// sql:  "INSERT INTO users (id, name, email) VALUES (?, ?, ?), (?, ?, ?)"
// args: [1 "Alice" "alice@example.com" 2 "Bob" "bob@example.com"]
```

## INSERT With Returning

```go
q := query.InsertInto("users").
    Columns("name", "email").
    Values("Alice", "alice@example.com").
    ReturningID()

sql, args := q.ToSQL()
// sql:  "INSERT INTO users (name, email) VALUES (?, ?) RETURNING id"
// args: ["Alice", "alice@example.com"]
```

### Basic SELECT

```go
sql, args := query.From("users").
    Select("id", "name").
    Build()

// sql:  "SELECT id, name FROM users"
// args: []
```

### Join

#### INNER JOIN

```go
sql, args := query.From("users").
    Select("users.id", "profiles.bio").
    Join("profiles", "users.id = profiles.user_id").
    Build()

// sql:  "SELECT users.id, profiles.bio FROM users INNER JOIN profiles ON users.id = profiles.user_id"
// args: []
```

#### LEFT JOIN

```go
sql, args := query.From("orders").
    Select("orders.id", "customers.name").
    LeftJoin("customers", "orders.customer_id = customers.id").
    Build()

// sql:  "SELECT orders.id, customers.name FROM orders LEFT JOIN customers ON orders.customer_id = customers.id"
// args: []
```

#### RIGHT JOIN

```go
sql, args := query.From("a").
    Select("a.x", "b.y").
    RightJoin("b", "a.id = b.a_id").
    Build()

// sql:  "SELECT a.x, b.y FROM a RIGHT JOIN b ON a.id = b.a_id"
// args: []
```

#### FULL JOIN

```go
sql, args := query.From("products").
    Select("products.id", "categories.name").
    FullJoin("categories", "products.category_id = categories.id").
    Build()

// sql:  "SELECT products.id, categories.name FROM products FULL JOIN categories ON products.category_id = categories.id"
// args: []
```

### SELECT with WHERE

```go
sql, args := query.From("users").
    Select("id", "email").
    Where("age > ?", 18).
    Build()

// sql:  "SELECT id, email FROM users WHERE age > ?"
// args: [18]
```

### Multiple WHERE conditions

```go
sql, args := query.From("users").
    Select("id").
    Where("age > ?", 18).
    Where("status = ?", "active").
    Build()

// sql:  "SELECT id FROM users WHERE age > ? AND status = ?"
// args: [18, "active"]
```

### ORDER BY, LIMIT, OFFSET

```go
sql, args := query.From("users").
    Select("id", "email").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(5).
    Build()

// sql:  "SELECT id, email FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?"
// args: [10, 5]
```

### Default SELECT *

```go
sql, args := query.From("products").Build()

// sql:  "SELECT * FROM products"
// args: []
```
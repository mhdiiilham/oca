# OCA

OCA is a lightweight, generic ORM for Go, designed out of necessity and nostalgia.
Coming from a background working with ORMs like Sequelize (JavaScript), JPA/Hibernate (Java), and Eloquent (Laravel/PHP), I always appreciated how straightforward and expressive they made database interactions.

When I started working more with Go, I tried GORM, but found it too feature-heavy for my personal taste. What I wanted was something in between:

- Simple enough to understand and extend without hidden magic.
- Generic and type-safe, making full use of Go’s modern generics.
- Practical for building real-world applications without boilerplate overload.

That’s why OCA exists: to bring back the developer experience I enjoyed from other ecosystems, but in a way that fits naturally into the Go world.

# Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
    - [Defining a Repository](#defining-a-repository)
    - [Insert](#insert)
- [LICENSE](./LICENCE)
- [Code of Conduct](./CODE_OF_CONDUCT.md)
- [Changelog](./CHANGELOG.md)

# Features

- Lightweight and minimal
- Built with generics for type safety
- Inspired by Sequelize, JPA, and Eloquent

# Requirements

- Go **1.18+** (generics support required)


# Installation

```bash
go get github.com/yourusername/oca
```

# Usage

OCA provides a generic Repository[T] to perform common CRUD operations on any struct-based model.
Below are the main APIs with examples.

## Defining a Repository

```go
repo := oca.NewRepository[User](db)
```

- T is your model type (e.g. User).
- db is a standard *sql.DB connection.
- You can now call Insert, Find, FindOne, Update, and Delete.

## Insert
To add a new record to the database, define your model struct with proper schema tags.
If the primary key is auto-increment, the ORM will update your struct’s ID after a successful insert.

```go
// User represents the "users" table in the database.
// The `schema` tags define column mappings and behavior.
//
// Special tags:
// - `pk`       marks this field as the primary key
// - `auto`     indicates auto-increment, so Insert will update this field
// - `default`  skips inserting this column if it's zero-valued in Go,
//              letting the DB default apply (e.g. NOW(), true, 'guest').
type User struct {
    ID        int       `schema:"id,pk,auto"`                 // Auto-increment primary key
    Name      string    `schema:"name"`                       // Column: name
    Age       int       `schema:"age"`                        // Column: age
    Active    bool      `schema:"active,default:true"`        // DB default: true
    Role      string    `schema:"role,default:'guest'"`       // DB default: 'guest'
    CreatedAt time.Time `schema:"created_at,default:NOW()"`   // DB default: current timestamp
}

func main() {
    ctx := context.Background()

    newUser := User{
        Name: "Alice",
        Age:  25,
        // Active, Role, CreatedAt omitted → DB will apply defaults
    }

    // Insert the new user into the database.
    // - On success, newUser.ID will be updated with the generated ID.
    // - Default fields (Active, Role, CreatedAt) will be populated
    //   with the values assigned by the database.
    if err := repo.Insert(ctx, newUser); err != nil {
        panic(err)
    }

    fmt.Println("Inserted user with ID:", newUser.ID)
    fmt.Println("Role:", newUser.Role)           // "guest"
    fmt.Println("Active:", newUser.Active)       // true
    fmt.Println("Created at:", newUser.CreatedAt) // timestamp from DB
}
```

_Notes_:
- `pk` is required for primary keys.
- `auto` updates the field after insert (e.g. auto-increment ID).
- `default` tells OCA to let the DB fill in the field if it’s zero-valued in Go.
- After Insert, both `auto` and default fields are read back into your struct.
- default tag supports SQL expressions (NOW()), booleans (true, false), and string literals ('guest').
- If a field is zero-valued in Go (e.g. 0, false, "", time.Time{}), OCA omits it from the INSERT and lets the DB assign the default.
- If you explicitly set the field in Go, your value overrides the DB default.

package oca

import (
	"context"
	"reflect"
	"strings"

	"github.com/mhdiiilham/oca/query"
)

// Insert inserts a new record into the database and sets auto-increment ID or other auto fields if applicable.
// It also respects schema defaults, e.g., `schema:"default:now()"`.
//
// Example:
//
//	type User struct {
//	    ID        int64  `db:"id,pk,auto"`
//	    Name      string `db:"name"`
//	    CreatedAt time.Time `db:"created_at" schema:"default:now()"`
//	}
//
//	repo := NewRepository[User](db)
//	u := &User{Name: "Alice"}
//	err := repo.Insert(ctx, u)
//
// u.ID and u.CreatedAt will be automatically set if applicable.
func (r *Repository[T]) Insert(ctx context.Context, entity *T) error {
	fields := parseFields(entity)

	var cols []string
	var args []any
	var autoFields []FieldMeta

	// Collect columns, args, and auto fields
	for _, f := range fields {
		if f.IsAuto {
			autoFields = append(autoFields, f.FieldMeta)
			continue
		}

		if defaultVal, ok := f.Schema["default"]; ok && defaultVal == "now()" {
			cols = append(cols, f.Column)
			args = append(args, query.Raw("NOW()"))
			continue
		}

		cols = append(cols, f.Column)
		args = append(args, f.Value)
	}

	var t T
	builder := query.InsertInto(resolveTableName(t)).
		Columns(cols...).
		Values(args...)

	// Add RETURNING for all auto fields
	if len(autoFields) > 0 {
		autoCols := make([]string, len(autoFields))
		for i, f := range autoFields {
			autoCols[i] = f.Column
		}
		builder = builder.Returning(autoCols...)
	}

	sqlStr, sqlArgs := builder.ToSQL()

	if len(autoFields) > 0 {
		// Use QueryRow because we expect RETURNING values
		row := r.db.QueryRowContext(ctx, sqlStr, sqlArgs...)

		// Prepare pointers to struct fields for scanning
		valOfEntity := reflect.ValueOf(entity).Elem()
		scanTargets := make([]interface{}, len(autoFields))
		for i, f := range autoFields {
			field := valOfEntity.Field(f.Index)
			if !field.CanAddr() {
				return &reflect.ValueError{Method: "Insert", Kind: field.Kind()}
			}
			scanTargets[i] = field.Addr().Interface()
		}

		if err := row.Scan(scanTargets...); err != nil {
			return err
		}

		return nil
	}

	// Normal exec (no auto-returning)
	_, err := r.db.ExecContext(ctx, sqlStr, sqlArgs...)
	return err
}

// resolveTableName returns the table name either from Tabler or from struct name.
func resolveTableName[T any](entity T) string {
	if t, ok := any(entity).(Tabler); ok {
		return t.TableName()
	}

	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Fallback: use lowercase struct name
	return strings.ToLower(typ.Name()) + "s"
}

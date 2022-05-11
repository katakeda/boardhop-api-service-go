package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Category struct {
	Id       *int    `db:"id"`
	ParentId *int    `db:"parent_id"`
	Name     *string `db:"name"`
	Path     *string `db:"path"`
}

func (r *Repository) GetCategories(ctx context.Context) (categories []Category, err error) {
	tx := ctx.Value(TxnKey).(pgx.Tx)
	if tx == nil {
		tx, _ = r.db.Begin(ctx)
		defer func() error {
			if err != nil {
				return tx.Rollback(ctx)
			}
			return tx.Commit(ctx)
		}()
	}

	cols := []string{
		"id",
		"parent_id",
		"name",
		"path",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From("category").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	for rows.Next() {
		c := Category{}
		rows.Scan(&c.Id, &c.ParentId, &c.Name, &c.Path)
		categories = append(categories, c)
	}

	return categories, nil
}

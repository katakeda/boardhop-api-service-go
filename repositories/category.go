package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Category struct {
	Id       *int    `json:"id" db:"id"`
	ParentId *int    `json:"parentId" db:"parent_id"`
	Path     *string `json:"path" db:"path"`
	Value    *string `json:"value" db:"value"`
	Label    *string `json:"label" db:"label"`
}

func (r *Repository) GetCategories(ctx context.Context) (categories []Category, err error) {
	tx, ok := ctx.Value(TxnKey).(pgx.Tx)
	if !ok || tx == nil {
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
		"path",
		"value",
		"label",
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
		if err := rows.Scan(&c.Id, &c.ParentId, &c.Path, &c.Value, &c.Label); err != nil {
			return nil, fmt.Errorf("failed to scan rows | %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type Category struct {
	Id       *int    `db:"id"`
	ParentId *int    `db:"parent_id"`
	Name     *string `db:"name"`
	Path     *string `db:"path"`
}

func (r *Repository) GetCategories() ([]Category, error) {
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
		return nil, fmt.Errorf("failed to build query: %s %w", sqlStmt, err)
	}

	var categories []Category
	{
		err := pgxscan.Select(context.Background(), r.db, &categories, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
		}
	}

	return categories, nil
}

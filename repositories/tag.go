package repositories

import (
	"context"
	"fmt"
	"net/url"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Tag struct {
	Id    int    `json:"id" db:"id"`
	Type  string `json:"type" db:"type"`
	Value string `json:"value" db:"value"`
	Label string `json:"label" db:"label"`
}

func (r *Repository) GetTags(ctx context.Context, params url.Values) (tags []Tag, err error) {
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
		"type",
		"value",
		"label",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(cols...).
		From("tag")

	if rootType := params.Get("type"); rootType == "snowboard" {
		psql = psql.Where(sq.Eq{"type": []string{"Snowboard Brand", "Skill Level"}})
	} else {
		psql = psql.Where(sq.Eq{"type": []string{"Surfboard Brand", "Skill Level"}})
	}

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	for rows.Next() {
		t := Tag{}
		rows.Scan(&t.Id, &t.Type, &t.Value, &t.Label)
		tags = append(tags, t)
	}

	return tags, nil
}

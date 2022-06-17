package repositories

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/katakeda/boardhop-api-service-go/utils"
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

	types := []string{"Skill Level"}
	typeArr := strings.Split(params.Get("type"), ",")
	typeMap := utils.StrArrayToMap(typeArr)
	if _, exists := typeMap["surfboard"]; exists {
		types = append(types, "Surfboard Brand")
	}
	if _, exists := typeMap["snowboard"]; exists {
		types = append(types, "Snowboard Brand")
	}
	psql = psql.Where(sq.Eq{"type": types})

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if err := pgxscan.ScanAll(&tags, rows); err != nil {
		return nil, fmt.Errorf("failed to scan rows | %w", err)
	}

	return tags, nil
}

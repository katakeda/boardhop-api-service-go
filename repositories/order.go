package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Order struct {
	Id        string     `json:"id" db:"id"`
	PostId    string     `json:"postId" db:"post_id"`
	UserId    string     `json:"userId" db:"user_id"`
	PaymentId string     `json:"paymentId" db:"payment_id"`
	Status    string     `json:"status" db:"status"`
	Message   string     `json:"message" db:"message"`
	Quantity  int8       `json:"quantity" db:"quantity"`
	Total     float32    `json:"total" db:"total"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type CreateOrderPayload struct {
	PostId    string  `json:"postId"`
	UserId    string  `json:"userId"`
	PaymentId string  `json:"paymentId"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
	Quantity  int8    `json:"quantity"`
	Total     float32 `json:"total"`
}

func (r *Repository) CreateOrder(ctx context.Context, payload CreateOrderPayload) (order *Order, err error) {
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
		"post_id",
		"user_id",
		"payment_id",
		"status",
		"message",
		"quantity",
		"total",
	}

	vals := []interface{}{
		payload.PostId,
		payload.UserId,
		payload.PaymentId,
		payload.Status,
		payload.Message,
		payload.Quantity,
		payload.Total,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Insert(`"order"`).
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	var newOrder Order
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(&newOrder.Id); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return &newOrder, nil
}

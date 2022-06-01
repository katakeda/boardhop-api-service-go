package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Order struct {
	Id        string     `json:"id" db:"id"`
	PostId    string     `json:"postId" db:"post_id"`
	UserId    string     `json:"userId" db:"user_id"`
	PaymentId string     `json:"paymentId" db:"payment_id"`
	Status    string     `json:"status" db:"status"`
	Quantity  int8       `json:"quantity" db:"quantity"`
	Total     float32    `json:"total" db:"total"`
	StartDate *time.Time `json:"startDate" db:"start_date"`
	EndDate   *time.Time `json:"endDate" db:"end_date"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	Post      Post       `json:"post"`
	Messages  []Message  `json:"messages"`
}

type CreateOrderPayload struct {
	PostId    string  `json:"postId"`
	UserId    string  `json:"userId"`
	PaymentId string  `json:"paymentId"`
	Status    string  `json:"status"`
	Quantity  int8    `json:"quantity"`
	Total     float32 `json:"total"`
	Message   *string `json:"message"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

type GetOrdersFilter struct {
	UserId *string
}

func (r *Repository) GetOrders(ctx context.Context, filter GetOrdersFilter) (orders []Order, err error) {
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

	if filter.UserId == nil {
		return nil, fmt.Errorf("userId cant be empty")
	}

	cols := []string{
		"id",
		"post_id",
		"user_id",
		"payment_id",
		"status",
		"quantity",
		"total",
		"created_at",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(cols...).
		From(`"order"`).
		Where(sq.Eq{"user_id": filter.UserId})

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if err := pgxscan.ScanAll(&orders, rows); err != nil {
		return nil, fmt.Errorf("failed to scan rows | %w", err)
	}

	return orders, nil
}

func (r *Repository) GetOrder(ctx context.Context, id string) (order *Order, err error) {
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
		"post_id",
		"user_id",
		"payment_id",
		"status",
		"quantity",
		"total",
		"created_at",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(cols...).
		From(`"order"`).
		Where(sq.Eq{"id": id})

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	order = &Order{}
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(
		&order.Id,
		&order.PostId,
		&order.UserId,
		&order.PaymentId,
		&order.Status,
		&order.Quantity,
		&order.Total,
		&order.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if err := r.setOrderPost(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to set order post | %w", err)
	}

	if err := r.setOrderMessages(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to set order message | %w", err)
	}

	return order, nil
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
		"quantity",
		"total",
		"start_date",
		"end_date",
	}

	vals := []interface{}{
		payload.PostId,
		payload.UserId,
		payload.PaymentId,
		payload.Status,
		payload.Quantity,
		payload.Total,
		payload.StartDate,
		payload.EndDate,
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

func (r *Repository) setOrderPost(ctx context.Context, order *Order) (err error) {
	post, err := r.GetPost(ctx, order.PostId)
	if err != nil {
		return fmt.Errorf("failed to get order post | %w", err)
	}

	order.Post = *post

	return nil
}

func (r *Repository) setOrderMessages(ctx context.Context, order *Order) (err error) {
	messages, err := r.GetMessagesByOrderId(ctx, order.Id)
	if err != nil {
		return fmt.Errorf("failed to get order messages | %w", err)
	}

	order.Messages = messages

	return nil
}

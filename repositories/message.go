package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Message struct {
	Id        int        `json:"id" db:"id"`
	UserId    string     `json:"userId" db:"user_id"`
	PostId    *string    `json:"postId" db:"post_id"`
	OrderId   *string    `json:"orderId" db:"order_id"`
	Message   *string    `json:"message" db:"message"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`

	AvatarUrl *string `json:"avatarUrl" db:"avatar_url"`
}

type CreateMessagePayload struct {
	UserId  string
	PostId  *string `json:"postId"`
	OrderId *string `json:"orderId"`
	Message *string `json:"message"`
}

func (r *Repository) GetMessagesByOrderId(ctx context.Context, orderId string) (messages []Message, err error) {
	return r.getMessages(ctx, orderId, "order")
}

func (r *Repository) GetMessagesByPostId(ctx context.Context, postId string) (messages []Message, err error) {
	return r.getMessages(ctx, postId, "post")
}

func (r *Repository) CreateMessage(ctx context.Context, payload CreateMessagePayload) (message *Message, err error) {
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
		"user_id",
		"post_id",
		"order_id",
		"message",
	}

	vals := []interface{}{
		payload.UserId,
		payload.PostId,
		payload.OrderId,
		payload.Message,
	}

	sqlStmt, sqlArgs, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("message").
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	var newMessage Message
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(&newMessage.Id); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return &newMessage, nil
}

func (r *Repository) getMessages(ctx context.Context, id string, by string) (messages []Message, err error) {
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
		"a.id",
		"a.user_id",
		"a.post_id",
		"a.order_id",
		"a.message",
		"a.created_at",
		"b.avatar_url",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(cols...).
		From("message a").
		Join(`"user" b ON a.user_id = b.id`)

	if by == "order" {
		psql = psql.Where(sq.Eq{"a.order_id": id})
	} else {
		psql = psql.Where(sq.Eq{"a.post_id": id})
	}

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if err := pgxscan.ScanAll(&messages, rows); err != nil {
		return nil, fmt.Errorf("failed to scan rows | %w", err)
	}

	return messages, nil
}

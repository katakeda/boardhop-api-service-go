package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func (r *Repository) GetMessagesByOrderId(ctx context.Context, orderId string) (messages []Message, err error) {
	return r.getMessages(ctx, orderId, "order")
}

func (r *Repository) GetMessagesByPostId(ctx context.Context, postId string) (messages []Message, err error) {
	return r.getMessages(ctx, postId, "post")
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

	for rows.Next() {
		m := Message{}
		if err := rows.Scan(&m.Id, &m.UserId, &m.PostId, &m.OrderId, &m.Message, &m.CreatedAt, &m.AvatarUrl); err != nil {
			return nil, fmt.Errorf("failed to scan rows | %w", err)
		}
		messages = append(messages, m)
	}

	return messages, nil
}

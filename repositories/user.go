package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Id           string  `json:"id" db:"id"`
	Email        string  `json:"email" db:"email"`
	FirstName    string  `json:"firstName" db:"first_name"`
	LastName     string  `json:"lastName" db:"last_name"`
	Phone        *string `json:"phone" db:"phone"`
	AvatarUrl    *string `json:"avatarUrl" db:"avatar_url"`
	GoogleAuthId *string `json:"googleAuthId" db:"google_auth_id"`
}

type UserSignupPayload struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	GoogleAuthId string `json:"googleAuthId"`
}

type UserLoginPayload struct {
	Email        string `json:"email"`
	GoogleAuthId string `json:"googleAuthId"`
}

func (r *Repository) UserSignup(ctx context.Context, payload UserSignupPayload) (user *User, err error) {
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
		"email",
		"first_name",
		"last_name",
		"google_auth_id",
	}

	vals := []interface{}{
		payload.Email,
		payload.FirstName,
		payload.LastName,
		payload.GoogleAuthId,
	}

	sqlStmt, sqlArgs, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(`"user"`).
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	var newUser User
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(&newUser.Id); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return &newUser, nil
}

func (r *Repository) GetUserByGoogleAuthId(ctx context.Context, googleAuthId interface{}) (*User, error) {
	cols := []string{
		"id",
		"email",
		"first_name",
		"last_name",
		"google_auth_id",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From(`"user"`).
		Where(sq.Eq{"google_auth_id": googleAuthId}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s %w", sqlStmt, err)
	}

	var user User
	{
		err := pgxscan.Get(ctx, r.db, &user, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
		}
	}

	return &user, nil
}

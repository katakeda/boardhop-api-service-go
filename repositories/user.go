package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
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

func (r *Repository) UserSignup(ctx context.Context, payload UserSignupPayload) (*User, error) {
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

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Insert(`"user"`).
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING id, email, first_name, last_name, phone, avatar_url, google_auth_id").
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

func (r *Repository) GetUserByGoogleAuthId(ctx context.Context, googleAuthId interface{}) (*User, error) {
	cols := []string{
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

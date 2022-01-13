package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) (*Repository, error) {
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

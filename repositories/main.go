package repositories

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IRepository interface {
	GetPosts(ctx context.Context, params url.Values) ([]Post, error)
	GetPost(ctx context.Context, id string) (*Post, error)
	CreatePost(ctx context.Context, payload CreatePostPayload) (*Post, error)

	GetCategories() ([]Category, error)

	UserSignup(ctx context.Context, payload UserSignupPayload) (*User, error)
	GetUserByGoogleAuthId(ctx context.Context, googleAuthId interface{}) (*User, error)
}

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

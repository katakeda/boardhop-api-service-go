package repositories

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CtxKey string

const (
	TxnKey CtxKey = "txnKey"
)

type IRepository interface {
	BeginTxn(ctx context.Context) (context.Context, error)
	CommitTxn(ctx context.Context) error
	RollbackTxn(ctx context.Context) error

	GetPosts(ctx context.Context, params url.Values) ([]Post, error)
	GetPost(ctx context.Context, id string) (*Post, error)
	CreatePost(ctx context.Context, payload CreatePost) (*Post, error)
	UpdatePost(ctx context.Context, id string, payload UpdatePost) (*Post, error)
	CreatePostTags(ctx context.Context, tags []CreatePostTag) error
	CreatePostMedias(ctx context.Context, medias []CreatePostMedia) error
	CreatePostCategories(ctx context.Context, categories []CreatePostCategory) error
	DeletePostTags(ctx context.Context, id string) error
	DeletePostCategories(ctx context.Context, id string) error

	GetTags(ctx context.Context, params url.Values) ([]Tag, error)

	GetCategories(ctx context.Context) ([]Category, error)

	UserSignup(ctx context.Context, payload UserSignupPayload) (*User, error)
	GetUserByGoogleAuthId(ctx context.Context, googleAuthId interface{}) (*User, error)

	GetOrders(ctx context.Context, filter GetOrdersFilter) ([]Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	CreateOrder(ctx context.Context, payload CreateOrderPayload) (*Order, error)

	CreateMessage(ctx context.Context, payload CreateMessagePayload) (*Message, error)
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

func (r Repository) BeginTxn(ctx context.Context) (context.Context, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, TxnKey, tx), nil
}

func (r Repository) CommitTxn(ctx context.Context) error {
	tx := ctx.Value(TxnKey)
	if tx == nil {
		return fmt.Errorf("failed to get txn from ctx")
	}

	return tx.(pgx.Tx).Commit(ctx)
}

func (r Repository) RollbackTxn(ctx context.Context) error {
	tx := ctx.Value(TxnKey)
	if tx == nil {
		return fmt.Errorf("failed to get txn from ctx")
	}

	return tx.(pgx.Tx).Rollback(ctx)
}

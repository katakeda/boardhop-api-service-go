package repositories

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

const (
	SURFBOARD_PATH = "root.1"
	SNOWBOARD_PATH = "root.2"
	PER_PAGE_MAX   = 25
)

type Post struct {
	Id              string     `json:"id" db:"id"`
	UserId          string     `json:"userId" db:"user_id"`
	Title           string     `json:"title" db:"title"`
	Price           float32    `json:"price" db:"price"`
	Rate            string     `json:"rate" db:"rate"`
	Description     *string    `json:"description" db:"description"`
	PickupLatitude  *float64   `json:"pickupLatitude" db:"pickup_latitude"`
	PickupLongitude *float64   `json:"pickupLongitude" db:"pickup_longitude"`
	CreatedAt       *time.Time `json:"createdAt" db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`

	Email      *string `json:"email" db:"email"`
	AvatarUrl  *string `json:"avatarUrl" db:"avatar_url"`
	Categories *string `json:"categories" db:"categories"`
	Tags       *string `json:"tags" db:"tags"`
}

type CreatePostPayload struct {
	UserId          string   `json:"userId"`
	Title           string   `json:"title"`
	Price           float32  `json:"price"`
	Rate            string   `json:"rate"`
	Description     *string  `json:"description"`
	PickupLatitude  *float64 `json:"pickupLatitude"`
	PickupLongitude *float64 `json:"pickupLongitude"`
}

// TODO: Replace params with filters
func (r *Repository) GetPosts(ctx context.Context, params url.Values) ([]Post, error) {
	var rootPath string
	if params.Get("type") == "snowboard" {
		rootPath = SNOWBOARD_PATH
	} else {
		rootPath = SURFBOARD_PATH
	}

	cols := []string{
		"a.id",
		"a.user_id",
		"a.title",
		"a.price",
		"a.rate",
		"a.pickup_latitude",
		"a.pickup_longitude",
		"a.created_at",
		"b.email",
		"b.avatar_url",
		`string_agg(DISTINCT d. "name", ',') AS categories`,
		`string_agg(DISTINCT f. "value", ',') AS tags`,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlBuilder := psql.Select(cols...).From("post a").
		Join(`"user" b ON a.user_id = b.id`).
		Join("post_category c ON a.id = c.post_id").
		Join("category d ON c.category_id = d.id").
		LeftJoin("post_tag e ON a.id = e.post_id").
		LeftJoin("tag f ON e.tag_id = f.id").
		LeftJoin("tag_type g ON f.type_id = g.id").
		Where("d.path <@ ?", rootPath)

	if categories := params.Get("cats"); categories != "" {
		sqlBuilder = sqlBuilder.Where(sq.Eq{"d.name": strings.Split(categories, ",")})
	}

	if tags := params.Get("tags"); tags != "" {
		sqlBuilder = sqlBuilder.Where(sq.Eq{"f.value": strings.Split(tags, ",")})
	}

	offset := 0
	if page := params.Get("p"); page != "" {
		var err error
		offset, err = strconv.Atoi(page)
		if err != nil {
			offset = 0
		}
	}

	var sqlStmt string
	var sqlArgs []interface{}
	{
		var err error
		sqlStmt, sqlArgs, err = sqlBuilder.Offset(uint64(offset)*PER_PAGE_MAX).
			Limit(PER_PAGE_MAX).
			GroupBy("a.id", "b.id").
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query: %s %w", sqlStmt, err)
		}
	}

	var posts []Post
	err := pgxscan.Select(ctx, r.db, &posts, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
	}

	return posts, nil
}

func (r *Repository) GetPost(ctx context.Context, id string) (*Post, error) {
	cols := []string{
		"a.id",
		"a.user_id",
		"a.title",
		"a.description",
		"a.price",
		"a.rate",
		"a.pickup_latitude",
		"a.pickup_longitude",
		"a.created_at",
		"b.email",
		"b.avatar_url",
		`string_agg(DISTINCT d. "name", ',') AS categories`,
		`string_agg(DISTINCT f. "value", ',') AS tags`,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From("post a").
		Join(`"user" b ON a.user_id = b.id`).
		Join("post_category c ON a.id = c.post_id").
		Join("category d ON c.category_id = d.id").
		LeftJoin("post_tag e ON a.id = e.post_id").
		LeftJoin("tag f ON e.tag_id = f.id").
		LeftJoin("tag_type g ON f.type_id = g.id").
		Where(sq.Eq{"a.id": id}).
		GroupBy("a.id", "b.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s %w", sqlStmt, err)
	}

	var post Post
	{
		err := pgxscan.Get(ctx, r.db, &post, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
		}
	}

	return &post, nil
}

func (r *Repository) CreatePost(ctx context.Context, payload CreatePostPayload) (*Post, error) {
	cols := []string{
		"user_id",
		"title",
		"price",
		"rate",
		"pickup_latitude",
		"pickup_longitude",
	}

	vals := []interface{}{
		payload.UserId,
		payload.Title,
		payload.Price,
		payload.Rate,
		payload.PickupLatitude,
		payload.PickupLongitude,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Insert("post").
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s %w", sqlStmt, err)
	}

	var post Post
	{
		err := pgxscan.Get(ctx, r.db, &post, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
		}
	}

	return &post, nil
}

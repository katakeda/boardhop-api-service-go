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
	Id              int        `db:"id"`
	UserId          int        `db:"user_id"`
	Title           string     `db:"title"`
	Price           float32    `db:"price"`
	Rate            string     `db:"rate"`
	Description     *string    `db:"description"`
	PickupLatitude  *float64   `db:"pickup_latitude"`
	PickupLongitude *float64   `db:"pickup_longitude"`
	CreatedAt       *time.Time `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`

	Username   *string `db:"username"`
	AvatarUrl  *string `db:"avatar_url"`
	Categories *string `db:"categories"`
	Tags       *string `db:"tags"`
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
		"b.username",
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
		"b.username",
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
		err := pgxscan.Get(context.Background(), r.db, &post, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s %w", sqlStmt, err)
		}
	}

	return &post, nil
}

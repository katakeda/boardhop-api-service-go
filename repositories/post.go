package repositories

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

const (
	DEFAULT_LIMIT = 50
)

var (
	CATEGORY_PATH = map[string]string{
		"all":       "root",
		"surfboard": "root.1",
		"snowboard": "root.2",
		"eps":       "root.1.3",
		"pu":        "root.1.3",
		"wideboard": "root.2.4",
	}
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

	Email      *string     `json:"email" db:"email"`
	AvatarUrl  *string     `json:"avatarUrl" db:"avatar_url"`
	FirstName  *string     `json:"firstName" db:"first_name"`
	LastName   *string     `json:"lastName" db:"last_name"`
	Categories *string     `json:"categories" db:"categories"`
	Tags       []Tag       `json:"tags" db:"tags"`
	Medias     []PostMedia `json:"medias" db:"medias"`
}

type PostMedia struct {
	Id        int        `json:"id" db:"id"`
	PostId    string     `json:"postId" db:"post_id"`
	MediaUrl  string     `json:"mediaUrl" db:"media_url"`
	Type      string     `json:"type" db:"type"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type CreatePostPayload struct {
	Data        CreatePost              `json:"data" form:"data" binding:"required"`
	Images      []*multipart.FileHeader `json:"images" form:"images"`
	TagIds      []string                `json:"tag_ids" form:"tag_ids"`
	CategoryIds []string                `json:"category_ids" form:"category_ids"`
}

type CreatePost struct {
	UserId          string
	Title           string   `json:"title"`
	Price           float32  `json:"price"`
	Rate            string   `json:"rate"`
	Description     *string  `json:"description"`
	PickupLatitude  *float64 `json:"pickupLatitude"`
	PickupLongitude *float64 `json:"pickupLongitude"`
}

type UpdatePostPayload struct {
	Data   UpdatePost              `json:"data" form:"data" binding:"required"`
	Images []*multipart.FileHeader `json:"images" form:"images"`
	TagIds []string                `json:"tag_ids" form:"tag_ids"`
}

type UpdatePost struct {
	Title           *string  `json:"title"`
	Price           *float32 `json:"price"`
	Rate            *string  `json:"rate"`
	Description     *string  `json:"description"`
	PickupLatitude  *float64 `json:"pickupLatitude"`
	PickupLongitude *float64 `json:"pickupLongitude"`
}

type CreatePostTag struct {
	PostId string
	TagId  string
}

type CreatePostMedia struct {
	PostId   string
	MediaUrl string
	Type     string
}

type CreatePostCategory struct {
	PostId     string
	CategoryId string
}

func (r *Repository) GetPosts(ctx context.Context, params url.Values) (posts []Post, err error) {
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
		"a.title",
		"a.price",
		"a.rate",
		"a.pickup_latitude",
		"a.pickup_longitude",
		"a.created_at",
		"b.email",
		"b.avatar_url",
		`string_agg(DISTINCT d. "value", ',') AS categories`,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(cols...).From("post a").
		Join(`"user" b ON a.user_id = b.id`).
		Join("post_category c ON a.id = c.post_id").
		Join("category d ON c.category_id = d.id").
		LeftJoin("post_tag e ON a.id = e.post_id").
		LeftJoin("tag f ON e.tag_id = f.id")

	if categories := params.Get("cats"); categories != "" {
		category := strings.Split(categories, ",")[0]
		if path, exists := CATEGORY_PATH[category]; exists {
			psql = psql.Where("d.path <@ ?", path)
		}
	}

	if tags := params.Get("tags"); tags != "" {
		psql = psql.Where(sq.Eq{"f.value": strings.Split(tags, ",")})
	}

	if userId := params.Get("uid"); userId != "" {
		psql = psql.Where(sq.Eq{"a.user_id": userId})
	}

	offset := 0
	if p := params.Get("p"); p != "" {
		offset, err = strconv.Atoi(p)
		if err != nil {
			offset = 0
		}
	}

	limit := DEFAULT_LIMIT
	if l := params.Get("l"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			limit = DEFAULT_LIMIT
		}
	}

	sqlStmt, sqlArgs, err := psql.Offset(uint64(offset)*uint64(limit)).
		Limit(uint64(limit)).
		GroupBy("a.id", "b.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if err := pgxscan.ScanAll(&posts, rows); err != nil {
		return nil, fmt.Errorf("failed to scan rows | %w", err)
	}

	for idx := range posts {
		if err := r.setPostMedias(ctx, &posts[idx]); err != nil {
			return nil, fmt.Errorf("failed to set post medias | %w", err)
		}
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
		"b.first_name",
		"b.last_name",
		`string_agg(DISTINCT d. "value", ',') AS categories`,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From("post a").
		Join(`"user" b ON a.user_id = b.id`).
		Join("post_category c ON a.id = c.post_id").
		Join("category d ON c.category_id = d.id").
		Where(sq.Eq{"a.id": id}).
		GroupBy("a.id", "b.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s | %w", sqlStmt, err)
	}

	var post Post
	{
		err := pgxscan.Get(ctx, r.db, &post, sqlStmt, sqlArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute: %s | %w", sqlStmt, err)
		}
	}

	if err := r.setPostMedias(ctx, &post); err != nil {
		return nil, fmt.Errorf("failed to set post medias | %w", err)
	}

	if err := r.setPostTags(ctx, &post); err != nil {
		return nil, fmt.Errorf("failed to set post tags | %w", err)
	}

	return &post, nil
}

func (r *Repository) CreatePost(ctx context.Context, payload CreatePost) (post *Post, err error) {
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
		"title",
		"price",
		"rate",
		"description",
		"pickup_latitude",
		"pickup_longitude",
	}

	vals := []interface{}{
		payload.UserId,
		payload.Title,
		payload.Price,
		payload.Rate,
		payload.Description,
		payload.PickupLatitude,
		payload.PickupLongitude,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlStmt, sqlArgs, err := psql.Insert("post").
		Columns(cols...).
		Values(vals...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	var newPost Post
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(&newPost.Id); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return &newPost, nil
}

func (r *Repository) UpdatePost(ctx context.Context, id string, payload UpdatePost) (post *Post, err error) {
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

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("post").
		Where(sq.Eq{"id": id})

	if payload.Title != nil {
		psql = psql.Set("title", payload.Title)
	}
	if payload.Price != nil {
		psql = psql.Set("price", payload.Price)
	}
	if payload.Rate != nil {
		psql = psql.Set("rate", payload.Rate)
	}
	if payload.Description != nil {
		psql = psql.Set("description", payload.Description)
	}
	if payload.PickupLatitude != nil {
		psql = psql.Set("pickup_latitude", payload.PickupLatitude)
	}
	if payload.PickupLongitude != nil {
		psql = psql.Set("pickup_longitude", payload.PickupLongitude)
	}

	sqlStmt, sqlArgs, err := psql.Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	var updatedPost Post
	if err := tx.QueryRow(ctx, sqlStmt, sqlArgs...).Scan(&updatedPost.Id); err != nil {
		return nil, fmt.Errorf("failed to execute: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return &updatedPost, nil
}

func (r *Repository) CreatePostTags(ctx context.Context, tags []CreatePostTag) (err error) {
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
		"tag_id",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("post_tag").
		Columns(cols...)
	for idx := range tags {
		psql = psql.Values(
			tags[idx].PostId,
			tags[idx].TagId,
		)
	}

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if _, err = tx.Exec(ctx, sqlStmt, sqlArgs...); err != nil {
		return fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return nil
}

func (r *Repository) CreatePostMedias(ctx context.Context, medias []CreatePostMedia) (err error) {
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
		"media_url",
		"type",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("post_media").
		Columns(cols...)
	for idx := range medias {
		psql = psql.Values(
			medias[idx].PostId,
			medias[idx].MediaUrl,
			medias[idx].Type,
		)
	}

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if _, err = tx.Exec(ctx, sqlStmt, sqlArgs...); err != nil {
		return fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return nil
}

func (r *Repository) CreatePostCategories(ctx context.Context, categories []CreatePostCategory) (err error) {
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
		"category_id",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("post_category").
		Columns(cols...)
	for idx := range categories {
		psql = psql.Values(
			categories[idx].PostId,
			categories[idx].CategoryId,
		)
	}

	sqlStmt, sqlArgs, err := psql.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	if _, err = tx.Exec(ctx, sqlStmt, sqlArgs...); err != nil {
		return fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	return nil
}

func (r *Repository) setPostMedias(ctx context.Context, post *Post) (err error) {
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
		"media_url",
		"type",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From("post_media").
		Where(sq.Eq{"post_id": post.Id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	postMedias := []PostMedia{}
	if err := pgxscan.ScanAll(&postMedias, rows); err != nil {
		return fmt.Errorf("failed to scan rows | %w", err)
	}

	post.Medias = postMedias

	return nil
}

func (r *Repository) setPostTags(ctx context.Context, post *Post) (err error) {
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
		"b.id",
		"b.type",
		"b.value",
		"b.label",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStmt, sqlArgs, err := psql.Select(cols...).
		From("post_tag a").
		Join(`"tag" b ON a.tag_id = b.id`).
		Where(sq.Eq{"a.post_id": post.Id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	rows, err := tx.Query(ctx, sqlStmt, sqlArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %s args: %v | %w", sqlStmt, sqlArgs, err)
	}

	postTags := []Tag{}
	if err := pgxscan.ScanAll(&postTags, rows); err != nil {
		return fmt.Errorf("failed to scan rows | %w", err)
	}

	post.Tags = postTags

	return nil
}

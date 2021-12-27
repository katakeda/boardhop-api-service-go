package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const SURFBOARD_PATH = "root.1"
const SNOWBOARD_PATH = "root.2"
const PER_PAGE_MAX = 25

type App struct {
	router *gin.Engine
	db     *pgxpool.Pool
}

type Post struct {
	Id              int       `json:"id"`
	UserId          int       `json:"userId"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Price           float32   `json:"price"`
	Rate            string    `json:"rate"`
	PickupLatitude  float64   `json:"pickupLatitude"`
	PickupLongitude float64   `json:"pickupLongitude"`
	CreatedAt       time.Time `json:"createdAt"`
	Username        string    `json:"username"`
	AvatarUrl       string    `json:"avatarUrl"`
	Categories      string    `json:"categories"`
	Tags            string    `json:"tags"`
}

func (app *App) Initialize() {
	gin.Logger()

	app.router = gin.Default()
	app.router.GET("/posts", app.getPosts)
	app.router.GET("/posts/:id", app.getPost)
	app.router.POST("/posts", app.createPost)

	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		log.Fatalln("Failed to connect with DB", err.Error())
		os.Exit(1)
	}

	app.db = db
}

func (app *App) Run() {
	defer app.db.Close()
	err := app.router.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalln("Failed to run app", err.Error())
		os.Exit(1)
	}
}

func (app *App) getPosts(c *gin.Context) {
	params := c.Request.URL.Query()

	var rootPath string
	{
		if params.Get("type") == "snowboard" {
			rootPath = SNOWBOARD_PATH
		} else {
			rootPath = SURFBOARD_PATH
		}
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
		Join("post_tag e ON a.id = e.post_id").
		Join("tag f ON e.tag_id = f.id").
		Join("tag_type g ON f.type_id = g.id").
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
			log.Fatalln("Failed to build query:", sqlStmt, err)
			c.IndentedJSON(http.StatusInternalServerError, "Error")
		}
	}

	var posts []*Post
	err := pgxscan.Select(context.Background(), app.db, &posts, sqlStmt, sqlArgs...)
	if err != nil {
		log.Fatalln("Failed to execute:", sqlStmt, err)
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	if len(posts) <= 0 {
		log.Println("No posts found")
		c.IndentedJSON(http.StatusNotFound, "No posts found")
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func (app *App) getPost(c *gin.Context) {
	id := c.Param("id")

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

	sqlStmt, sqlArgs, err := psql.Select(cols...).From("post a").
		Join(`"user" b ON a.user_id = b.id`).
		Join("post_category c ON a.id = c.post_id").
		Join("category d ON c.category_id = d.id").
		Join("post_tag e ON a.id = e.post_id").
		Join("tag f ON e.tag_id = f.id").
		Join("tag_type g ON f.type_id = g.id").
		Where(sq.Eq{"a.id": id}).
		GroupBy("a.id", "b.id").
		ToSql()
	if err != nil {
		log.Fatalln("Failed to build query:", sqlStmt, err)
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	var post Post
	{
		err := pgxscan.Get(context.Background(), app.db, &post, sqlStmt, sqlArgs...)
		if err != nil {
			log.Println("Failed to execute:", sqlStmt, err)
			c.IndentedJSON(http.StatusNotFound, nil)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, post)
}

func (app *App) createPost(c *gin.Context) {
	log.Println("createPost")
	c.IndentedJSON(http.StatusOK, "OK")
}

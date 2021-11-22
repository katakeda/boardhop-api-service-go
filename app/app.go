package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

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
	app.router.GET("/post/:id", app.getPost)
	app.router.POST("/post", app.createPost)

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
	log.Println("getPosts")
	sqlStmt := `
	SELECT
		a.id,
		a.user_id,
		a.title,
		a.description,
		a.price,
		a.rate,
		a.pickup_latitude,
		a.pickup_longitude,
		a.created_at,
		b.username,
		b.avatar_url,
		string_agg(DISTINCT d. "name", ',') AS categories,
		string_agg(DISTINCT f. "value", ',') AS tags
	FROM
		post a
		JOIN "user" b ON a.user_id = b.id
		JOIN post_category c ON a.id = c.post_id
		JOIN category d ON c.category_id = d.id
		JOIN post_tag e ON a.id = e.post_id
		JOIN tag f ON e.tag_id = f.id
		JOIN tag_type g ON f.type_id = g.id
	WHERE
		d.path <@ 'root.1'
	GROUP BY
		a.id, b.id;
	`

	var posts []*Post
	err := pgxscan.Select(context.Background(), app.db, &posts, sqlStmt)
	if err != nil {
		log.Fatalln("Failed to execute:", sqlStmt, err)
		c.IndentedJSON(http.StatusOK, "Error")
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func (app *App) getPost(c *gin.Context) {
	log.Println("getPost")
	c.IndentedJSON(http.StatusOK, c.Param("id"))
}

func (app *App) createPost(c *gin.Context) {
	log.Println("createPost")
	c.IndentedJSON(http.StatusOK, "OK")
}

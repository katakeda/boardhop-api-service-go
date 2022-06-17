package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/katakeda/boardhop-api-service-go/repositories"
	"github.com/katakeda/boardhop-api-service-go/services"
)

type App struct {
	router *gin.Engine
}

func (app *App) Initialize() {
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Fatalln("Failed to connect with DB", err)
	}

	repo, err := repositories.NewRepository(db)
	if err != nil {
		log.Fatalln("Failed to initialize repository", err)
	}

	svc, err := services.NewService(repo)
	if err != nil {
		log.Fatalln("Failed to initialize service", err)
	}

	app.router = gin.Default()
	app.router.GET("/posts", svc.GetPosts)
	app.router.GET("/posts/:id", svc.GetPost)
	app.router.GET("/tags", svc.GetTags)
	app.router.GET("/categories", svc.GetCategories)
	app.router.GET("/user", AuthRequired(), svc.GetUser)
	app.router.GET("/orders", AuthRequired(), svc.GetOrders)
	app.router.GET("/orders/:id", AuthRequired(), svc.GetOrder)

	app.router.POST("/user/signup", svc.UserSignup)
	app.router.POST("/user/login", svc.UserLogin)
	app.router.POST("/posts", AuthRequired(), svc.CreatePost)
	app.router.POST("/orders", AuthRequired(), svc.CreateOrder)
	app.router.POST("/messages", AuthRequired(), svc.CreateMessage)
}

func (app *App) Run() {
	err := app.router.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalln("Failed to run app", err)
	}
}

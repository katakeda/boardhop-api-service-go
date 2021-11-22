package main

import (
	"github.com/joho/godotenv"
	"github.com/katakeda/boardhop-api-service-go/app"
)

func main() {
	godotenv.Load()

	app := app.App{}
	app.Initialize()
	app.Run()
}

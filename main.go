package main

import (
	"github.com/katakeda/boardhop-api-service-go/app"
)

func main() {
	app := app.App{}
	app.Initialize()
	app.Run()
}

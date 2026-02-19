package main

import (
	"context"
	"fmt"

	"github.com/viktare/go-shortener/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start server", err)
	}
}


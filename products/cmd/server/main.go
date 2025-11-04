package main

import (
	"fmt"
	"os"

	app2 "github.com/jva44ka/ozon-simulator-go/internal/app"
)

func main() {
	fmt.Println("app starting")

	app, err := app2.NewApp(os.Getenv("ROUTE_256_WS_1"))
	if err != nil {
		panic(err)
	}

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}

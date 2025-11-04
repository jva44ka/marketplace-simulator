package main

import (
	"fmt"
	app2 "gitlab.ozon.dev/16/students/week-1-workshop/internal/app"
	"os"
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

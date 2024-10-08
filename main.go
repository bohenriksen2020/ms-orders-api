package main

import (
	"fmt"
	"context"
	"github.com/bohenriksen2020/ms-orders-api/application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}

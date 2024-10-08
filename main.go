package main

import (
	"fmt"
	"context"
	"os/signal"
	"github.com/bohenriksen2020/ms-orders-api/application"
)

func main() {
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel() // cancel context when main function ends

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}


	
}

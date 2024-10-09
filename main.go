package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bohenriksen2020/ms-orders-api/application"
)

func main() {
	// Load the config
	config := application.LoadConfig()

	// Initialize the application (App) with just the config
	app := application.New(config)

	// Handle context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel() // cancel context when main function ends

	// Start the application
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}

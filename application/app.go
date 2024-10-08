package application

import (
	"context"
	"fmt"
	"net/http"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb	*redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{ Addr: "localhost:6379", DB: 0 }),
	}
	return app
}


defer func() {

	if err := a.rdb.Close(); err != nil {
		fmt.Println("failed to close redis connection: ", err)
	}
}()

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}
	
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}
	fmt.Println("Server is starting")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen and serve: %w", err)
		}
		close(ch)
	}()

	select {
		err := <-ch:
			return err
		case <-ctx.Done():
			timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			
			return server.Shutdown(timeout)
	}
	
	return nil
}


package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bohenriksen2020/ms-orders-api/repository/order"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	repo   order.Repo // Use the Repo interface
	config Config
}

// New creates a new App instance and chooses the correct repository (Redis or Postgres)
func New(config Config) *App {
	var repo order.Repo

	// Initialize the correct repository based on the config
	if config.DatabaseType == "postgres" {
		// Initialize PostgresRepo
		db, err := sql.Open("postgres", config.PostgresDSN)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %v", err))
		}
		repo = order.NewPostgresRepo(db)
		fmt.Println("Using Postgres as the database")
	} else {
		// Initialize RedisRepo
		rdb := redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		})
		repo = order.NewRedisRepo(rdb)
		fmt.Println("Using Redis as the database")
	}

	app := &App{
		repo:   repo,  // Use the selected Repo implementation (Redis or Postgres)
		config: config,
	}
	app.loadRoutes() // Define your routes here
	return app
}

func (a *App) Start(ctx context.Context) error {
	// If the repo is a RedisRepo, handle Redis-specific shutdown
	if redisRepo, ok := a.repo.(*order.RedisRepo); ok {
		defer func() {
			if err := redisRepo.Close(); err != nil {
				fmt.Println("failed to close Redis connection: ", err)
			}
		}()

		// Ping Redis to ensure connection
		err := redisRepo.Ping(ctx).Err()
		if err != nil {
			return fmt.Errorf("failed to ping Redis: %w", err)
		}
	}

	// If the repo is a PostgresRepo, handle Postgres-specific shutdown
	if postgresRepo, ok := a.repo.(*order.PostgresRepo); ok {
		defer func() {
			if err := postgresRepo.Close(); err != nil {
				fmt.Println("failed to close PostgreSQL connection: ", err)
			}
		}()
	}

	// Setup and start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	fmt.Println("Server is starting")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to listen and serve: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bohenriksen2020/ms-orders-api/handler"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Load order-related routes
	router.Route("/orders", a.loadOrderRoutes)

	// Assign the router to the App's router field
	a.router = router
}

func (a *App) loadOrderRoutes(route chi.Router) {
	// Initialize the order handler using the repository provided in a.repo
	orderHandler := &handler.Order{
		Repo: a.repo, // Use the repository passed to the App (could be Redis or Postgres)
	}

	// Define routes for order operations
	route.Post("/", orderHandler.Create)
	route.Get("/", orderHandler.List)
	route.Get("/{id}", orderHandler.GetByID)
	route.Put("/{id}", orderHandler.UpdateByID)
	route.Delete("/{id}", orderHandler.DeleteByID)
}

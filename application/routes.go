package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/bohenriksen2020/ms-orders-api/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRoutes) 

	return router
}


func loadOrderRoutes(route chi.Router) {
	orderHandler := &handler.Order{}

	route.Post("/", orderHandler.Create)
	route.Get("/", orderHandler.List)
	route.Get("/{id}", orderHandler.GetByID)
	route.Put("/{id}", orderHandler.UpdateByID)
	route.Delete("/{id}", orderHandler.DeleteByID)

}
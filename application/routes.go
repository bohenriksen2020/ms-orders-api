package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/bohenriksen2020/ms-orders-api/handler"
	"github.com/bohenriksen2020/ms-orders-api/repository/order"


)

func (a *App) loadRoutes()  {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", a.loadOrderRoutes) 

	a.router = router
}


func (a *App) loadOrderRoutes(route chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo {
			Client: a.rdb,
		},
	}

	route.Post("/", orderHandler.Create)
	route.Get("/", orderHandler.List)
	route.Get("/{id}", orderHandler.GetByID)
	route.Put("/{id}", orderHandler.UpdateByID)
	route.Delete("/{id}", orderHandler.DeleteByID)

}
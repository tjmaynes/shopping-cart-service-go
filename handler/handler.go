package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	handlers "github.com/tjmaynes/shopping-cart-service-go/handler/http"
)

// Initialize ..
func Initialize(itemHandler *handlers.ItemHandler, healthCheckHandler *handlers.HealthCheckHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	router.Route("/", func(rt chi.Router) {
		rt.Mount("/items", addItemRouter(itemHandler))
		rt.Get("/health", healthCheckHandler.GetHealthCheckHandler)
	})

	return router
}

func addItemRouter(itemHandler *handlers.ItemHandler) http.Handler {
	router := chi.NewRouter()

	router.Get("/", itemHandler.GetItems)
	router.Get("/{id}", itemHandler.GetItemByID)
	router.Post("/", itemHandler.AddItem)
	router.Put("/{id}", itemHandler.UpdateItem)
	router.Delete("/{id}", itemHandler.RemoveItem)

	return router
}

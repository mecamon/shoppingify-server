package items_summary

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	itemsHandler := chi.NewRouter()
	itemsHandler.Use(middlewares.TokenValidation)
	itemsHandler.Get("/{year}", handler.GetByMonth)
	itemsHandler.Get("/", handler.GetByYear)
	return itemsHandler
}

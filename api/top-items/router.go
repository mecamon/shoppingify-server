package top_items

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	topItemsRouter := chi.NewRouter()
	topItemsRouter.Use(middlewares.TokenValidation)
	topItemsRouter.Get("/", handler.GetTop)
	return topItemsRouter
}

package top_categories

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	topCategoriesRouter := chi.NewRouter()
	topCategoriesRouter.Use(middlewares.TokenValidation)
	topCategoriesRouter.Get("/", handler.GetTop)
	return topCategoriesRouter
}

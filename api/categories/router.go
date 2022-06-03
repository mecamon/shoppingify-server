package categories

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	categoriesRouter := chi.NewRouter()
	categoriesRouter.Use(middlewares.TokenValidation)
	categoriesRouter.Post("/", handler.Create)
	categoriesRouter.Get("/", handler.GetAllByName)
	return categoriesRouter
}

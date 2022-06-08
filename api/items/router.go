package items

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	itemsHandler := chi.NewRouter()
	itemsHandler.Use(middlewares.TokenValidation)
	itemsHandler.Post("/", handler.Create)
	itemsHandler.Get("/", handler.GetByCategoryGroups)
	return itemsHandler
}

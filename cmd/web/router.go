package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/auth"
	"github.com/mecamon/shoppingify-server/api/categories"
	"github.com/mecamon/shoppingify-server/api/items"
	"github.com/mecamon/shoppingify-server/config"
	"net/http"
)

func makeRouter() http.Handler {
	r := chi.NewRouter()
	appConfig := config.Get()

	auth.InitHandler(appConfig)
	r.Mount("/api/auth", auth.Routes())
	items.InitHandler(appConfig)
	r.Mount("/api/items", items.Routes())
	categories.InitHandler(appConfig)
	r.Mount("/api/categories", categories.Routes())
	return r
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/auth"
	"github.com/mecamon/shoppingify-server/config"
	"net/http"
)

func makeRouter() http.Handler {
	r := chi.NewRouter()

	auth.InitHandler(config.Get())
	r.Mount("/api/auth", auth.Routes())
	return r
}

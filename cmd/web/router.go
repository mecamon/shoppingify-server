package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mecamon/shoppingify-server/api/auth"
	"github.com/mecamon/shoppingify-server/api/categories"
	"github.com/mecamon/shoppingify-server/api/items"
	"github.com/mecamon/shoppingify-server/api/lists"
	top_categories "github.com/mecamon/shoppingify-server/api/top-categories"
	top_items "github.com/mecamon/shoppingify-server/api/top-items"
	"github.com/mecamon/shoppingify-server/config"
	_ "github.com/mecamon/shoppingify-server/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func makeRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // TODO: Change it if you want to allow a specific domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Accept-Language", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	appConfig := config.Get()

	auth.InitHandler(appConfig)
	r.Mount("/api/auth", auth.Routes())
	items.InitHandler(appConfig)
	r.Mount("/api/items", items.Routes())
	categories.InitHandler(appConfig)
	r.Mount("/api/categories", categories.Routes())
	lists.InitHandler(appConfig)
	r.Mount("/api/lists", lists.Routes())
	top_categories.InitHandler(appConfig)
	r.Mount("/api/top-categories", top_categories.Routes())
	top_items.InitHandler(appConfig)
	r.Mount("/api/top-items", top_items.Routes())

	// Documenting the api
	var domain string
	if appConfig.IsProd {
		domain = "https://shoppingify-be.onrender.com"
	} else {
		domain = "http://localhost:8080"
	}
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(domain+"/swagger/doc.json")))

	return r
}

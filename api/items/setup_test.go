//go:build integration
// +build integration

package items

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"log"
	"net/http"
	"os"
	"testing"
)

var Router http.Handler

func TestMain(m *testing.M) {
	conn := setup()
	code := m.Run()
	shutdown(conn)
	os.Exit(code)
}

func setup() *sql.DB {
	config.Set()
	appConfig := config.Get()

	err := appi18n.InitLocales()
	if err != nil {
		panic(err.Error())
	}
	
	conn, err := db.InitDB(appConfig)
	if err != nil {
		log.Println(err.Error())
	}
	repositories.InitRepos(conn, appConfig)
	itemsHandler := InitHandler(appConfig)
	r := chi.NewRouter()
	r.Use(middlewares.TokenValidation)
	r.Post("/api/items/", itemsHandler.Create)
	r.Get("/api/items/", itemsHandler.GetByCategoryGroups)
	r.Get("/api/items/{id}", itemsHandler.GetDetailsByID)
	Router = r

	return conn
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

package main

import (
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/services/storage"
	"log"
	"net/http"
	"os"
)

// @title shoppingify-server APIs
// @version 1.0
// @description shopping list site APIs
// @BasePath /
func main() {
	run()
}

func run() {
	config.Set()
	conf := config.Get()

	err := appi18n.InitLocales()
	if err != nil {
		panic(err)
	}

	conn, err := db.InitDB(conf)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	repositories.InitRepos(conn, conf)
	_, err = storage.InitStorage()
	if err != nil {
		conf.Loggers.Error.Println(err.Error())
	}

	router := makeRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("running on port %s...", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("Could not start server", router)
	}
}

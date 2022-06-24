package main

import (
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/services/storage"
	"log"
	"net/http"
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

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Could not start server", router)
	}
}

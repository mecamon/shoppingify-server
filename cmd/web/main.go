package main

import (
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	"log"
	"net/http"
)

func main() {
	run()
}

func run() {
	config.Set()
	conf := config.Get()

	database, err := db.InitDB(conf)
	if err != nil {
		panic(err.Error())
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}
	log.Println("DB pinged!!!")

	router := makeRouter()

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Could not start server", router)
	}
}

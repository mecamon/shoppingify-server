package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
)

func main() {
	log.Println("It is working...")
	config.Set()
	conf := config.Get()

	db, err := db.InitDB(conf)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	log.Println("DB pinged!!!")

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome!")
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
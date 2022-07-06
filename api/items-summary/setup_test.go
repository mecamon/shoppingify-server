//go:build integration
// +build integration

package items_summary

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	"github.com/mecamon/shoppingify-server/db"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	Router        http.Handler
	HandlerUserID int64
	Token         string
)

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
	summaryHandler := InitHandler(appConfig)
	r := chi.NewRouter()
	r.Use(middlewares.TokenValidation)
	r.Get("/api/summary/{year}", summaryHandler.GetByMonth)
	r.Get("/api/summary", summaryHandler.GetByYear)
	Router = r

	createUser()

	return conn
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

func createUser() {
	//Registering new user
	validPass := "Validpass09"
	passwordHashed, err := utils.GenerateHash(validPass)
	if err != nil {
		log.Println(err.Error())
	}
	HandlerUserID, err = repositories.Main.AuthRepoImpl.Register(models.User{
		Name:      "Get by month H",
		Lastname:  "lastname",
		Email:     "getbymonth@handler.com",
		Password:  passwordHashed,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})

	Token, err = json_web_token.Generate(HandlerUserID, "")
	if err != nil {
		log.Println(err.Error())
	}
}

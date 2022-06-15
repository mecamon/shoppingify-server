//go:build integration
// +build integration

package top_items

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	"github.com/mecamon/shoppingify-server/db"
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
	userIdForTest int64
	tokenForTests string
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
	conn, err := db.InitDB(appConfig)
	if err != nil {
		log.Println(err.Error())
	}
	repositories.InitRepos(conn, appConfig)
	topItemsHandler := InitHandler(appConfig)
	r := chi.NewRouter()
	r.Use(middlewares.TokenValidation)
	r.Get("/api/top-items", topItemsHandler.GetTop)

	Router = r
	generateTokenForTests()
	return conn
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

func generateTokenForTests() {
	user := models.User{
		Name:      "user test topitem 1",
		Lastname:  "lastname",
		Email:     "usertest@topitem.com",
		Password:  "ValidPassword123",
		IsActive:  true,
		IsVisitor: false,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, _ := utils.GenerateHash(user.Password)
	user.Password = hashedPass

	insertedUserID, _ := repositories.Main.AuthRepoImpl.Register(user)
	userIdForTest = insertedUserID
	tokenForTests, _ = json_web_token.Generate(insertedUserID, "")
}

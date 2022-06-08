//go:build integration
// +build integration

package categories

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
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var Router http.Handler
var userToken string

func TestMain(m *testing.M) {
	conn := setup()
	code := m.Run()
	shutdown(conn)
	os.Exit(code)
}

func setup() *sql.DB {
	config.Set()
	appConfig := config.Get()
	db, err := db.InitDB(appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	repositories.InitRepos(db, appConfig)
	categoriesHandler := InitHandler(appConfig)
	r := chi.NewRouter()
	r.Use(middlewares.TokenValidation)
	r.Post("/api/categories/", categoriesHandler.Create)
	r.Get("/api/categories/", categoriesHandler.GetAll)
	r.Get("/api/categories/by-name", categoriesHandler.GetAllByName)
	Router = r

	seedingData() //Seeding data for the handlers tests

	return db
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

func seedingData() {
	newUIID := uuid.NewString()
	user := models.User{
		Name:      "User Cat Handler",
		Lastname:  "User Cat Handler LM",
		Email:     "cat@handler.com",
		Password:  "CatHandlerPass123",
		IsActive:  true,
		IsVisitor: false,
		LoginCode: newUIID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	repo := repositories.Main
	ID, err := repo.AuthRepoImpl.Register(user)
	if err != nil {
		log.Fatal(err.Error())
	}

	userToken, err = json_web_token.Generate(ID, user.Email)
	if err != nil {
		log.Fatal(err.Error())
	}

	cat1 := models.Category{
		Name:      "Meats",
		UserID:    ID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	cat2 := models.Category{
		Name:      "Snow",
		UserID:    ID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	repo.CategoriesRepoImpl.RegisterCategory(cat1)
	repo.CategoriesRepoImpl.RegisterCategory(cat2)
}

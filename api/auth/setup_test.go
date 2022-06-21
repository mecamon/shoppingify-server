//go:build integration
// +build integration

package auth

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	"github.com/mecamon/shoppingify-server/models"
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

	conn, err := db.InitDB(appConfig)
	if err != nil {
		panic(err.Error())
	}

	err = conn.Ping()
	if err != nil {
		panic(err.Error())
	}

	repositories.InitRepos(conn, appConfig)
	authHandler := InitHandler(appConfig)
	seedDataForIntTests(authHandler.repos)

	r := chi.NewRouter()
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/visitor-register", authHandler.VisitorRegister)
	Router = r
	return conn
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

func seedDataForIntTests(repos repositories.MainRepo) {
	seedUser := models.UserDTO{
		Name:     LoginUserData.Name,
		Lastname: LoginUserData.Lastname,
		Email:    LoginUserData.Email,
		Password: LoginUserData.Password,
	}
	completeSeedUser, _ := completeUserInformation(seedUser)
	repos.AuthRepoImpl.Register(completeSeedUser)
}

package repositories

import (
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
)

var Main MainRepo

type MainRepo struct {
	AuthRepoImpl AuthRepo
}

type AuthRepo interface {
	Register(user models.User) error
	SearchUserByEmail(email string) (models.User, error)
	CheckUserPassword(email, password string) (bool, error)
}

func InitRepos(conn *sql.DB, app *config.App) {
	authRepo := initAuthRepo(conn, app)
	Main = MainRepo{AuthRepoImpl: authRepo}
}

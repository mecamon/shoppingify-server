package repositories

import (
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
)

var Main MainRepo

type MainRepo struct {
	AuthRepoImpl       AuthRepo
	CategoriesRepoImpl CategoriesRepo
}

func InitRepos(conn *sql.DB, app *config.App) {
	authRepo := initAuthRepo(conn, app)
	categoriesRepo := initCategoriesRepo(conn, app)
	Main = MainRepo{
		AuthRepoImpl:       authRepo,
		CategoriesRepoImpl: categoriesRepo,
	}
}

type AuthRepo interface {
	Register(user models.User) (int64, error)
	SearchUserByEmail(email string) (models.User, error)
	CheckUserPassword(email, password string) (bool, error)
}

type CategoriesRepo interface {
	RegisterCategory(category models.Category) (int64, error)
	SearchCategoryByName(q string, skip, take int) ([]models.CategoryDTO, error)
}

type ItemsRepo interface {
	//RegisterItem(item models.Item) (int64, error)
	//GetItems(take, skip int64) ([]models.Item, error)
}

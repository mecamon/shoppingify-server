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
	ItemsRepoIpml      ItemsRepo
	ListsRepoImpl      ListsRepo
}

func InitRepos(conn *sql.DB, app *config.App) {
	authRepo := initAuthRepo(conn, app)
	categoriesRepo := initCategoriesRepo(conn, app)
	itemsRepo := initItemsRepo(conn, app)
	listsRepo := initListsRepo(conn, app)

	Main = MainRepo{
		AuthRepoImpl:       authRepo,
		CategoriesRepoImpl: categoriesRepo,
		ItemsRepoIpml:      itemsRepo,
		ListsRepoImpl:      listsRepo,
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
	GetAll(take, skip int) ([]models.CategoryDTO, error)
	GetAllWithItemName(q string, take, skip int) ([]models.CategoryDTO, error)
	Count(filter ...string) (int64, error)
}

type ItemsRepo interface {
	Register(item models.Item) (int64, error)
	GetAll(take, skip int) ([]models.ItemDTO, error)
	GetDetails(id int64) (models.ItemDTO, error)
	GetByID(itemID int64) (models.ItemDTO, error)
	GetAllByCategoryID(categoryId int64) ([]models.ItemDTO, error)
	Count() (int64, error)
}

type ListsRepo interface {
	Create(list models.List) (int64, error)
	GetActive() (models.ListDTO, error)
	AddItemToList(item models.SelectedItem) (int64, error)
	IsListActive(listID int64) (bool, error)
	DeleteItemFromList(itemSelID int64) error
	CompleteItemSelected(itemSelID int64) error
	CancelActive() error
	CompleteActive() error
}

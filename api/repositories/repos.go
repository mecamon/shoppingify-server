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
	TopCategoriesImpl  TopCategoriesRepo
	TopItemsImpl       TopItemsRepo
}

func InitRepos(conn *sql.DB, app *config.App) {
	authRepo := initAuthRepo(conn, app)
	categoriesRepo := initCategoriesRepo(conn, app)
	itemsRepo := initItemsRepo(conn, app)
	listsRepo := initListsRepo(conn, app)
	topCategories := initTopCategoriesRepo(conn, app)
	topItemsRepo := initTopItemsRepo(conn, app)

	Main = MainRepo{
		AuthRepoImpl:       authRepo,
		CategoriesRepoImpl: categoriesRepo,
		ItemsRepoIpml:      itemsRepo,
		ListsRepoImpl:      listsRepo,
		TopCategoriesImpl:  topCategories,
		TopItemsImpl:       topItemsRepo,
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
	UpdateActiveListName(userID int64, name string) error
	GetActive(userID int64) (models.ListDTO, error)
	AddItemToList(item models.SelectedItem) (int64, error)
	IsListActive(listID int64) (bool, error)
	DeleteItemFromList(itemSelID int64) error
	CompleteItemSelected(itemSelID int64) error
	UpdateItemsSelected(items []models.UpdateSelItemDTO) error
	CancelActive(userID int64) error
	CompleteActive(userID int64) error
}

type TopCategoriesRepo interface {
	Add(userID, categoryID int64) (int64, error)
	Update(userID, categoryID int64) error
	GetTop(userID int64, take int) ([]models.TopCategoryDTO, error)
	GetAll(userID int64) ([]models.TopCategoryDTO, error)
}

type TopItemsRepo interface {
	Add(userID, itemID int64) (int64, error)
	Update(userID, itemID int64) error
	GetTop(userID int64, take int) ([]models.TopItemDTO, error)
	GetAll(userID int64) ([]models.TopItemDTO, error)
}

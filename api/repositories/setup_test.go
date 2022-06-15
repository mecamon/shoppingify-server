//go:build integration
// +build integration

package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"os"
	"testing"
	"time"
)

var (
	authRepo               AuthRepo
	categoriesRepo         CategoriesRepo
	itemsRepo              ItemsRepo
	listsRepo              ListsRepo
	topCategoriesRepo      TopCategoriesRepo
	topItemsRepo           TopItemsRepo
	userIdForTestRepos     int64
	categoryIDForTestRepos int64
	conn                   *sql.DB
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	config.Set()
	conf := config.Get()
	conn, _ = db.InitDB(conf)
	authRepo = initAuthRepo(conn, conf)
	categoriesRepo = initCategoriesRepo(conn, conf)
	itemsRepo = initItemsRepo(conn, conf)
	listsRepo = initListsRepo(conn, conf)
	topCategoriesRepo = initTopCategoriesRepo(conn, conf)
	topItemsRepo = initTopItemsRepo(conn, conf)
	creatingUserForTestRepos()
}

func shutdown() {
	err := conn.Close()
	if err != nil {
		log.Println("error shutting down db connection: ", err.Error())
	}
}

func creatingUserForTestRepos() {
	hashedPass, _ := utils.GenerateHash("ValidPass123")
	userIdForTestRepos, _ = authRepo.Register(models.User{
		Name:      "User for test repos",
		Lastname:  "lastname",
		Email:     "usertest@repos.com",
		Password:  hashedPass,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	categoryIDForTestRepos, _ = categoriesRepo.RegisterCategory(models.Category{
		Name:      "Category for test repos",
		UserID:    userIdForTestRepos,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
}

//go:build integration
// +build integration

package lists

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
	Router          http.Handler
	userIdForTest   int64
	tokenForTests   string
	insertedCatID1  int64
	insertedCatID2  int64
	insertedItemID1 int64
	insertedItemID2 int64
	insertedItemID3 int64
	insertedItemID4 int64
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
	listsHandler := InitHandler(appConfig)
	r := chi.NewRouter()
	r.Use(middlewares.TokenValidation)
	r.Post("/api/lists/create", listsHandler.Create)
	r.Get("/api/lists/active", listsHandler.GetActive)
	r.Patch("/api/lists/name", listsHandler.UpdateActiveListName)
	r.Post("/api/lists/add-item", listsHandler.AddItemToList)
	r.Patch("/api/lists/update-items", listsHandler.UpdateItemsSelected)
	r.Delete("/api/lists/selected-items/{itemID}", listsHandler.DeleteItemFromList)
	r.Put("/api/lists/selected-items", listsHandler.CompleteItemSelected)
	r.Delete("/api/lists/cancel-active", listsHandler.CancelActive)
	r.Patch("/api/lists/complete-active", listsHandler.CompleteActive)
	r.Get("/api/lists/old-lists", listsHandler.GetOldLists)
	r.Get("/api/lists/{listID}", listsHandler.GetByID)
	Router = r

	generateTokenForTests()
	seedCategories()
	seedItems()

	return conn
}

func shutdown(conn *sql.DB) {
	conn.Close()
}

func generateTokenForTests() {
	user := models.User{
		Name:      "user test lists 1",
		Lastname:  "lastname",
		Email:     "usertest@lists.com",
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

func seedCategories() {
	insertedCatID1, _ = repositories.Main.CategoriesRepoImpl.RegisterCategory(models.Category{
		Name:      "Category test lists handler 1",
		UserID:    userIdForTest,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	insertedCatID2, _ = repositories.Main.CategoriesRepoImpl.RegisterCategory(models.Category{
		Name:      "Category test lists handler 2",
		UserID:    userIdForTest,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
}

func seedItems() {
	insertedItemID1, _ = repositories.Main.ItemsRepoIpml.Register(models.Item{
		Name:       "Item test lists handler 1",
		Note:       "This is a note",
		CategoryID: insertedCatID1,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	insertedItemID2, _ = repositories.Main.ItemsRepoIpml.Register(models.Item{
		Name:       "Item test lists handler 2",
		Note:       "This is a note",
		CategoryID: insertedCatID1,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	insertedItemID3, _ = repositories.Main.ItemsRepoIpml.Register(models.Item{
		Name:       "Item test lists handler 3",
		Note:       "This is a note",
		CategoryID: insertedCatID2,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	insertedItemID4, _ = repositories.Main.ItemsRepoIpml.Register(models.Item{
		Name:       "Item test lists handler 4",
		Note:       "This is a note",
		CategoryID: insertedCatID2,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
}

//go:build integration
// +build integration

package repositories

import (
	"fmt"
	fixtures_categoryies_repo "github.com/mecamon/shoppingify-server/__test__/fixtures/repos/categoryies-repo"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
	"time"
)

func TestItemsRepoPostgres_RegisterCategory(t *testing.T) {
	user := fixtures_categoryies_repo.User
	hashedPass, _ := utils.GenerateHash(fixtures_categoryies_repo.User.Password)
	user.Password = hashedPass
	userID, _ := authRepo.Register(user)

	cat := fixtures_categoryies_repo.Cat1
	cat.UserID = userID

	_, err := categoriesRepo.RegisterCategory(cat)
	if err != nil {
		t.Error("error registering category: ", err.Error())
	}
}

func TestItemsRepoPostgres_SearchCategoryByName(t *testing.T) {
	user := fixtures_categoryies_repo.UserForSearchCat
	hashedPass, _ := utils.GenerateHash(fixtures_categoryies_repo.UserForSearchCat.Password)
	user.Password = hashedPass
	userID, _ := authRepo.Register(user)

	cat := fixtures_categoryies_repo.Cat2
	cat.UserID = userID
	categoriesRepo.RegisterCategory(cat)

	var skip = 0
	var take = 6
	q := "eat"
	_, err := categoriesRepo.SearchCategoryByName(q, take, skip)
	if err != nil {
		t.Error("error searching for categories: ", err.Error())
	}
}

func TestCategoriesRepoPostgres_GetAll(t *testing.T) {
	user := fixtures_categoryies_repo.UserForGetAll
	hashedPass, _ := utils.GenerateHash(user.Password)
	user.Password = hashedPass

	userID, _ := authRepo.Register(user)

	var catGeneric = models.Category{
		Name:      "CatGet",
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	for i := 0; i < 10; i++ {
		catGeneric.Name = fmt.Sprintf("%s-%d", catGeneric.Name, i)
		_, err := categoriesRepo.RegisterCategory(catGeneric)
		if err != nil {
			t.Error(err.Error())
		}
	}

	categories, err := categoriesRepo.GetAll(6, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if len(categories) > 6 {
		t.Errorf("expected %d categories but got %d", 6, len(categories))
	}
}

func TestCategoriesRepoPostgres_GetAllWithItemName(t *testing.T) {
	user := fixtures_categoryies_repo.UserForGetAllWithItemName
	hashedPass, _ := utils.GenerateHash(user.Password)
	user.Password = hashedPass
	userID, err := authRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	// adding two categories
	cat1 := fixtures_categoryies_repo.Cat1ForItemsName
	cat1.UserID = userID
	insertedCat1ID, err := categoriesRepo.RegisterCategory(cat1)
	if err != nil {
		t.Error(err.Error())
	}
	cat2 := fixtures_categoryies_repo.Cat1ForItemsName
	cat2.UserID = userID
	insertedCat2ID, err := categoriesRepo.RegisterCategory(cat2)
	if err != nil {
		t.Error(err.Error())
	}

	// adding two items for the previously added categories
	item1 := fixtures_categoryies_repo.Item1ForItemsName
	item1.CategoryID = insertedCat1ID
	_, err = itemsRepo.Register(item1)
	if err != nil {
		t.Error(err.Error())
	}
	item2 := fixtures_categoryies_repo.Item2ForItemsName
	item2.CategoryID = insertedCat2ID
	_, err = itemsRepo.Register(item2)
	if err != nil {
		t.Error(err.Error())
	}

	// asserting data
	categories, err := categoriesRepo.GetAllWithItemName("Item2", 8, 0)
	if err != nil {
		t.Error(err.Error())
	}
	for _, c := range categories {
		if c.ID != item2.CategoryID {
			t.Errorf("expected %d as category_id but got %d", item2.CategoryID, c.ID)
		}
	}
}

func TestCategoriesRepoPostgres_Count(t *testing.T) {
	qty, err := categoriesRepo.Count()
	if err != nil {
		t.Error(err.Error())
	}
	if qty == 0 {
		t.Error("expected more than 0 as a result but got 0")
	}
}

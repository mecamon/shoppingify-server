//go:build integration
// +build integration

package repositories

import (
	"fmt"
	fixtures_items_repo "github.com/mecamon/shoppingify-server/__test__/fixtures/repos/items-repo"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
)

func TestItemsRepoPostgres_Register(t *testing.T) {
	insertedUserId, _ := authRepo.Register(fixtures_items_repo.User)
	cat := fixtures_items_repo.Cat
	cat.UserID = insertedUserId

	insertedCatID, _ := categoriesRepo.RegisterCategory(cat)
	item1 := fixtures_items_repo.Item1
	item1.CategoryID = insertedCatID

	_, err := itemsRepo.Register(item1)
	if err != nil {
		t.Error(err.Error())
	}

	item2 := fixtures_items_repo.Item2
	item2.CategoryID = insertedCatID

	_, err = itemsRepo.Register(item2)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestItemsRepoPostgres_GetAll(t *testing.T) {
	items, err := itemsRepo.GetAll(4, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if len(items) > 4 {
		t.Error("pagination is not working properly")
	}
}

func TestItemsRepoPostgres_GetByID(t *testing.T) {
	_, err := itemsRepo.GetByID(1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestItemsRepoPostgres_GetAllByCategoryID(t *testing.T) {
	user := fixtures_items_repo.UserForGetAllByID
	hashedPass, _ := utils.GenerateHash(user.Password)
	user.Password = hashedPass

	userId, err := authRepo.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	cat := fixtures_items_repo.CatForGetByID
	cat.UserID = userId

	insertedCatID, err := categoriesRepo.RegisterCategory(cat)
	if err != nil {
		t.Error(err.Error())
	}

	item := fixtures_items_repo.GenericItem
	item.CategoryID = insertedCatID
	for i := 0; i < 10; i++ {
		item.Name = fmt.Sprintf("%s-%d", item.Name, i)
		_, err := itemsRepo.Register(item)
		if err != nil {
			t.Error(err.Error())
		}
	}

	_, err = itemsRepo.GetAllByCategoryID(insertedCatID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestItemsRepoPostgres_Count(t *testing.T) {
	count, err := itemsRepo.Count()
	if err != nil {
		t.Error(err.Error())
	}
	if count == 0 {
		t.Error("expected more than 0 but got 0")
	}
}

func TestItemsRepoPostgres_GetDetails(t *testing.T) {
	_, err := itemsRepo.GetDetails(1)
	if err != nil {
		t.Error(err.Error())
	}
}

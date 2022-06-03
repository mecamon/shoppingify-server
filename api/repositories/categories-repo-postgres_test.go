package repositories

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/__test__/fixtures"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
	"time"
)

func TestItemsRepoPostgres_RegisterCategory(t *testing.T) {
	newUUID := uuid.NewString()
	hashedPass, _ := utils.GenerateHash(fixtures.UserForRegisterCat.Password)

	user := models.User{
		Name:      fixtures.UserForRegisterCat.Name,
		Lastname:  fixtures.UserForRegisterCat.Lastname,
		Email:     fixtures.UserForRegisterCat.Email,
		Password:  hashedPass,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: newUUID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	userID, _ := authRepo.Register(user)
	cat := models.Category{
		Name:      "Fruits",
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	_, err := categoriesRepo.RegisterCategory(cat)
	if err != nil {
		t.Error("error registering category: ", err.Error())
	}
}

func TestItemsRepoPostgres_SearchCategoryByName(t *testing.T) {
	newUUID := uuid.NewString()
	hashedPass, _ := utils.GenerateHash(fixtures.UserForSearchCat.Password)

	user := models.User{
		Name:      fixtures.UserForSearchCat.Name,
		Lastname:  fixtures.UserForSearchCat.Lastname,
		Email:     fixtures.UserForSearchCat.Email,
		Password:  hashedPass,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: newUUID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	userID, _ := authRepo.Register(user)

	cat := models.Category{
		Name:      "Meats",
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	categoriesRepo.RegisterCategory(cat)

	var skip = 0
	var take = 6
	q := "eat"
	_, err := categoriesRepo.SearchCategoryByName(q, take, skip)
	if err != nil {
		t.Error("error searching for categories: ", err.Error())
	}
}

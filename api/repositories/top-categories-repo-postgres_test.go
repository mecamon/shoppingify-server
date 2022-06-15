//go:build integration
// +build integration

package repositories

import (
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
	"time"
)

func TestInitTopCategoriesRepo(t *testing.T) {
	var i interface{}
	i = initTopCategoriesRepo(conn, config.Get())
	if _, ok := i.(TopCategoriesRepo); !ok {
		t.Error("wrong type returned")
	}
}

func TestTopCategoriesRepoImpl_Add_Error(t *testing.T) {
	_, err := topCategoriesRepo.Add(userIdForTestRepos, 67876)
	if err == nil {
		t.Error("expected error but did not get it")
	}
}

func TestTopCategoriesRepoImpl_Add_Success(t *testing.T) {
	catID, err := categoriesRepo.RegisterCategory(models.Category{
		Name:      "Cat for test topCategories 1",
		UserID:    userIdForTestRepos,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}
	_, err = topCategoriesRepo.Add(userIdForTestRepos, catID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTopCategoriesRepoImpl_Update_Error(t *testing.T) {
	err := topCategoriesRepo.Update(userIdForTestRepos, 87656)
	if err == nil {
		t.Error("error was expected but did not get it")
	}
}

func TestTopCategoriesRepoImpl_Update_Success(t *testing.T) {
	catID, err := categoriesRepo.RegisterCategory(models.Category{
		Name:      "Cat for test topCategories 2",
		UserID:    userIdForTestRepos,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}
	_, err = topCategoriesRepo.Add(userIdForTestRepos, catID)
	if err != nil {
		t.Error(err.Error())
	}

	err = topCategoriesRepo.Update(userIdForTestRepos, catID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTopCategoriesRepoImpl_GetTop(t *testing.T) {
	topCategories, err := topCategoriesRepo.GetTop(userIdForTestRepos, 3)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topCategories) == 0 {
		t.Error("expected length of topCategories must be longer than 0")
	}
}

func TestTopCategoriesRepoImpl_GetAll(t *testing.T) {
	topCategories, err := topCategoriesRepo.GetAll(userIdForTestRepos)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topCategories) == 0 {
		t.Error("expected length of topCategories must be longer than 0")
	}
}

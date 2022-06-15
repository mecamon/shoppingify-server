//go:build integration
// +build integration

package repositories

import (
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
	"time"
)

func TestInitTopItemsRepo(t *testing.T) {
	var i interface{}
	i = initTopItemsRepo(conn, config.Get())
	if _, ok := i.(TopItemsRepo); !ok {
		t.Error("wrong type returned")
	}
}

func TestTopItemsRepoImpl_Add_Error(t *testing.T) {
	_, err := topItemsRepo.Add(userIdForTestRepos, 887)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
}

func TestTopItemsRepoImpl_Add_Success(t *testing.T) {
	insertedID, err := itemsRepo.Register(models.Item{
		Name:       "Item for tests top item repos",
		Note:       "this is a note",
		CategoryID: categoryIDForTestRepos,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}
	_, err = topItemsRepo.Add(userIdForTestRepos, insertedID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTopItemsRepoImpl_Update_Error(t *testing.T) {
	err := topItemsRepo.Update(userIdForTestRepos, 87645)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
}

func TestTopItemsRepoImpl_Update_Success(t *testing.T) {
	insertedID, err := itemsRepo.Register(models.Item{
		Name:       "Item top item repos 2",
		Note:       "this is a note",
		CategoryID: categoryIDForTestRepos,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}
	_, err = topItemsRepo.Add(userIdForTestRepos, insertedID)
	if err != nil {
		t.Error(err.Error())
	}
	err = topItemsRepo.Update(userIdForTestRepos, insertedID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestTopItemsRepoImpl_GetTop(t *testing.T) {
	topItems, err := topItemsRepo.GetTop(userIdForTestRepos, 3)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topItems) == 0 {
		t.Error("top items length is expected to be longer than 0")
	}
}

func TestTopItemsRepoImpl_GetAll(t *testing.T) {
	topItems, err := topItemsRepo.GetAll(userIdForTestRepos)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topItems) == 0 {
		t.Error("top items length is expected to be longer than 0")
	}
}

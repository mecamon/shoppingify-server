package repositories

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
	"time"
)

//TODO add the build integration tag

var (
	userIDForListRepo      int64
	insertedListID         int64
	insertedItemSelectedID int64
	userForListRepo        models.User
)

func TestListsRepoPostgres_Create(t *testing.T) {
	var err error
	userForListRepo = models.User{
		Name:      "User test List repo",
		Lastname:  "Lastname",
		Email:     "user@testlistrepo.com",
		Password:  "ValidPass1234",
		IsActive:  true,
		IsVisitor: false,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	hashedPass, _ := utils.GenerateHash(userForListRepo.Password)
	userForListRepo.Password = hashedPass
	userIDForListRepo, err = authRepo.Register(userForListRepo)
	if err != nil {
		t.Error(err.Error())
	}

	list := models.List{
		Name:        "Test list 1",
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userIDForListRepo,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	}

	insertedListID, err = listsRepo.Create(list)
	if err != nil {
		t.Error(err.Error())
	}
	if insertedListID == 0 {
		t.Error("expected a list id but got 0 instead")
	}
}

func TestListsRepoPostgres_GetActive(t *testing.T) {

	listDTO, err := listsRepo.GetActive()
	if err != nil {
		t.Error(err.Error())
	}
	if listDTO.ID != insertedListID {
		t.Errorf("expected list with id: %d but instead got %d", insertedListID, listDTO.ID)
	}
}

func TestListsRepoPostgres_AddItemToList(t *testing.T) {
	var err error
	category := models.Category{
		Name:      "Category for listsRepo",
		UserID:    userIDForListRepo,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	insertedCategoryID, err := categoriesRepo.RegisterCategory(category)
	if err != nil {
		t.Error(err.Error())
	}

	item := models.Item{
		Name:       "Item for listsRepo",
		Note:       "This is a note",
		CategoryID: insertedCategoryID,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	insertedItemID, err := itemsRepo.Register(item)
	if err != nil {
		t.Error(err.Error())
	}

	itemSelected := models.SelectedItem{
		ItemID:      insertedItemID,
		Quantity:    3,
		IsCompleted: false,
		ListID:      insertedListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedItemSelectedID, err = listsRepo.AddItemToList(itemSelected)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_IsListActive(t *testing.T) {
	isActive, err := listsRepo.IsListActive(insertedListID)
	if err != nil {
		t.Error(err.Error())
	}
	if !isActive {
		t.Errorf("the list with the id: %d is inactive and expected was active", insertedListID)
	}
}

func TestListsRepoPostgres_CompleteItemSelected(t *testing.T) {
	err := listsRepo.CompleteItemSelected(insertedItemSelectedID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_DeleteItemFromList(t *testing.T) {
	err := listsRepo.DeleteItemFromList(insertedListID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_CompleteActive(t *testing.T) {
	err := listsRepo.CompleteActive()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_CancelActive(t *testing.T) {
	list := models.List{
		Name:        "Test list 2",
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userIDForListRepo,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	}
	_, err := listsRepo.Create(list)
	if err != nil {
		t.Error(err.Error())
	}
	err = listsRepo.CancelActive()
	if err != nil {
		t.Error(err.Error())
	}
}

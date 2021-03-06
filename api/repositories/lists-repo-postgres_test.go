//go:build integration
// +build integration

package repositories

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
	"time"
)

var (
	userIDForListRepo      int64
	insertedListID         int64
	insertedItemSelectedID int64
	userForListRepo        models.User
	insertedItemID         int64
)

func TestListsRepoPostgres_Create_Success(t *testing.T) {
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

func TestListsRepoPostgres_UpdateActiveListName_Success(t *testing.T) {
	err := listsRepo.UpdateActiveListName(userIDForListRepo, "new list name")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_UpdateActiveListName_Error(t *testing.T) {
	//Cancelling the active list to get an error in the upcoming step
	err := listsRepo.CancelActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
	err = listsRepo.UpdateActiveListName(userIDForListRepo, "new list name")
	if err == nil {
		t.Error("expected and error but did not get it")
	}
}

func TestListsRepoPostgres_GetActive_Success(t *testing.T) {
	var err error
	//Creating the cancelled list
	list := models.List{
		Name:        "Test list 2",
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
	listDTO, err := listsRepo.GetActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
	if listDTO.ID != insertedListID {
		t.Errorf("expected list with id: %d but instead got %d", insertedListID, listDTO.ID)
	}
}

func TestListsRepoPostgres_GetActive_Error1(t *testing.T) {
	//Cancelling the previously created list
	err := listsRepo.CancelActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = listsRepo.GetActive(userIDForListRepo)
	if err == nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_AddItemToList_Error1(t *testing.T) {
	var err error
	itemSelected := models.SelectedItem{
		ItemID:      0,
		Quantity:    3,
		IsCompleted: false,
		ListID:      insertedListID, //Cancelled at TestListsRepoPostgres_GetActive_Error1
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	insertedItemSelectedID, err = listsRepo.AddItemToList(itemSelected)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
	if err.Error() != "cannot add item to inactive list" {
		t.Error("got the wrong error message")
	}
}

func TestListsRepoPostgres_AddItemToList_Success(t *testing.T) {
	var err error

	//Creating the cancelled list
	list := models.List{
		Name:        "Test list 3",
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
		IsActive:   true,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	insertedItemID, err = itemsRepo.Register(item)
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

func TestListsRepoPostgres_AddItemToList_Error2(t *testing.T) {
	var err error
	itemSelected := models.SelectedItem{
		ItemID:      insertedItemID,
		Quantity:    3,
		IsCompleted: false,
		ListID:      insertedListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	_, err = listsRepo.AddItemToList(itemSelected)
	if err == nil {
		t.Error("expected error but did not get it")
	}
}

func TestListsRepoPostgres_UpdateItemsSelected_Error(t *testing.T) {
	items := []models.UpdateSelItemDTO{
		{ItemID: 877, Quantity: 4},
		{ItemID: 837, Quantity: 2},
	}
	err := listsRepo.UpdateItemsSelected(items)
	if err == nil {
		t.Error("expected error but did not get it")
	}
}

func TestListsRepoPostgres_UpdateItemsSelected_Success(t *testing.T) {
	items := []models.UpdateSelItemDTO{
		{ItemID: insertedItemSelectedID, Quantity: 8},
	}
	err := listsRepo.UpdateItemsSelected(items)
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

func TestListsRepoPostgres_CompleteItemSelected_Error(t *testing.T) {
	var wrongItemSelectedId int64 = 500
	err := listsRepo.CompleteItemSelected(wrongItemSelectedId)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
	if err.Error() != "item does not exist" {
		t.Error("wrong error message")
	}
}

func TestListsRepoPostgres_CompleteItemSelected_Success(t *testing.T) {
	err := listsRepo.CompleteItemSelected(insertedItemSelectedID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_DeleteItemFromList_Error(t *testing.T) {
	var wrongItemSelectedId int64 = 500
	err := listsRepo.DeleteItemFromList(wrongItemSelectedId)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
	if err.Error() != "item does not exist" {
		t.Error("wrong error message")
	}
}

func TestListsRepoPostgres_DeleteItemFromList_Success(t *testing.T) {
	err := listsRepo.DeleteItemFromList(insertedItemSelectedID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_CompleteActive_Error(t *testing.T) {
	err := listsRepo.CancelActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
	err = listsRepo.CompleteActive(userIDForListRepo)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
}

func TestListsRepoPostgres_CompleteActive_Success(t *testing.T) {
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

	err = listsRepo.CompleteActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_CancelActive_Success(t *testing.T) {
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
	err = listsRepo.CancelActive(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListsRepoPostgres_CancelActive_Error(t *testing.T) {
	err := listsRepo.CancelActive(userIDForListRepo)
	if err == nil {
		t.Error("expected an error but did not get it")
	}
}

func TestListsRepoPostgres_GetOldOnes(t *testing.T) {
	oldList, err := listsRepo.GetOldOnes(userIDForListRepo)
	if err != nil {
		t.Error(err.Error())
	}
	if len(oldList) == 0 {
		t.Error("expected length of oldList was longer than 0")
	}
}

func TestListsRepoPostgres_GetByID_Success(t *testing.T) {
	list, err := listsRepo.GetByID(userIDForListRepo, insertedListID)
	if err != nil {
		t.Error(err.Error())
	}
	if list.ID == 0 {
		t.Error("did not find any list with that ID")
	}
}

func TestListsRepoPostgres_GetByID_Error(t *testing.T) {
	_, err := listsRepo.GetByID(userIDForListRepo, 5674)
	if err == nil {
		t.Error("error expected but did not get it")
	}
}

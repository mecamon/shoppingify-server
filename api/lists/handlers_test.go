//go:build integration
// +build integration

package lists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var currentActiveListID int64

func TestInitHandler(t *testing.T) {
	var i interface{}
	appConfig := config.Get()
	i = InitHandler(appConfig)
	if _, ok := i.(*Handler); !ok {
		t.Error("expected type *Handler but got another")
	}
}

var testsCreateList = []struct {
	testName           string
	listName           string
	expectedStatusCode int
}{
	{testName: "invalid-list-name", listName: "", expectedStatusCode: http.StatusBadRequest},
	{testName: "valid", listName: "List to test handler 2", expectedStatusCode: http.StatusCreated},
	{testName: "already-one-active", listName: "List to test handler 3", expectedStatusCode: http.StatusConflict},
}

func TestHandler_Create(t *testing.T) {
	for _, tt := range testsCreateList {
		t.Log(tt.testName)
		body := struct {
			Name string `json:"name"`
		}{Name: tt.listName}
		marshalledBody, err := json.Marshal(body)
		if err != nil {
			t.Error(err.Error())
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/lists/create", bytes.NewReader(marshalledBody))
		req.Header.Set("Authorization", tokenForTests)
		req.Header.Set("Accept-Language", "en-EN")
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected status code was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

var testsGetActive = []struct {
	testName           string
	expectedStatusCode int
}{
	{testName: "active", expectedStatusCode: http.StatusOK},
	{testName: "not-active", expectedStatusCode: http.StatusNotFound},
}

func TestHandler_GetActive(t *testing.T) {
	for _, tt := range testsGetActive {
		if tt.testName == "not-active" {
			err := repositories.Main.ListsRepoImpl.CancelActive(userIdForTest)
			if err != nil {
				t.Error(err.Error())
			}
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/lists/active", nil)
		req.Header.Set("Authorization", tokenForTests)
		req.Header.Set("Accept-Language", "en-EN")
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

var testsUpdateActiveListName = []struct {
	testName           string
	listName           string
	expectedStatusCode int
}{
	{testName: "invalid-name", listName: "", expectedStatusCode: http.StatusBadRequest},
	{testName: "valid-name", listName: "list name updated", expectedStatusCode: http.StatusOK},
}

func TestHandler_UpdateActiveListName(t *testing.T) {
	var err error
	//Creating a list since it was cancelled in the previous test
	list := models.List{
		Name:        "List number for update handler",
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userIdForTest,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	}
	currentActiveListID, err = repositories.Main.ListsRepoImpl.Create(list)
	if err != nil {
		t.Error(err.Error())
	}

	for _, tt := range testsUpdateActiveListName {
		body := struct {
			Name string `json:"name"`
		}{Name: tt.listName}
		marshalledBody, _ := json.Marshal(body)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/lists/name", bytes.NewReader(marshalledBody))
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_AddItemToList(t *testing.T) {
	var testsAddItemToList = []struct {
		testName           string
		itemToAdd          models.SelectedItem
		expectedStatusCode int
	}{
		{testName: "uncompleted-fields", itemToAdd: models.SelectedItem{
			ItemID:   0,
			Quantity: 0,
			ListID:   0,
		}, expectedStatusCode: http.StatusBadRequest},
		{testName: "valid", itemToAdd: models.SelectedItem{
			ItemID:   insertedItemID2,
			Quantity: 3,
			ListID:   currentActiveListID,
		}, expectedStatusCode: http.StatusOK},
		{testName: "duplicate item", itemToAdd: models.SelectedItem{
			ItemID:   insertedItemID2,
			Quantity: 3,
			ListID:   currentActiveListID,
		}, expectedStatusCode: http.StatusConflict},
		{testName: "on-cancelled-list", itemToAdd: models.SelectedItem{
			ItemID:   insertedItemID3,
			Quantity: 2,
			ListID:   currentActiveListID,
		}, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range testsAddItemToList {
		t.Log(tt.testName)
		if tt.testName == "on-cancelled-list" {
			err := repositories.Main.ListsRepoImpl.CancelActive(userIdForTest)
			if err != nil {
				t.Error(err.Error())
			}
		}
		marshalled, err := json.Marshal(tt.itemToAdd)
		if err != nil {
			t.Error(err.Error())
		}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/lists/add-item", bytes.NewReader(marshalled))
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_UpdateItemsSelected(t *testing.T) {
	currentActiveListID, _ = repositories.Main.ListsRepoImpl.Create(models.List{
		Name:        "list for test lists handler 3",
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userIdForTest,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	})
	insertedItemToListID1, _ := repositories.Main.ListsRepoImpl.AddItemToList(models.SelectedItem{
		ItemID:      insertedItemID1,
		Quantity:    1,
		IsCompleted: false,
		ListID:      currentActiveListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})
	insertedItemToListID2, _ := repositories.Main.ListsRepoImpl.AddItemToList(models.SelectedItem{
		ItemID:      insertedItemID2,
		Quantity:    1,
		IsCompleted: false,
		ListID:      currentActiveListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})

	var testsUpdateItemsSelected = []struct {
		testName           string
		items              []models.UpdateSelItemDTO
		expectedStatusCode int
	}{
		{testName: "invalid-itemID", items: []models.UpdateSelItemDTO{
			{ItemID: 0, Quantity: 3},
			{ItemID: insertedItemToListID1, Quantity: 3},
		}, expectedStatusCode: http.StatusBadRequest},
		{testName: "invalid-quantity", items: []models.UpdateSelItemDTO{
			{ItemID: insertedItemToListID1, Quantity: 0},
			{ItemID: insertedItemToListID2, Quantity: 2},
		}, expectedStatusCode: http.StatusBadRequest},
		{testName: "valid", items: []models.UpdateSelItemDTO{
			{ItemID: insertedItemToListID1, Quantity: 15},
			{ItemID: insertedItemToListID2, Quantity: 4},
		}, expectedStatusCode: http.StatusOK},
	}

	for _, tt := range testsUpdateItemsSelected {
		t.Log(tt.testName)
		marshalled, err := json.Marshal(tt.items)
		if err != nil {
			t.Error(err.Error())
		}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/lists/update-items", bytes.NewReader(marshalled))
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_DeleteItemFromList(t *testing.T) {
	category := models.Category{
		Name:      "Cat for del from list 1",
		UserID:    userIdForTest,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	insertedCatID, err := repositories.Main.CategoriesRepoImpl.RegisterCategory(category)
	if err != nil {
		t.Error(err.Error())
	}

	item := models.Item{
		Name:       "item for del list 1",
		Note:       "item for del list 1",
		CategoryID: insertedCatID,
		IsActive:   true,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	insertedItemID, err := repositories.Main.ItemsRepoIpml.Register(item)
	if err != nil {
		t.Error(err.Error())
	}

	insertedItemToListID, err := repositories.Main.ListsRepoImpl.AddItemToList(models.SelectedItem{
		ItemID:      insertedItemID,
		Quantity:    1,
		IsCompleted: false,
		ListID:      currentActiveListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}

	var testsDeleteItemFromList = []struct {
		testName           string
		selItemID          int64
		expectedStatusCode int
	}{
		{testName: "not-saved-item", selItemID: 899, expectedStatusCode: http.StatusNotFound},
		{testName: "valid-request", selItemID: insertedItemToListID, expectedStatusCode: http.StatusOK},
	}

	for _, tt := range testsDeleteItemFromList {
		t.Log(tt.testName)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/lists/selected-items/%d", tt.selItemID), nil)
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)

		Router.ServeHTTP(rr, req)
		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_CompleteItemSelected(t *testing.T) {
	category := models.Category{
		Name:      "Cat for del from list 2",
		UserID:    userIdForTest,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	insertedCatID, err := repositories.Main.CategoriesRepoImpl.RegisterCategory(category)
	if err != nil {
		t.Error(err.Error())
	}

	item := models.Item{
		Name:       "item for del list 2",
		Note:       "item for del list 2",
		CategoryID: insertedCatID,
		IsActive:   true,
		ImageURL:   "",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	insertedItemID, err := repositories.Main.ItemsRepoIpml.Register(item)
	if err != nil {
		t.Error(err.Error())
	}

	insertedItemToListID, _ := repositories.Main.ListsRepoImpl.AddItemToList(models.SelectedItem{
		ItemID:      insertedItemID,
		Quantity:    1,
		IsCompleted: false,
		ListID:      currentActiveListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})

	var testsCompleteItemSelected = []struct {
		testName           string
		selItemID          int64
		expectedStatusCode int
	}{
		{testName: "required-field", selItemID: 0, expectedStatusCode: http.StatusBadRequest},
		{testName: "not-saved-item", selItemID: 899, expectedStatusCode: http.StatusNotFound},
		{testName: "valid-request", selItemID: insertedItemToListID, expectedStatusCode: http.StatusOK},
	}

	for _, tt := range testsCompleteItemSelected {
		t.Log(tt.testName)
		body := struct {
			ItemSelID int64 `json:"item_sel_id"`
		}{ItemSelID: tt.selItemID}
		marshalled, err := json.Marshal(body)
		if err != nil {
			t.Error(err.Error())
		}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/lists/selected-items", bytes.NewReader(marshalled))
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)

		Router.ServeHTTP(rr, req)
		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_CancelActive(t *testing.T) {
	var testsCancelActive = []struct {
		testName           string
		expectedStatusCode int
	}{
		{testName: "valid-request", expectedStatusCode: http.StatusOK},
		{testName: "not-active-list", expectedStatusCode: http.StatusNotFound},
	}

	for _, tt := range testsCancelActive {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/lists/cancel-active", nil)
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_CompleteActive(t *testing.T) {
	//Creating a new list because the last one was cancelled in the previous test
	currentActiveListID, _ = repositories.Main.ListsRepoImpl.Create(models.List{
		Name:        "List for test lists handlers 4",
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userIdForTest,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	})

	var testsCompleteActive = []struct {
		testName           string
		expectedStatusCode int
	}{
		{testName: "valid-request", expectedStatusCode: http.StatusOK},
		{testName: "not-active-list", expectedStatusCode: http.StatusNotFound},
	}

	for _, tt := range testsCompleteActive {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/lists/cancel-active", nil)
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_GetOldLists(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/lists/old-lists", nil)
	req.Header.Set("Accept-Language", "en-EN")
	req.Header.Set("Authorization", tokenForTests)
	Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected statusCode was %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestHandler_GetByID(t *testing.T) {
	var testsGetByID = []struct {
		testName           string
		routeParam         int64
		expectedStatusCode int
	}{
		{testName: "valid-request", routeParam: currentActiveListID, expectedStatusCode: http.StatusOK},
		{testName: "invalid-route-param", routeParam: 3442343, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range testsGetByID {
		t.Log(tt.testName)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/lists/%d", tt.routeParam), nil)
		req.Header.Set("Accept-Language", "en-EN")
		req.Header.Set("Authorization", tokenForTests)
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected statusCode was %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

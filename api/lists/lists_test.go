//go:build !integration
// +build !integration

package lists

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

var validItemToAddTests = []struct {
	testName       string
	itemToAdd      models.SelectedItem
	isValid        bool
	expectedErrors int
}{
	{testName: "not-itemID-id", isValid: false, expectedErrors: 1, itemToAdd: models.SelectedItem{
		ItemID:   0,
		Quantity: 3,
		ListID:   1,
	}},
	{testName: "not-listID-id", isValid: false, expectedErrors: 1, itemToAdd: models.SelectedItem{
		ItemID:   3,
		Quantity: 3,
		ListID:   0,
	}},
	{testName: "not-quantity-id", isValid: false, expectedErrors: 1, itemToAdd: models.SelectedItem{
		ItemID:   10,
		Quantity: 0,
		ListID:   1,
	}},
	{testName: "valid-item", isValid: true, expectedErrors: 0, itemToAdd: models.SelectedItem{
		ItemID:   4,
		Quantity: 3,
		ListID:   7,
	}},
}

func TestDomLists_validateItemToAdd(t *testing.T) {
	for _, tt := range validItemToAddTests {
		t.Log(tt.testName)
		domLists := DomLists{
			appLocales: appi18n.GetLocales("en-EN"),
			itemToAdd:  tt.itemToAdd,
		}
		isValid, errMap := domLists.validateItemToAdd()
		if isValid != tt.isValid {
			t.Errorf("expected valid: %v but got valid: %v", tt.isValid, isValid)
		}
		if len(errMap) != tt.expectedErrors {
			t.Errorf("expected errors were: %d but got %d", tt.expectedErrors, len(errMap))
		}
	}
}

func TestDomLists_completeItemToAdd(t *testing.T) {
	item := models.SelectedItem{
		ItemID:   2,
		Quantity: 5,
		ListID:   7,
	}
	dom := DomLists{
		appLocales: appi18n.GetLocales("en-EN"),
		itemToAdd:  item,
	}
	completedItem := dom.completeItemToAdd()
	if completedItem.UpdatedAt == 0 || completedItem.CreatedAt == 0 {
		t.Error("item is not completed")
	}
}

var testsItemToUpdate = []struct {
	testName       string
	itemToUpdate   []models.UpdateSelItemDTO
	expectedErrors int
	isValid        bool
}{
	{testName: "not-itemID", expectedErrors: 2, isValid: false, itemToUpdate: []models.UpdateSelItemDTO{
		{ItemID: 0, Quantity: 2},
		{ItemID: 0, Quantity: 1},
	}},
	{testName: "not-quantity", expectedErrors: 2, isValid: false, itemToUpdate: []models.UpdateSelItemDTO{
		{ItemID: 2, Quantity: 0},
		{ItemID: 4, Quantity: 0},
	}},
	{testName: "valid-item", expectedErrors: 0, isValid: true, itemToUpdate: []models.UpdateSelItemDTO{
		{ItemID: 2, Quantity: 7},
		{ItemID: 4, Quantity: 2},
	}},
}

func TestDomLists_validateItemToUpdate(t *testing.T) {
	for _, tt := range testsItemToUpdate {
		dom := DomLists{
			appLocales:    appi18n.GetLocales("en-EN"),
			itemsToUpdate: tt.itemToUpdate,
		}
		isValid, errMap := dom.validateItemToUpdate()
		if isValid != tt.isValid {
			t.Errorf("expected validity of %v but got %v", tt.isValid, isValid)
		}
		if len(errMap) != tt.expectedErrors {
			t.Errorf("expected errors %d but got %d", tt.expectedErrors, len(errMap))
		}
	}
}

//go:build !integration
// +build !integration

package top_categories

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

func TestAddPercentages(t *testing.T) {
	err := appi18n.InitLocales()
	if err != nil {
		t.Error(err.Error())
	}

	topCats := []models.TopCategoryDTO{
		{ID: 1, CategoryID: 3, Name: "Cat 1", SumQuantity: 3},
		{ID: 2, CategoryID: 1, Name: "Cat 2", SumQuantity: 3},
	}
	allCats := []models.TopCategoryDTO{
		{ID: 1, CategoryID: 3, Name: "Cat 1", SumQuantity: 3},
		{ID: 2, CategoryID: 7, Name: "Cat 2", SumQuantity: 3},
		{ID: 3, CategoryID: 3, Name: "Cat 3", SumQuantity: 2},
		{ID: 4, CategoryID: 9, Name: "Cat 4", SumQuantity: 2},
	}
	topCatsWithPercentages := addPercentages(topCats, allCats)

	if topCatsWithPercentages[0].Name == "" {
		t.Error("Does not contain name")
	}
	if topCatsWithPercentages[0].Percentage != 30 {
		t.Error("wrong percentage")
	}
	if topCatsWithPercentages[1].Percentage != 30 {
		t.Error("wrong percentage")
	}
}

//go:build !integration
// +build !integration

package top_items

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

	topItems := []models.TopItemDTO{
		{ID: 1, ItemID: 3, Name: "Cat 1", SumQuantity: 3},
		{ID: 2, ItemID: 1, Name: "Cat 2", SumQuantity: 3},
	}
	allItems := []models.TopItemDTO{
		{ID: 1, ItemID: 3, Name: "Cat 1", SumQuantity: 3},
		{ID: 2, ItemID: 7, Name: "Cat 2", SumQuantity: 3},
		{ID: 3, ItemID: 3, Name: "Cat 3", SumQuantity: 2},
		{ID: 4, ItemID: 9, Name: "Cat 4", SumQuantity: 2},
	}
	topCatsWithPercentages := addPercentages(topItems, allItems)

	if topCatsWithPercentages[0].Percentage != 30 {
		t.Error("wrong percentage")
	}
	if topCatsWithPercentages[1].Percentage != 30 {
		t.Error("wrong percentage")
	}
}

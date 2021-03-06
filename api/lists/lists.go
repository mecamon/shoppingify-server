package lists

import (
	"fmt"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type DomLists struct {
	appLocales    appi18n.AppLocales
	itemToAdd     models.AddSelectedItemDTO
	itemsToUpdate []models.UpdateSelItemDTO
}

func (d *DomLists) validateItemToAdd() (bool, models.ErrorMap) {
	var errMap = models.ErrorMap{}
	if d.itemToAdd.ItemID == 0 {
		td := map[string]interface{}{"Field": "item_id"}
		msg := d.appLocales.GetMsg("RequiredField", td)
		errMap["itemID"] = msg
	}
	if d.itemToAdd.ListID == 0 {
		td := map[string]interface{}{"Field": "list_id"}
		msg := d.appLocales.GetMsg("RequiredField", td)
		errMap["listID"] = msg
	}
	if d.itemToAdd.Quantity == 0 {
		td := map[string]interface{}{"Field": "quantity"}
		msg := d.appLocales.GetMsg("RequiredField", td)
		errMap["quantity"] = msg
	}
	return len(errMap) == 0, errMap
}

func (d *DomLists) completeItemToAdd() models.SelectedItem {
	itemToAddCompleted := models.SelectedItem{
		ItemID:      d.itemToAdd.ItemID,
		Quantity:    d.itemToAdd.Quantity,
		IsCompleted: false,
		ListID:      d.itemToAdd.ListID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	return itemToAddCompleted
}

func (d *DomLists) validateItemToUpdate() (bool, models.ErrorMap) {
	var errMap = models.ErrorMap{}
	for i, item := range d.itemsToUpdate {
		if item.ItemID == 0 {
			td := map[string]interface{}{"Field": "item_id"}
			msg := d.appLocales.GetMsg("RequiredField", td)
			errMap[fmt.Sprintf("line: %d", i)] = msg
		}
		if item.Quantity == 0 {
			td := map[string]interface{}{"Field": "quantity"}
			msg := d.appLocales.GetMsg("RequiredField", td)
			errMap[fmt.Sprintf("line: %d", i)] = msg
		}
	}
	return len(errMap) == 0, errMap
}

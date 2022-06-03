package categories

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type CategoryDom struct {
	appLocales appi18n.AppLocales
	category   models.Category
}

func (c *CategoryDom) validCategory() (bool, models.ErrorMap) {
	errMap := models.ErrorMap{}
	if len(c.category.Name) == 0 {
		errMap["name"] = c.appLocales.GetMsg("CategoryNameError", nil)
	}
	if c.category.UserID == 0 {
		errMap["user_id"] = c.appLocales.GetMsg("UserIDRequired", nil)
	}
	return len(errMap) == 0, errMap
}

func (c *CategoryDom) completeCategory() models.Category {
	return models.Category{
		Name:      c.category.Name,
		UserID:    c.category.UserID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

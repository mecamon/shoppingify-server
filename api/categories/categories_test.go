//go:build !integration
// +build !integration

package categories

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

var validCategoryTests = []struct {
	testName       string
	category       models.Category
	isValid        bool
	expectedErrors int
}{
	{testName: "invalid-category", category: models.Category{
		Name:   "",
		UserID: 0,
	},
		isValid:        false,
		expectedErrors: 2,
	},
	{testName: "valid-category", category: models.Category{
		Name:   "Vegetables",
		UserID: 3,
	},
		isValid:        true,
		expectedErrors: 0,
	},
}

func TestCategories_ValidCategory(t *testing.T) {
	err := appi18n.InitLocales()
	if err != nil {
		t.Error(err.Error())
	}

	for _, tt := range validCategoryTests {
		t.Log(tt.testName)
		{
			cd := CategoryDom{
				appLocales: appi18n.GetLocales("en-EN"),
				category:   tt.category,
			}
			isCorrect, errMap := cd.validCategory()
			if isCorrect != tt.isValid {
				t.Errorf("got %v when expected was %v", isCorrect, tt.isValid)
			}
			if len(errMap) != tt.expectedErrors {
				t.Errorf("expected %d errors but got %d", tt.expectedErrors, len(errMap))
			}
		}
	}
}

func TestCategories_CompleteCategory(t *testing.T) {
	cd := CategoryDom{
		appLocales: appi18n.GetLocales("en-EN"),
		category: models.Category{
			Name:   "Fruits",
			UserID: 4,
		},
	}
	completedCat := cd.completeCategory()
	if completedCat.CreatedAt == 0 || completedCat.UpdatedAt == 0 {
		t.Error("category is not completed")
	}
}

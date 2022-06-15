//go:build integration
// +build integration

package top_categories

import (
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInitHandler(t *testing.T) {
	var i interface{}
	i = InitHandler(config.Get())
	if _, ok := i.(*Handler); !ok {
		t.Error("wrong conversion type")
	}
}

func TestHandler_GetTop(t *testing.T) {
	insertedCatID, err := repositories.Main.CategoriesRepoImpl.RegisterCategory(models.Category{
		Name:      "Cat test topcat1",
		UserID:    userIdForTest,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		t.Error(err.Error())
	}

	_, err = repositories.Main.TopCategoriesImpl.Add(userIdForTest, insertedCatID)
	if err != nil {
		t.Error(err.Error())
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/top-categories", nil)
	req.Header.Set("Accept-Language", "en-EN")
	req.Header.Set("Authorization", tokenForTests)
	Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected statusCode was %d but got %d", http.StatusOK, rr.Code)
	}
}

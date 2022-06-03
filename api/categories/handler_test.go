//go:build integration
// +build integration

package categories

import (
	"bytes"
	"encoding/json"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInitHandler(t *testing.T) {
	var i interface{}
	appConfig := config.Get()
	i = InitHandler(appConfig)

	if _, ok := i.(*Handler); !ok {
		t.Error("wrong type")
	}
}

var createTests = []struct {
	testName     string
	cat          models.Category
	expectedCode int
}{
	{testName: "invalid-category-name", cat: models.Category{Name: ""}, expectedCode: http.StatusBadRequest},
	{testName: "valid-category", cat: models.Category{Name: "Plastics"}, expectedCode: http.StatusCreated},
}

func TestHandler_Create(t *testing.T) {
	for _, tt := range createTests {
		t.Log(tt.testName)
		body, _ := json.Marshal(&tt.cat)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/categories/", bytes.NewReader(body))
		req.Header.Set("Authorization", userToken)

		Router.ServeHTTP(rr, req)
		if rr.Code != tt.expectedCode {
			t.Errorf("got %v statusCode when expected was %v", rr.Code, tt.expectedCode)
		}
	}
}

func TestHandler_GetAllByName(t *testing.T) {
	q := "ea"

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/?q="+q, nil)
	req.Header.Set("Authorization", userToken)

	Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %d when expected was %d", rr.Code, http.StatusOK)
	}

	var body []models.CategoryDTO
	err := json.NewDecoder(rr.Result().Body).Decode(&body)
	if err != nil {
		t.Error(err.Error())
	}

	for _, cat := range body {
		if !strings.Contains(cat.Name, q) {
			t.Error("categories do not contain the query searched")
			break
		}
	}
}

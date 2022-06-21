//go:build integration
// +build integration

package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mecamon/shoppingify-server/__test__/fixtures/items"
	"github.com/mecamon/shoppingify-server/api/repositories"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var createItemsTests = []struct {
	testName           string
	item               models.Item
	expectedStatusCode int
	expectedErrors     int
}{
	{
		testName:           "invalid-item",
		item:               models.Item{Name: wordsGenerator(40), Note: wordsGenerator(113)},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrors:     2,
	},
	{
		testName:           "invalid-item-name",
		item:               models.Item{Name: "", Note: wordsGenerator(113)},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrors:     1,
	},
	{
		testName:           "valid-item",
		item:               models.Item{Name: "Pineapple", Note: "This is just a fruit"},
		expectedStatusCode: http.StatusCreated,
		expectedErrors:     0,
	},
	{
		testName:           "duplicated-item-name",
		item:               models.Item{Name: "Pineapple", Note: "This is just a fruit"},
		expectedStatusCode: http.StatusConflict,
		expectedErrors:     1,
	},
}

func TestHandler_Create(t *testing.T) {
	repos := repositories.Main
	insertedUserID, err := repos.AuthRepoImpl.Register(fixtures_items.User)
	if err != nil {
		t.Error(err.Error())
	}
	token, err := json_web_token.Generate(insertedUserID, "")
	if err != nil {
		t.Error(err.Error())
	}
	cat := fixtures_items.Cat
	cat.UserID = insertedUserID
	insertedCatId, err := repos.CategoriesRepoImpl.RegisterCategory(cat)

	for _, tt := range createItemsTests {
		t.Log(tt.testName)
		{
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", tt.item.Name)
			writer.WriteField("note", tt.item.Note)
			writer.WriteField("category_id", strconv.FormatInt(insertedCatId, 10))
			writer.Close()

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/items/", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", token)

			Router.ServeHTTP(rr, req)
			if tt.expectedStatusCode != rr.Code {
				t.Errorf("expected %d but got %d", tt.expectedStatusCode, rr.Code)
			}
		}
	}
}

func TestHandler_GetByCategoryGroups(t *testing.T) {
	//Creating user
	user := fixtures_items.UserForGet
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	repos := repositories.Main
	insertedUserId, err := repos.AuthRepoImpl.Register(user)
	if err != nil {
		t.Error(err.Error())
	}

	token, err := json_web_token.Generate(insertedUserId, "")
	if err != nil {
		t.Error(err.Error())
	}

	cat1 := fixtures_items.GenericCat
	cat1.UserID = insertedUserId
	cat1.Name = fmt.Sprintf("%s-1", cat1.Name)
	insertedCatID1, err := repos.CategoriesRepoImpl.RegisterCategory(cat1)
	if err != nil {
		t.Error(err.Error())
	}

	cat2 := fixtures_items.GenericCat
	cat2.UserID = insertedUserId
	cat2.Name = fmt.Sprintf("%s-2", cat2.Name)
	insertedCatID2, err := repos.CategoriesRepoImpl.RegisterCategory(cat2)
	if err != nil {
		t.Error(err.Error())
	}

	//Seeding information
	for i := 0; i < 10; i++ {
		var cat int64
		if i < 6 {
			cat = insertedCatID1
		} else {
			cat = insertedCatID2
		}

		item := fixtures_items.GenericItem
		item.Name = fmt.Sprintf("%s-%d", item.Name, i)
		item.CategoryID = cat
		_, err := repos.ItemsRepoIpml.Register(item)
		if err != nil {
			t.Error(err.Error())
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/items/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected statusCode %d but got %d", http.StatusOK, rr.Code)
	}

	value := rr.Result().Header.Get("X-Total-Count")
	count, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Error(err.Error())
	}
	if count == 0 {
		t.Error("expected more than 0 from X-Total-Count header but got 0")
	}

	var body []models.CategoriesGroup
	err = json.NewDecoder(rr.Result().Body).Decode(&body)
	if err != nil {
		t.Error(err.Error())
	} else {
		log.Println(len(body))
	}
}

var token string

func TestHandler_GetDetailsByID(t *testing.T) {
	user := fixtures_items.UserForGetOne
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		t.Error(err.Error())
	}
	user.Password = hashedPass
	repos := repositories.Main
	insertedUserId, err := repos.AuthRepoImpl.Register(user)
	if err != nil {
		t.Error(err.Error())
	}
	token, err = json_web_token.Generate(insertedUserId, "")
	if err != nil {
		t.Error(err.Error())
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/items/1"), nil)
	req.Header.Set("Authorization", token)
	Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestHandler_GetDetailsByID2(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/items/"+"aa", nil)
	req.Header.Set("Authorization", token)
	Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected %d but got %d", http.StatusNotFound, rr.Code)
	}
}

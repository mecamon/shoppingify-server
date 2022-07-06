//go:build integration
// +build integration

package items_summary

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/api/repositories"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_GetByMonth(t *testing.T) {
	var err error

	var getByMonthTests = []struct {
		testName           string
		expectedStatusCode int
		year               int
	}{
		{testName: "invalid-route-param", expectedStatusCode: http.StatusBadRequest, year: 0},
		{testName: "no-content-found", expectedStatusCode: http.StatusNoContent, year: 2345},
		{testName: "successful request", expectedStatusCode: http.StatusOK, year: time.Now().Year()},
	}

	err = repositories.Main.ItemsSummaryRepoImpl.Add(HandlerUserID, models.ItemsSummary{
		Month:    int(time.Now().Month()),
		Year:     time.Now().Year(),
		Quantity: 5,
	})
	if err != nil {
		t.Error(err.Error())
	}

	for _, tt := range getByMonthTests {
		t.Log(tt.testName)
		rr := httptest.NewRecorder()
		var req *http.Request
		if tt.testName == "invalid-route-param" {
			req = httptest.NewRequest(http.MethodGet, "/api/summary/abc", nil)
		} else {
			req = httptest.NewRequest(http.MethodGet, "/api/summary/"+fmt.Sprintf("%d", tt.year), nil)
		}

		req.Header.Set("Authorization", Token)
		req.Header.Set("Content-Type", "application/json")
		Router.ServeHTTP(rr, req)
		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected status %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

func TestHandler_GetByYear(t *testing.T) {

	//Registering new user
	validPass := "Validpass09"
	passwordHashed, err := utils.GenerateHash(validPass)
	if err != nil {
		log.Println(err.Error())
	}
	HandlerUserID, err = repositories.Main.AuthRepoImpl.Register(models.User{
		Name:      "Get by month H2",
		Lastname:  "lastname",
		Email:     "getbymonth@handler2.com",
		Password:  passwordHashed,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})

	newUserToken, err := json_web_token.Generate(HandlerUserID, "")
	if err != nil {
		log.Println(err.Error())
	}

	var getByYearTests = []struct {
		testName           string
		expectedStatusCode int
	}{
		{testName: "no-content-found", expectedStatusCode: http.StatusNoContent},
		{testName: "successful request", expectedStatusCode: http.StatusOK},
	}

	for _, tt := range getByYearTests {
		t.Log(tt.testName)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
		req.Header.Set("Content-Type", "application/json")
		if tt.testName == "successful request" {
			req.Header.Set("Authorization", Token)
		} else {
			req.Header.Set("Authorization", newUserToken)
		}

		Router.ServeHTTP(rr, req)
		if rr.Code != tt.expectedStatusCode {
			t.Errorf("expected status %d but got %d", tt.expectedStatusCode, rr.Code)
		}
	}
}

//go:build integration
// +build integration

package repositories

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"testing"
	"time"
)

func TestInitRepos(t *testing.T) {
	var i interface{}
	i = initAuthRepo(conn, config.Get())

	if _, ok := i.(AuthRepo); !ok {
		t.Errorf("wrong return type in authRepo")
	}
}

var authRepoTestUser = models.User{
	Name:      "Test Repo",
	Lastname:  "Test Repo",
	Email:     "test@repo.com",
	Password:  "TestRepoPass09876",
	IsActive:  true,
	IsVisitor: false,
}

func TestAuthRepoPostgres_Register(t *testing.T) {
	newUUID := uuid.NewString()
	hashedPass, _ := utils.GenerateHash(authRepoTestUser.Password)

	user := models.User{
		Name:      authRepoTestUser.Name,
		Lastname:  authRepoTestUser.Lastname,
		Email:     authRepoTestUser.Email,
		Password:  hashedPass,
		IsActive:  authRepoTestUser.IsActive,
		IsVisitor: authRepoTestUser.IsVisitor,
		LoginCode: newUUID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	_, err := authRepo.Register(user)
	if err != nil {
		t.Error("error inserting new user:", err.Error())
	}
}

func TestAuthRepoPostgres_SearchUserByEmail(t *testing.T) {
	user, err := authRepo.SearchUserByEmail(authRepoTestUser.Email)
	if user.Email != authRepoTestUser.Email {
		t.Errorf("got %s when expected was %s", user.Email, authRepoTestUser.Email)
	}
	if err != nil {
		t.Error("error searching user for email: ", err.Error())
	}
}

func TestAuthRepoPostgres_CheckUserPassword(t *testing.T) {
	isCorrect, err := authRepo.CheckUserPassword(authRepoTestUser.Email, authRepoTestUser.Password)
	if !isCorrect {
		t.Error("invalid user password")
	}
	if err != nil {
		t.Errorf("error checking pass: %s", err.Error())
	}
}

package auth

import (
	"github.com/google/uuid"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"time"
)

func validCredentials(user models.User, lang string) (bool, models.ErrorMap) {
	appLocales := appi18n.GetLocales(lang)
	errMap := models.ErrorMap{}
	nameMinimalLength := 2

	if validEmail := utils.HasValidEmail(user.Email); !validEmail {
		errMap["email"] = appLocales.GetMsg("InvalidEmail", nil)
	}
	if validPassword := utils.HasValidPass(user.Password); !validPassword {
		errMap["password"] = appLocales.GetMsg("InvalidPassword", nil)
	}
	if validName := len(user.Name) > nameMinimalLength; !validName {
		td := map[string]interface{}{"Label": "name", "Count": nameMinimalLength}
		errMap["name"] = appLocales.GetMsg("TooShortField", td)
	}
	if validLastname := len(user.Name) > nameMinimalLength; !validLastname {
		td := map[string]interface{}{"Label": "lastname", "Count": nameMinimalLength}
		errMap["lastname"] = appLocales.GetMsg("TooShortField", td)
	}

	return len(errMap) == 0, errMap
}

func completeUserInformation(user models.User) (models.User, error) {
	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		return models.User{}, err
	}

	newUUID := uuid.NewString()
	completedUser := models.User{
		Name:      user.Name,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  hashedPass,
		IsActive:  true,
		IsVisitor: false,
		LoginCode: newUUID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	return completedUser, nil
}

func createVisitorInformation() models.User {
	return models.User{
		IsVisitor: true,
		LoginCode: uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

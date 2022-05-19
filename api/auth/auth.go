package auth

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
)

func validCredentials(user models.User, lang string) (bool, models.ErrorMap) {
	appLocales := appi18n.GetLocales(lang)
	errMap := models.ErrorMap{}

	if validEmail := utils.HasValidEmail(user.Email); !validEmail {
		errMap["email"] = appLocales.GetMsg("InvalidEmail")
	}
	if validPassword := utils.HasValidPass(user.Password); !validPassword {
		errMap["password"] = appLocales.GetMsg("InvalidPassword")
	}
	if validUsername := len(user.Username) > 4; !validUsername {
		errMap["username"] = appLocales.GetMsg("invalidUsername")
	}

	return len(errMap) == 0, errMap
}

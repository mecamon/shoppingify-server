package auth

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func validCredentials(user models.User, lang string) (bool, models.ErrorMap) {
	locales := appi18n.GetLocales(lang)
	errMap := models.ErrorMap{}

	if validEmail := utils.HasValidEmail(user.Email); !validEmail {
		errMap["email"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidEmail"})
	}
	if validPassword := utils.HasValidPass(user.Password); !validPassword {
		errMap["password"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "InvalidPassword"})
	}
	if validUsername := len(user.Username) > 4; !validUsername {
		errMap["username"] = locales.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalidUsername"})
	}

	return len(errMap) == 0, errMap
}

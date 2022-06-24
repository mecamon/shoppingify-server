package appi18n

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var multiLocales *[]AppLocales

type AppLocales struct {
	Lang    string
	Locales *i18n.Localizer
}

func InitLocales() error {
	var enPath string
	var bundle *i18n.Bundle
	var localesSlice []AppLocales

	if flag.Lookup("test.v") == nil {
		enPath = "./i18n/catalog/active.en.toml"
	} else {
		enPath = "../../i18n/catalog/active.en.toml"
	}

	//English built
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile(enPath)
	if err != nil {
		return err
	}

	localesSlice = append(localesSlice, AppLocales{
		Lang:    "en-EN",
		Locales: i18n.NewLocalizer(bundle, "en-EN"),
	})

	multiLocales = &localesSlice
	
	return nil
}

func GetLocales(lang string) AppLocales {
	for _, appLocales := range *multiLocales {
		if appLocales.Lang == lang {
			return appLocales
		}
	}
	return (*multiLocales)[0]
}

func (a *AppLocales) GetMsg(messageId string, tempData map[string]interface{}) string {
	msg := a.Locales.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: tempData,
	})
	return msg
}

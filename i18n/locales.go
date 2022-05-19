package appi18n

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

func GetLocales(lang string) *i18n.Localizer {
	var enPath string
	var bundle *i18n.Bundle

	if flag.Lookup("test.v") == nil {
		enPath = "./i18n/catalog/active.en.toml"
	} else {
		enPath = "../../i18n/catalog/active.en.toml"
	}

	switch lang {
	case "en-EN":
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		_, err := bundle.LoadMessageFile(enPath)
		if err != nil {
			log.Println(err)
		}
	default:
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		_, err := bundle.LoadMessageFile(enPath)
		if err != nil {
			log.Println(err)
		}
	}

	return i18n.NewLocalizer(bundle, lang)
}

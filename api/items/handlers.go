package items

import (
	"encoding/json"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"net/http"
	"strconv"
)

type Handler struct {
	app   *config.App
	repos repositories.MainRepo
}

var handler *Handler

func InitHandler(conf *config.App) *Handler {
	return &Handler{
		app:   conf,
		repos: repositories.Main,
	}
}

// WAITING FOR CATEGORIES TO GET DONE
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	appLocales := appi18n.GetLocales(lang)

	categoryID, _ := strconv.ParseInt(r.Form.Get("category_id"), 10, 64)
	item := models.Item{
		Name:       r.Form.Get("name"),
		Note:       r.Form.Get("note"),
		CategoryID: categoryID,
	}

	_, _, err := r.FormFile("file")
	if err != nil {
		//TODO
	}
	fileInfo := FileInfo{
		//Size:        fileHeader.Size,
		//ContentType: fileHeader.Header.Get("Content-Type"),
	}

	itemDom := ItemDom{
		appLocales: appLocales,
		item:       item,
		fileInfo:   fileInfo,
	}

	isValid, errMap := itemDom.validItem()
	if !isValid {
		output, _ := json.MarshalIndent(errMap, "", "    ")
		utils.Response(w, http.StatusBadRequest, output)
		return
	}
}

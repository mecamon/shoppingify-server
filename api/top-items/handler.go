package top_items

import (
	"encoding/json"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/utils"
	"net/http"
)

var handler *Handler

type Handler struct {
	app   *config.App
	repos repositories.MainRepo
}

func InitHandler(conf *config.App) *Handler {
	handler = &Handler{app: conf, repos: repositories.Main}
	return handler
}

func (h *Handler) GetTop(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Accept-Language")
	userID := r.Context().Value("ID").(int64)

	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 3)

	topItems, err := h.repos.TopItemsImpl.GetTop(userID, take)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	allItems, err := h.repos.TopItemsImpl.GetAll(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	completedItems := addPercentages(topItems, allItems)
	output, err := json.Marshal(completedItems)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

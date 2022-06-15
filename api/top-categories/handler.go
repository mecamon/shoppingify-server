package top_categories

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

	topCategories, err := h.repos.TopCategoriesImpl.GetTop(userID, take)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	allCategories, err := h.repos.TopCategoriesImpl.GetAll(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	completedCategories := addPercentages(topCategories, allCategories)
	output, err := json.Marshal(completedCategories)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

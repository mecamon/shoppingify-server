package categories

import (
	"encoding/json"
	"fmt"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"net/http"
)

type Handler struct {
	App   *config.App
	Repos repositories.MainRepo
}

var handler *Handler

func InitHandler(app *config.App) *Handler {
	handler = &Handler{App: app, Repos: repositories.Main}
	return handler
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	appLocales := appi18n.GetLocales(lang)
	userId := r.Context().Value("ID").(int64)
	var category models.Category

	json.NewDecoder(r.Body).Decode(&category)
	category.UserID = userId

	categoryDom := CategoryDom{
		appLocales: appLocales,
		category:   category,
	}

	isValid, errMap := categoryDom.validCategory()
	if !isValid {
		output, _ := json.MarshalIndent(errMap, "", "    ")
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	completedCategory := categoryDom.completeCategory()
	ID, err := h.Repos.CategoriesRepoImpl.RegisterCategory(completedCategory)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
		panic(w)
	}
	_, err = h.Repos.TopCategoriesImpl.Add(userId, ID)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	res := map[string]interface{}{"InsertedID": ID}
	output, _ := json.MarshalIndent(res, "", "    ")
	utils.Response(w, http.StatusCreated, output)
}

func (h *Handler) GetAllByName(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 4)

	skipQuery := r.URL.Query().Get("skip")
	skip, _ := utils.QueryConvertInt(skipQuery, 0)

	categories, err := h.Repos.CategoriesRepoImpl.SearchCategoryByName(q, skip, take)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
		panic(w)
	}
	count, err := h.Repos.CategoriesRepoImpl.Count(q)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	output, _ := json.MarshalIndent(categories, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 12)

	skipQuery := r.URL.Query().Get("skip")
	skip, _ := utils.QueryConvertInt(skipQuery, 0)

	categories, err := h.Repos.CategoriesRepoImpl.GetAll(take, skip)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
		panic(w)
	}
	count, err := h.Repos.CategoriesRepoImpl.Count()
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	output, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

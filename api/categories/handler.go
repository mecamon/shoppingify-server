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
	"strings"
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

// ShowAccount godoc
// @Summary      Creates a new category
// @Description  Creates a new category. It needs to have an unique name
// @Tags         categories
// @Param        user    body     models.CategoryDTO  true  "categorty info"
// @Accept       json
// @Produce      json
// @Success      200  {object} models.Created
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      409  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/categories [post]
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
		if strings.Contains(err.Error(), "unique constraint") {
			h.App.Loggers.Info.Println(err.Error())
			errMsg := appLocales.GetMsg("CategoryNameInUse", nil)
			errMap := models.ErrorMap{"error": errMsg}
			output, _ := json.MarshalIndent(errMap, "", "    ")
			utils.Response(w, http.StatusConflict, output)
			return
		} else {
			h.App.Loggers.Error.Println(err.Error())
			panic(w)
		}
	}
	_, err = h.Repos.TopCategoriesImpl.Add(userId, ID)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	res := models.Created{InsertedID: ID}
	output, _ := json.MarshalIndent(res, "", "    ")
	utils.Response(w, http.StatusCreated, output)
}

// ShowAccount godoc
// @Summary      Get categories by name
// @Description  Get categories by name. Pagination is able
// @Tags         categories
// @Param        q    query     string  true  "category search query"
// @Param        take    query     int  false  "items to take in query"
// @Param        skip    query     int  false  "items to skip in query"
// @Header       200              {string}  X-Total-Count  "total of items"
// @Accept       json
// @Produce      json
// @Success      200  {array} models.CategoryDTO
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/categories/by-name [get]
func (h *Handler) GetAllByName(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	q := r.URL.Query().Get("q")

	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 4)

	skipQuery := r.URL.Query().Get("skip")
	skip, _ := utils.QueryConvertInt(skipQuery, 0)

	categories, err := h.Repos.CategoriesRepoImpl.SearchCategoryByName(q, skip, take, userID)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
		panic(w)
	}
	count, err := h.Repos.CategoriesRepoImpl.Count(userID, q)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	output, _ := json.MarshalIndent(categories, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// ShowAccount godoc
// @Summary      Get categories
// @Description  Get categories. Pagination is able
// @Tags         categories
// @Param        take    query     int  false  "items to take in query"
// @Param        skip    query     int  false  "items to skip in query"
// @Header       200              {string}  X-Total-Count  "total of items"
// @Accept       json
// @Produce      json
// @Success      200  {array} models.CategoryDTO
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/categories [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)

	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 12)

	skipQuery := r.URL.Query().Get("skip")
	skip, _ := utils.QueryConvertInt(skipQuery, 0)

	categories, err := h.Repos.CategoriesRepoImpl.GetAll(take, skip, userID)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
		panic(w)
	}
	count, err := h.Repos.CategoriesRepoImpl.Count(userID)
	if err != nil {
		h.App.Loggers.Error.Println(err.Error())
	}

	output, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

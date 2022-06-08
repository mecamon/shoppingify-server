package items

import (
	"encoding/json"
	"fmt"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/services/storage"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	app   *config.App
	repos repositories.MainRepo
}

var handler *Handler

func InitHandler(conf *config.App) *Handler {
	handler = &Handler{app: conf, repos: repositories.Main}
	return handler
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	appLocales := appi18n.GetLocales(lang)
	ID := r.Context().Value("ID")

	err := r.ParseMultipartForm(128)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}
	categoryID, _ := strconv.ParseInt(r.Form.Get("category_id"), 10, 64)
	item := models.Item{
		Name:       r.Form.Get("name"),
		Note:       r.Form.Get("note"),
		CategoryID: categoryID,
	}

	var fileInfo FileInfo
	file, fileHeader, _ := r.FormFile("file")
	if file != nil {
		fileInfo = FileInfo{
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
		}
	}

	itemDom := ItemDom{
		appLocales: appLocales,
		item:       item,
		fileInfo:   fileInfo,
	}

	isValid, errMap := itemDom.validItem()
	if !isValid {
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	if file != nil {
		storageService, err := storage.GetStorage()
		if err != nil {
			h.app.Loggers.Error.Println(err.Error())
			log.Println(err.Error()) //TODO REMOVE THIS
		}

		imageURL, err := storageService.UploadImage(file, fmt.Sprintf("%s-%d", itemDom.item.Name, ID))
		if err != nil {
			h.app.Loggers.Error.Println(err.Error())
		}
		itemDom.item.ImageURL = imageURL
	}

	completedItem := itemDom.completedItemInfo()
	insertedID, err := h.repos.ItemsRepoIpml.Register(completedItem)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}

	res := map[string]interface{}{"insertedID": insertedID}
	output, _ := json.Marshal(res)
	utils.Response(w, http.StatusCreated, output)
}

func (h *Handler) GetByCategoryGroups(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	takeQuery := r.URL.Query().Get("take")
	take, _ := utils.QueryConvertInt(takeQuery, 4)

	skipQuery := r.URL.Query().Get("skip")
	skip, _ := utils.QueryConvertInt(skipQuery, 0)

	var categoriesGroups []models.CategoriesGroup
	var categories []models.CategoryDTO
	var err error
	var count int64

	if q != "" {
		categories, err = h.repos.CategoriesRepoImpl.GetAllWithItemName(q, take, skip)
		count, err = h.repos.CategoriesRepoImpl.Count(q)
	} else {
		categories, err = h.repos.CategoriesRepoImpl.GetAll(take, skip)
		count, err = h.repos.CategoriesRepoImpl.Count(q)
	}
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}

	for _, cat := range categories {
		cg := models.CategoriesGroup{
			CategoryID:   cat.ID,
			CategoryName: cat.Name,
		}

		items, err := h.repos.ItemsRepoIpml.GetAllByCategoryID(cat.ID)
		if err != nil {
			h.app.Loggers.Error.Println(err.Error())
			panic(w)
		}
		if len(items) > 0 {
			cg.Items = items
			categoriesGroups = append(categoriesGroups, cg)
		}
	}

	output, err := json.Marshal(categoriesGroups)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func (h *Handler) GetDetailsByID(w http.ResponseWriter, r *http.Request) {
	itemId := strings.TrimPrefix(r.URL.Path, "/api/items/")
	id, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	item, err := h.repos.ItemsRepoIpml.GetDetails(id)
	if err != nil {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	output, err := json.Marshal(item)
	if err != nil {
		utils.Response(w, http.StatusServiceUnavailable, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

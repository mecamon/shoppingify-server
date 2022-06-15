package lists

import (
	"encoding/json"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"net/http"
	"strings"
	"time"
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
	locales := appi18n.GetLocales(lang)
	userID := r.Context().Value("ID").(int64)
	body := struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	activeList, _ := h.repos.ListsRepoImpl.GetActive(userID)
	if activeList.ID != 0 {
		msg := locales.GetMsg("ThereIsAnActiveList", nil)
		errMap := models.ErrorMap{"list": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusConflict, output)
		return
	}

	if len(body.Name) == 0 {
		td := map[string]interface{}{"Field": "list name"}
		msg := locales.GetMsg("RequiredField", td)
		errMap := models.ErrorMap{"list": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	list := models.List{
		Name:        body.Name,
		IsCompleted: false,
		IsCancelled: false,
		UserID:      userID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	}

	insertedID, err := h.repos.ListsRepoImpl.Create(list)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	res := map[string]interface{}{"insertedID": insertedID}
	output, _ := json.Marshal(res)
	utils.Response(w, http.StatusCreated, output)
}

func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	list, err := h.repos.ListsRepoImpl.GetActive(userID)
	if err != nil {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}
	output, err := json.Marshal(list)
	utils.Response(w, http.StatusOK, output)
}

func (h *Handler) UpdateActiveListName(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	body := struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	if len(body.Name) == 0 {
		td := map[string]interface{}{"Field": "list name"}
		msg := locales.GetMsg("RequiredField", td)
		errMap := models.ErrorMap{"name": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	err = h.repos.ListsRepoImpl.UpdateActiveListName(userID, body.Name)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *Handler) AddItemToList(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	var item models.SelectedItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	dom := DomLists{
		appLocales: locales,
		itemToAdd:  item,
	}

	isValid, errMap := dom.validateItemToAdd()
	if !isValid {
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	completedItem := dom.completeItemToAdd()
	insertedID, err := h.repos.ListsRepoImpl.AddItemToList(completedItem)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		if strings.Contains(err.Error(), "inactive list") {
			msg := locales.GetMsg("NoOperationsOnInactiveList", nil)
			errMap := models.ErrorMap{"list": msg}
			output, _ := json.Marshal(errMap)
			utils.Response(w, http.StatusBadRequest, output)
		}

		return
	}

	res := map[string]interface{}{"insertedID": insertedID}
	output, _ := json.Marshal(res)
	utils.Response(w, http.StatusOK, output)
}

func (h *Handler) UpdateItemsSelected(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)
	var items []models.UpdateSelItemDTO

	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	dom := DomLists{
		appLocales:    locales,
		itemsToUpdate: items,
	}
	isValid, errMap := dom.validateItemToUpdate()
	if !isValid {
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	err = h.repos.ListsRepoImpl.UpdateItemsSelected(dom.itemsToUpdate)
	if err != nil {
		h.app.Loggers.Info.Println(err.Error())
		td := map[string]interface{}{"Item": err.Error()}
		msg := locales.GetMsg("DoesNotExist", td)
		errMap := models.ErrorMap{"item": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *Handler) DeleteItemFromList(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	body := struct {
		ItemSelID int64 `json:"item_sel_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	if body.ItemSelID == 0 {
		td := map[string]interface{}{"Field": "item_sel_id"}
		msg := locales.GetMsg("RequiredField", td)
		errMap := models.ErrorMap{"idRequired": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	err = h.repos.ListsRepoImpl.DeleteItemFromList(body.ItemSelID)
	if err != nil {
		h.app.Loggers.Info.Println(err.Error())
		td := map[string]interface{}{"Item": body.ItemSelID}
		msg := locales.GetMsg("DoesNotExist", td)
		errMap := models.ErrorMap{"item": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *Handler) CompleteItemSelected(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	body := struct {
		ItemSelID int64 `json:"item_sel_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusNotAcceptable, nil)
		return
	}

	if body.ItemSelID == 0 {
		td := map[string]interface{}{"Field": "item_sel_id"}
		msg := locales.GetMsg("RequiredField", td)
		errMap := models.ErrorMap{"idRequired": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	err = h.repos.ListsRepoImpl.CompleteItemSelected(body.ItemSelID)
	if err != nil {
		h.app.Loggers.Info.Println(err.Error())
		td := map[string]interface{}{"Item": err.Error()}
		msg := locales.GetMsg("DoesNotExist", td)
		errMap := models.ErrorMap{"item": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *Handler) CancelActive(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	err := h.repos.ListsRepoImpl.CancelActive(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		msg := locales.GetMsg("NoActiveList", nil)
		errMap := map[string]interface{}{"list": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *Handler) CompleteActive(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	err := h.repos.ListsRepoImpl.CompleteActive(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		msg := locales.GetMsg("NoActiveList", nil)
		errMap := map[string]interface{}{"list": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

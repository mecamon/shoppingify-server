package lists

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"net/http"
	"strconv"
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

// ShowAccount godoc
// @Summary      Creates a new list
// @Description  Creates a new list. Only one can be active at a time
// @Tags         lists
// @Param        list    body     models.CreateListDTO  true  "list info"
// @Accept       json
// @Produce      json
// @Success      200  {object} models.Created
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      409  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)
	userID := r.Context().Value("ID").(int64)
	body := models.CreateListDTO{}
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

// ShowAccount godoc
// @Summary      Get active list
// @Description  Get active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Success      200  {object} models.ListDTO
// @Failure      404  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/active [get]
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

// ShowAccount godoc
// @Summary      Update active list name
// @Description  Update active list name
// @Tags         lists
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/name [patch]
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

// ShowAccount godoc
// @Summary      Add item to active list
// @Description  Add item to active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @param selectedItem body models.AddSelectedItemDTO true "item to add to the active list"
// @Success      200 {object} models.Created
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/add-item [post]
func (h *Handler) AddItemToList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	var item models.AddSelectedItemDTO
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

	//Updating tops information
	err = h.repos.TopItemsImpl.Update(userID, completedItem.ItemID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
	}
	itemInfo, err := h.repos.ItemsRepoIpml.GetByID(completedItem.ItemID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
	}
	err = h.repos.TopCategoriesImpl.Update(userID, itemInfo.CategoryID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
	}
	res := models.Created{
		InsertedID: insertedID,
	}
	output, _ := json.Marshal(res)
	utils.Response(w, http.StatusOK, output)
}

// ShowAccount godoc
// @Summary      Update item in active list
// @Description  Update item in active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @param selectedItem body models.UpdateSelItemDTO true "update an item in the active list"
// @Success      200
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/update-items [patch]
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

// ShowAccount godoc
// @Summary      Delete item in active list
// @Description  Delete item in active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @param selectedItem body models.ItemSelIDDTO true "id from item to delete"
// @Success      200
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      404  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/selected-items [delete]
func (h *Handler) DeleteItemFromList(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)

	err = h.repos.ListsRepoImpl.DeleteItemFromList(itemID)
	if err != nil {
		h.app.Loggers.Info.Println(err.Error())
		td := map[string]interface{}{"Item": itemID}
		msg := locales.GetMsg("DoesNotExist", td)
		errMap := models.ErrorMap{"item": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusNotFound, output)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

// ShowAccount godoc
// @Summary      Completes a item in active list
// @Description  Completes a item in active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @param 		selectedItem body models.ItemSelIDDTO true "id from item to complete"
// @Success      200
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      404  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/selected-items [put]
func (h *Handler) CompleteItemSelected(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	var body models.ItemSelIDDTO
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

// ShowAccount godoc
// @Summary      Cancel the active list
// @Description  Cancel the active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      404  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/cancel-active [delete]
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

// ShowAccount godoc
// @Summary      Completes the active list
// @Description  Completes the active list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      404  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/complete-active [patch]
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

// ShowAccount godoc
// @Summary      Get old list
// @Description  Get old list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Success      200 {array} models.OldListDTO
// @Failure      500
// @Router       /api/lists/old-lists [get]
func (h *Handler) GetOldLists(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)
	_ = r.Header.Get("Accept-Language")

	oldLists, err := h.repos.ListsRepoImpl.GetOldOnes(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	output, err := json.Marshal(oldLists)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

// ShowAccount godoc
// @Summary      Get list by id
// @Description  Get list by id
// @Tags         lists
// @Accept       json
// @Param        listID    path     string  true  "list ID"
// @Produce      json
// @Success      200 {object} models.ListDTO
// @Failure      400  {object}  models.ErrorMapDTO
// @Failure      500
// @Router       /api/lists/{listId} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)

	userID := r.Context().Value("ID").(int64)
	listIDStr := chi.URLParam(r, "listID")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		msg := locales.GetMsg("InvalidRouteParam", nil)
		errMap := models.ErrorMap{"param": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	list, err := h.repos.ListsRepoImpl.GetByID(userID, listID)
	if err != nil {
		h.app.Loggers.Warning.Println(err.Error())
		td := map[string]interface{}{"Item": listID}
		msg := locales.GetMsg("DoesNotExist", td)
		errMap := models.ErrorMap{"listID": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	output, err := json.Marshal(list)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

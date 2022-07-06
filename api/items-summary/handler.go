package items_summary

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
)

type Handler struct {
	app   *config.App
	repos repositories.MainRepo
}

var handler *Handler

func InitHandler(app *config.App) *Handler {
	handler = &Handler{
		app:   app,
		repos: repositories.Main,
	}
	return handler
}

// ShowAccount godoc
// @Summary      Get summary group by month
// @Description  Get summary group by month
// @Tags         summary
// @Accept       json
// @Produce      json
// @Param        year    path     string  true  "year"
// @Success      200 {object} models.ItemsSummaryByMonthDTO
// @Success      204
// @Failure      500
// @Failure      400 {object} models.ErrorMap
// @Router       /api/summary/{year} [get]
func (h *Handler) GetByMonth(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	appLocales := appi18n.GetLocales(lang)
	userID := r.Context().Value("ID").(int64)
	yearStr := chi.URLParam(r, "year")
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		msg := appLocales.GetMsg("InvalidRouteParam", nil)
		errMap := models.ErrorMap{"param": msg}
		output, _ := json.Marshal(errMap)
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	itemsSumByMonth, err := h.repos.ItemsSummaryRepoImpl.GetMonthly(userID, year)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(itemsSumByMonth.Months) == 0 {
		utils.Response(w, http.StatusNoContent, nil)
		return
	}

	output, err := json.Marshal(itemsSumByMonth)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, output)
}

// ShowAccount godoc
// @Summary      Get summary group by year
// @Description  Get summary group by year
// @Tags         summary
// @Accept       json
// @Produce      json
// @Success      200 {array} models.ItemsSummaryByYearDTO
// @Success      204
// @Failure      500
// @Failure      400
// @Router       /api/summary [get]
func (h *Handler) GetByYear(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("ID").(int64)

	itemsByYear, err := h.repos.ItemsSummaryRepoImpl.GetYearly(userID)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(itemsByYear) == 0 {
		utils.Response(w, http.StatusNoContent, nil)
		return
	}

	output, err := json.Marshal(itemsByYear)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, output)
}

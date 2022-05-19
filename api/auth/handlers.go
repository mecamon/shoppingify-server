package auth

import (
	"encoding/json"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"net/http"
)

type Handler struct {
	Config *config.App
}

var handler *Handler

func InitHandler(conf *config.App) *Handler {
	handler = &Handler{Config: conf}
	return handler
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	valid, errMap := validCredentials(user, lang)
	if !valid {
		output, _ := json.MarshalIndent(errMap, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
	}
}

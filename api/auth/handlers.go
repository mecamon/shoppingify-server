package auth

import (
	"encoding/json"
	"github.com/mecamon/shoppingify-server/api/repositories"
	"github.com/mecamon/shoppingify-server/config"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
	"log"
	"net/http"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	valid, errMap := validCredentials(user, lang)
	if !valid {
		output, _ := json.MarshalIndent(errMap, "", "    ")
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	completedUser, err := completeUserInformation(user)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}

	id, err := h.repos.AuthRepoImpl.Register(completedUser)
	if err != nil {
		log.Println(err.Error())
		h.app.Loggers.Info.Println(err.Error())
		if strings.Contains(err.Error(), "unique constraint") {
			errMsg := locales.GetMsg("EmailAddressTaken", nil)
			errMap := models.ErrorMap{"error": errMsg}
			output, _ := json.MarshalIndent(errMap, "", "    ")
			utils.Response(w, http.StatusConflict, output)
			return
		} else {
			panic(w)
		}
	}

	token, err := json_web_token.Generate(id, user.Email)
	if err != nil {
		h.app.Loggers.Info.Println(err.Error())
		panic(w)
	}
	tokenMap := map[string]string{"token": token}
	output, _ := json.MarshalIndent(tokenMap, "", "    ")
	utils.Response(w, http.StatusCreated, output)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	locales := appi18n.GetLocales(lang)
	auth := models.Auth{}
	json.NewDecoder(r.Body).Decode(&auth)

	user, err := h.repos.AuthRepoImpl.SearchUserByEmail(auth.Email)
	if err != nil {
		h.app.Loggers.Warning.Println("wrong email")
		errMsg := locales.GetMsg("InvalidEmailOrPassword", nil)
		errMap := models.ErrorMap{"error": errMsg}
		output, _ := json.MarshalIndent(errMap, "", "    ")
		utils.Response(w, http.StatusBadRequest, output)
		return
	}
	rightCred, err := h.repos.AuthRepoImpl.CheckUserPassword(auth.Email, auth.Password)
	if !rightCred {
		h.app.Loggers.Warning.Println("wrong password")
	}
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
	}
	if !rightCred || err != nil {
		errMsg := locales.GetMsg("InvalidEmailOrPassword", nil)
		errMap := models.ErrorMap{"error": errMsg}
		output, _ := json.MarshalIndent(errMap, "", "    ")
		utils.Response(w, http.StatusBadRequest, output)
		return
	}

	token, err := json_web_token.Generate(user.ID, user.Email)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}
	tokenMap := map[string]string{"token": token}
	output, _ := json.MarshalIndent(tokenMap, "", "    ")
	utils.Response(w, http.StatusOK, output)
}

func (h *Handler) VisitorRegister(w http.ResponseWriter, r *http.Request) {
	visitor := createVisitorInformation()
	id, err := h.repos.AuthRepoImpl.Register(visitor)
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}

	token, err := json_web_token.Generate(id, "")
	if err != nil {
		h.app.Loggers.Error.Println(err.Error())
		panic(w)
	}
	tokenMap := map[string]string{"token": token}
	output, _ := json.MarshalIndent(tokenMap, "", "    ")
	utils.Response(w, http.StatusCreated, output)
}

package auth

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes() http.Handler {
	authRouter := chi.NewRouter()
	authRouter.Post("/register", handler.Register)
	authRouter.Post("/login", handler.Login)
	authRouter.Post("/visitor-register", handler.VisitorRegister)
	return authRouter
}

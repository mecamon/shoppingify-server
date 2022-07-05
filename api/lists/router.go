package lists

import (
	"github.com/go-chi/chi/v5"
	"github.com/mecamon/shoppingify-server/api/middlewares"
	"net/http"
)

func Routes() http.Handler {
	listsRouter := chi.NewRouter()
	listsRouter.Use(middlewares.TokenValidation)
	listsRouter.Post("/create", handler.Create)
	listsRouter.Get("/active", handler.GetActive)
	listsRouter.Patch("/name", handler.UpdateActiveListName)
	listsRouter.Post("/add-item", handler.AddItemToList)
	listsRouter.Patch("/update-items", handler.UpdateItemsSelected)
	listsRouter.Delete("/selected-items/{itemID}", handler.DeleteItemFromList)
	listsRouter.Put("/selected-items", handler.CompleteItemSelected)
	listsRouter.Delete("/cancel-active", handler.CancelActive)
	listsRouter.Patch("/complete-active", handler.CompleteActive)
	listsRouter.Get("/old-lists", handler.GetOldLists)
	listsRouter.Get("/{listID}", handler.GetByID)
	return listsRouter
}

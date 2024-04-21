package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/DanglingDynamo/chronotube/internal/handlers"
)

type Router struct {
	SearchHandler handlers.SearchHandler
}

func SetupRoutes(app *chi.Mux, router Router) {
	app.Get("/videos", router.SearchHandler.FetchVideosPaginated)
}

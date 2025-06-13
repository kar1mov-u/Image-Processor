package api

import (
	"image-processor/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type API struct {
	Routes http.Handler
	log    *log.Logger
	db     *database.Database
}

func New(log *log.Logger, db *database.Database) *API {
	api := &API{
		log: log,
		db:  db,
	}
	api.SetRoutes()
	return api
}

func (api *API) SetRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", api.login)
			r.Post("/signin", api.signIn)
		})
	})
	api.Routes = router
}

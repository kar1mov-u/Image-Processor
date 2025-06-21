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
	jwtKey string
}

func New(log *log.Logger, db *database.Database, jwtKey string) *API {
	api := &API{
		log:    log,
		db:     db,
		jwtKey: jwtKey,
	}
	api.SetRoutes()
	return api
}

func (api *API) SetRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(rAuth chi.Router) {
			rAuth.Post("/signin", api.signIn)
			rAuth.Post("/login", api.login)
		})

		r.Route("/images", func(r chi.Router) {
			r.Use(JwtMiddleware(api.jwtKey))
			r.Post("/", api.postImage)
			r.Get("/", api.getImages)
		})
	})

	api.Routes = router
}

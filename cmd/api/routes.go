package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Post("/register", app.registerUser)
	router.Post("/login", app.loginUser)

	router.Route("/v1/", func(r1 chi.Router) {
		r1.Get("/{url_key}", app.redirectToTargetUrl)

		r1.Route("/user/urls", func(r2 chi.Router) {
			r2.Use(app.authenticate)
			r2.Get("/", app.listUrls)
			r2.Post("/", app.createUrl)
			r2.Get("/{url_key}", app.getUrlInfo)
			r2.Delete("/{url_key}", app.deleteUrl)
		})
	})

	return router
}

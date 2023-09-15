package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func routes(app *application) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Post("/register", registerUser(app))
	router.Post("/login", loginUser(app))
	router.Get("/healtz", healthcheckHandler(app))

	router.Route("/v1/", func(r1 chi.Router) {
		r1.Get("/{url_key}", redirectToTargetUrl(app))

		r1.Route("/user/urls", func(r2 chi.Router) {
			r2.Use(authenticate(app))
			r2.Get("/{url_key}", getUrlInfo(app))
			r2.Delete("/{url_key}", deleteUrl(app))
			r2.Get("/", listUrls(app))
			r2.Post("/", createUrl(app))

		})
	})

	return router
}

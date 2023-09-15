package main

import (
	"errors"
	"net/http"
	"strings"
	"urlshortner/models"
)

func authenticate(app *application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader == "" {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			token := headerParts[1]
			user, err := app.models.Users.GetForToken(models.ScopeAuthentication, token)
			if err != nil {
				switch {
				case errors.Is(err, models.ErrRecordNotFound):
					app.invalidAuthenticationTokenResponse(w, r)
				default:
					app.serverErrorResponse(w, r, err)
				}
				return
			}
			r = app.contextSetUser(r, user)
			next.ServeHTTP(w, r)
		})
	}
}

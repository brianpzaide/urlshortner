package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"urlshortner/models"

	"github.com/go-chi/chi/v5"
)

const ttl = 24 * time.Hour

func redirectToTargetUrl(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		urlKey := chi.URLParam(r, "url_key")
		urlInfo, err := app.models.Urls.GetTargetUrl(urlKey, 0, false)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		http.Redirect(w, r, urlInfo.TargetUrl, http.StatusSeeOther)
	}
}

func listUrls(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user == nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		urls, err := app.models.Urls.ListUrls(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.writeJSON(w, http.StatusOK, envelope{"urls": urls}, nil)
	}
}

func createUrl(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user == nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		var input struct {
			TargetUrl string `json:"target_url"`
		}
		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		if strings.TrimSpace(input.TargetUrl) == "" {
			app.badRequestResponse(w, r, errors.New("empty target_url"))
			return
		}
		urlInfo := models.Url{
			TargetUrl: input.TargetUrl,
			UserId:    int64(user.ID),
		}
		err = app.models.Urls.Insert(&urlInfo)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.writeJSON(w, http.StatusCreated, envelope{"urlInfo": urlInfo}, nil)
	}
}

func getUrlInfo(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user == nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		urlKey := chi.URLParam(r, "url_key")

		urlInfo, err := app.models.Urls.GetTargetUrl(urlKey, int64(user.ID), true)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		app.writeJSON(w, http.StatusOK, envelope{"urlInfo": urlInfo}, nil)
	}
}

func deleteUrl(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user == nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		urlKey := chi.URLParam(r, "url_key")
		err := app.models.Urls.DeleteUrl(urlKey, int64(user.ID))
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "url successfully deleted"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func registerUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
			app.badRequestResponse(w, r, errors.New("empty email or password"))
			return
		}
		user := &models.User{
			Email: input.Email,
		}
		err = user.Password.Set(input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.models.Users.Insert(user)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrDuplicateEmail):
				app.badRequestResponse(w, r, errors.New("a user with this email address already exists"))
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func loginUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		user, err := app.models.Users.GetByEmail(input.Email)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.invalidCredentialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		match, err := user.Password.Matches(input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !match {
			app.invalidCredentialsResponse(w, r)
			return
		}
		token, err := app.models.Tokens.New(user.ID, ttl, models.ScopeAuthentication)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"urlshortner/models"
	"urlshortner/models/mock"

	"github.com/go-chi/chi/v5"
)

func Test_registerUser(t *testing.T) {

	// setting up mock database
	mock.DB = mock.MockDB{
		Users:       make(map[int64]*models.User),
		EmailLookup: map[string]int64{},
		Urls:        make(map[string]*models.Url),
		Tokens:      make(map[string]*models.Token),
	}
	mock.Count = 0

	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid json", `{"email":"abc@example.com","password":"123456"}`, http.StatusCreated},
		{"not json", `I'm not JSON`, http.StatusBadRequest},
		{"empty json", `{}`, http.StatusBadRequest},
		{"empty email", `{"email":""}`, http.StatusBadRequest},
		{"empty password", `{"email":"admin@example.com"}`, http.StatusBadRequest},
	}

	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/register", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.registerUser)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedStatusCode == http.StatusCreated {
			if tempId := mock.DB.EmailLookup["abc.example.com"]; tempId != 0 {
				t.Errorf("user not inserted into database")
			}
		}
	}
}

func Test_loginUser(t *testing.T) {

	// setting up mock database
	mock.DB = mock.MockDB{
		Users:       make(map[int64]*models.User),
		EmailLookup: map[string]int64{},
		Urls:        make(map[string]*models.Url),
		Tokens:      make(map[string]*models.Token),
	}
	mock.Count = 0

	testUser := models.User{
		Email: "abc@example.com",
	}
	testUser.Password.Set("123456")
	err := app.models.Users.Insert(&testUser)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}

	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"abc@example.com","password":"123456"}`, http.StatusCreated},
		{"not json", `I'm not JSON`, http.StatusBadRequest},
		{"empty json", `{}`, http.StatusUnauthorized},
		{"empty email", `{"email":""}`, http.StatusUnauthorized},
		{"empty password", `{"email":"abc@example.com"}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"abc@example.com","password":"secret"}`, http.StatusUnauthorized},
	}

	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req := httptest.NewRequest("POST", "/login", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.loginUser)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_redirectToTargetUrl(t *testing.T) {
	// setting up mock database
	mock.DB = mock.MockDB{
		Users:       make(map[int64]*models.User),
		EmailLookup: map[string]int64{},
		Urls:        make(map[string]*models.Url),
		Tokens:      make(map[string]*models.Token),
	}
	mock.Count = 0
	testUser := models.User{
		Email: "abc@example.com",
	}
	testUser.Password.Set("123456")
	err := app.models.Users.Insert(&testUser)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}

	urlInfo := models.Url{
		TargetUrl: "https://stackoverflow.com/",
		UserId:    int64(0),
	}
	app.models.Urls.Insert(&urlInfo)

	urlkeys := make([]string, 0)
	for url_key := range mock.DB.Urls {
		urlkeys = append(urlkeys, url_key)
	}

	router := chi.NewRouter()
	router.Get("/v1/{url_key}", app.redirectToTargetUrl)

	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
		expectedURL        string
	}{
		{"404", "/v1/fish", http.StatusNotFound, "/v1/fish"},
		{"stackOverflow", fmt.Sprintf("/v1/%s", urlkeys[0]), 303, "https://stackoverflow.com/"},
	}

	for _, e := range theTests {

		req := httptest.NewRequest("GET", e.url, nil)
		w := httptest.NewRecorder()

		// Call the handler with the request
		router.ServeHTTP(w, req)

		// Check the response
		resp := w.Result()
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("Expected status code %d, but got %d", e.expectedStatusCode, resp.StatusCode)
		} else {
			loc, err := resp.Location()
			if err == nil {
				if loc.String() != e.expectedURL {
					t.Errorf("redirected to %s instead of %s", loc.String(), e.expectedURL)
				}
			}
			if url, ok := mock.DB.Urls[e.url]; ok {
				if url.Visits == 0 {
					t.Errorf("number of visits must increase by 1")
				}
			}

		}
	}
}

func Test_apiHandlers(t *testing.T) {
	// setting up mock database
	mock.DB = mock.MockDB{
		Users:       make(map[int64]*models.User),
		EmailLookup: map[string]int64{},
		Urls:        make(map[string]*models.Url),
		Tokens:      make(map[string]*models.Token),
	}
	mock.Count = 0
	userTokens := make([]string, 0)
	testUser1 := models.User{
		Email: "abc@example.com",
	}
	testUser1.Password.Set("123456")
	err := app.models.Users.Insert(&testUser1)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	token, err := app.models.Tokens.New(0, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error logging the user", err)
	}
	userTokens = append(userTokens, token.Plaintext)

	urlInfo := models.Url{
		TargetUrl: "https://stackoverflow.com/",
		UserId:    int64(0),
	}
	app.models.Urls.Insert(&urlInfo)

	testUser2 := models.User{
		Email: "123@example.com",
	}
	testUser2.Password.Set("secret")
	err = app.models.Users.Insert(&testUser2)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	token, err = app.models.Tokens.New(1, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error logging the user", err)
	}
	userTokens = append(userTokens, token.Plaintext)

	urlkeys := make([]string, 0)
	for url_key := range mock.DB.Urls {
		urlkeys = append(urlkeys, url_key)
	}

	var tests = []struct {
		name           string
		userToken      string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"listurls", userTokens[0], "GET", "", "", app.listUrls, http.StatusOK},
		{"listurlsNotAuthenticated", "", "GET", "", "", app.listUrls, http.StatusUnauthorized},
		{"getUrlNotAuthenticated", "", "GET", "", urlkeys[0], app.getUrlInfo, http.StatusUnauthorized},
		{"getUrl", userTokens[0], "GET", "", urlkeys[0], app.getUrlInfo, http.StatusOK},
		{"getUrlNotYetCreatedByUser", userTokens[0], "GET", "", "fish", app.getUrlInfo, http.StatusNotFound},
		{"getUrlCreatedByOtherUser", userTokens[1], "GET", "", urlkeys[0], app.getUrlInfo, http.StatusNotFound},
		{"deleteUrlNotAuthenticated", "", "DELETE", "", urlkeys[0], app.deleteUrl, http.StatusUnauthorized},
		{"deleteUrl", userTokens[0], "DELETE", "", urlkeys[0], app.deleteUrl, http.StatusOK},
		{"deleteUrlNotCreatedByUser", userTokens[1], "DELETE", "", urlkeys[0], app.deleteUrl, http.StatusNotFound},
		{"createUrlNotAuthenticated", "", "POST", `{"target_url": "https://stackoverflow.com"}`, "", app.createUrl, http.StatusUnauthorized},
		{"createUrl", userTokens[0], "POST", `{"target_url": "https://stackoverflow.com"}`, "", app.createUrl, http.StatusCreated},
		{"createUrlBadJSON", userTokens[0], "POST", `{}`, "", app.createUrl, http.StatusBadRequest},
		{"createUrlNotJSON", userTokens[0], "POST", `not json`, "", app.createUrl, http.StatusBadRequest},
	}

	router := chi.NewRouter()
	router.Route("/v1/", func(r1 chi.Router) {
		r1.Route("/user/urls", func(r2 chi.Router) {
			r2.Use(app.authenticate)
			r2.Get("/", app.listUrls)
			r2.Post("/", app.createUrl)
			r2.Get("/{url_key}", app.getUrlInfo)
			r2.Delete("/{url_key}", app.deleteUrl)
		})
	})

	for _, e := range tests {
		var req *http.Request

		reqUrl := "/v1/user/urls"

		if e.paramID != "" {
			reqUrl = fmt.Sprintf("/v1/user/urls/%s", e.paramID)
		}

		if e.json == "" {
			req = httptest.NewRequest(e.method, reqUrl, nil)
		} else {
			req = httptest.NewRequest(e.method, reqUrl, strings.NewReader(e.json))
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.userToken))

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

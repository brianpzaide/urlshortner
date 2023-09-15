package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"urlshortner/models"

	"github.com/go-chi/chi/v5"
)

func Test_registerUser(t *testing.T) {

	testApp, testdb := setup(t)

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
		handler := http.HandlerFunc(registerUser(testApp))

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedStatusCode == http.StatusCreated {
			tempId, ok := testdb.EmailLookup["abc@example.com"]
			if !ok {
				t.Errorf("user not inserted into mock database")
			} else {
				tempUser, ok := testdb.Users[tempId]
				if !ok {
					t.Errorf("user not inserted into mock database")
				} else {
					log.Println(tempUser.Email, tempUser.ID, "got created")
				}
			}
		}
	}
}

func Test_loginUser(t *testing.T) {

	testApp, testdb := setup(t)

	testUser := models.User{
		Email: "abc@example.com",
	}
	testUser.Password.Set("123456")
	err := testApp.models.Users.Insert(&testUser)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	tempUsr, ok := testdb.Users[int64(0)]
	if !ok || tempUsr.Email != testUser.Email {
		t.Fatal("error inserting user into mock database")
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
		handler := http.HandlerFunc(loginUser(testApp))

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_redirectToTargetUrl(t *testing.T) {
	testApp, testdb := setup(t)
	testUser := models.User{
		Email: "abc@example.com",
	}
	testUser.Password.Set("123456")
	err := testApp.models.Users.Insert(&testUser)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	tempUsr, ok := testdb.Users[int64(0)]
	if !ok || tempUsr.Email != testUser.Email {
		t.Fatal("error inserting user into mock database")
	}

	urlInfo := models.Url{
		TargetUrl: "https://stackoverflow.com/",
		UserId:    int64(0),
	}
	err = testApp.models.Urls.Insert(&urlInfo)
	if err != nil {
		t.Fatal("error inserting user into mock database")
	}

	urlkeys := make([]string, 0)
	for url_key := range testdb.Urls {
		urlkeys = append(urlkeys, url_key)
	}

	router := chi.NewRouter()
	router.Get("/v1/{url_key}", redirectToTargetUrl(testApp))

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
				tempUrl, err := testApp.models.Urls.GetTargetUrl(urlkeys[0], int64(0), true)
				if err != nil {
					t.Error(err)
				}
				if tempUrl == nil || tempUrl.Visits == 0 {
					t.Errorf("case: %s tempUrl is nil (or) number of visits must increase by 1", e.url)
				}
			}

		}
	}
}

func Test_apiHandlers(t *testing.T) {
	testApp, testdb := setup(t)
	userTokens := make([]string, 0)
	testUser1 := models.User{
		Email: "abc@example.com",
	}
	testUser1.Password.Set("123456")
	err := testApp.models.Users.Insert(&testUser1)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	tempUsr, ok := testdb.Users[int64(0)]
	if !ok || tempUsr.Email != testUser1.Email {
		t.Fatal("error inserting user into mock database")
	}
	token1, err := testApp.models.Tokens.New(tempUsr.ID, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error logging the user", err)
	}
	tempToken1, ok := testdb.Tokens[token1.Plaintext]
	if !ok || tempToken1.Plaintext != token1.Plaintext {
		t.Fatal("error token not inserted into mock database for testUser1")
	}
	userTokens = append(userTokens, token1.Plaintext)

	urlInfo := models.Url{
		TargetUrl: "https://stackoverflow.com/",
		UserId:    tempUsr.ID,
	}
	testApp.models.Urls.Insert(&urlInfo)

	testUser2 := models.User{
		Email: "123@example.com",
	}
	testUser2.Password.Set("secret")
	err = testApp.models.Users.Insert(&testUser2)
	if err != nil {
		t.Fatal("error inserting new user", err)
	}
	tempUsr2, ok := testdb.Users[int64(1)]
	if !ok || tempUsr2.Email != testUser2.Email {
		t.Fatal("error inserting user into mock database")
	}
	token2, err := testApp.models.Tokens.New(tempUsr2.ID, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error logging the user", err)
	}
	tempToken2, ok := testdb.Tokens[token2.Plaintext]
	if !ok || tempToken2.Plaintext != token2.Plaintext || tempToken2.UserId != token2.UserId {
		t.Fatal("error token not inserted into mock database for testUser2")
	}
	userTokens = append(userTokens, token2.Plaintext)

	urlkeys := make([]string, 0)
	for url := range testdb.Urls {
		urlkeys = append(urlkeys, url)
	}

	var tests = []struct {
		name           string
		userToken      string
		method         string
		json           string
		paramID        string
		expectedStatus int
	}{
		{"listurls", userTokens[0], "GET", "", "", http.StatusOK},
		{"listurlsNotAuthenticated", "", "GET", "", "", http.StatusUnauthorized},
		{"getUrlNotAuthenticated", "", "GET", "", urlkeys[0], http.StatusUnauthorized},
		{"getUrl", userTokens[0], "GET", "", urlkeys[0], http.StatusOK},
		{"getUrlNotYetCreatedByUser", userTokens[0], "GET", "", "fish", http.StatusNotFound},
		{"getUrlCreatedByOtherUser", userTokens[1], "GET", "", urlkeys[0], http.StatusNotFound},
		{"deleteUrlNotAuthenticated", "", "DELETE", "", urlkeys[0], http.StatusUnauthorized},
		{"deleteUrl", userTokens[0], "DELETE", "", urlkeys[0], http.StatusOK},
		{"deleteUrlNotCreatedByUser", userTokens[1], "DELETE", "", urlkeys[0], http.StatusNotFound},
		{"createUrlNotAuthenticated", "", "POST", `{"target_url": "https://stackoverflow.com"}`, "", http.StatusUnauthorized},
		{"createUrl", userTokens[0], "POST", `{"target_url": "https://stackoverflow.com"}`, "", http.StatusCreated},
		{"createUrlBadJSON", userTokens[0], "POST", `{}`, "", http.StatusBadRequest},
		{"createUrlNotJSON", userTokens[0], "POST", `not json`, "", http.StatusBadRequest},
	}

	router := chi.NewRouter()
	router.Route("/v1/user/urls", func(r1 chi.Router) {
		r1.Use(authenticate(testApp))
		r1.Get("/{url_key}", getUrlInfo(testApp))
		r1.Delete("/{url_key}", deleteUrl(testApp))
		r1.Get("/", listUrls(testApp))
		r1.Post("/", createUrl(testApp))
	})

	for _, e := range tests {
		var req *http.Request

		reqUrl := "/v1/user/urls"

		if e.paramID != "" {
			reqUrl = fmt.Sprintf("/v1/user/urls/%s", e.paramID)
		}
		log.Printf("testname: %s url: %s", e.name, reqUrl)

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

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"urlshortner/models"
	"urlshortner/models/mock"
)

func Test_authenticate(t *testing.T) {

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

	tokenValid, err := app.models.Tokens.New(testUser.ID, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error generating token", err)
	}

	tokenExpired, err := models.GenerateToken(testUser.ID, 0*time.Hour, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error generating token", err)
	}
	tokenExpired.Expiry = time.Now().Add(time.Duration(-1) * time.Hour)
	err = app.models.Tokens.Insert(tokenExpired)
	if err != nil {
		t.Fatal("error inserting token", err)
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	var tests = []struct {
		name             string
		token            string
		expectAuthorized bool
		setHeader        bool
	}{
		{name: "valid token", token: fmt.Sprintf("Bearer %s", tokenValid.Plaintext), expectAuthorized: true, setHeader: true},
		{name: "no token", token: "", expectAuthorized: false, setHeader: false},
		{name: "invalid token", token: fmt.Sprintf("Bearer %s", tokenExpired.Plaintext), expectAuthorized: false, setHeader: true},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}
		rr := httptest.NewRecorder()

		handlerToTest := app.authenticate(nextHandler)
		handlerToTest.ServeHTTP(rr, req)

		if e.expectAuthorized && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got code 401, and should not have", e.name)
		}

		if !e.expectAuthorized && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get code 401, and should have", e.name)
		}
	}
}

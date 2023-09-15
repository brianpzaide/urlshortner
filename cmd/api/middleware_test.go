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

func setup(t *testing.T) (*application, *mock.MockDB) {
	testApp := &application{}
	// setting up mock database
	testdb := &mock.MockDB{
		Users:       make(map[int64]*models.User),
		EmailLookup: make(map[string]int64),
		Urls:        make(map[string]*models.Url),
		Tokens:      make(map[string]*models.Token),
	}

	testApp.models = models.Models{
		Users:  mock.NewUserModel(testdb),
		Urls:   mock.NewUrlModel(testdb),
		Tokens: mock.NewTokenModel(testdb),
	}
	return testApp, testdb
}

func Test_authenticate(t *testing.T) {

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

	tokenValid, err := testApp.models.Tokens.New(testUser.ID, ttl, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error generating token", err)
	}
	tempToken, ok := testdb.Tokens[tokenValid.Plaintext]
	if !ok || tempToken.Plaintext != tokenValid.Plaintext {
		t.Fatal("error inserting token into mock database")
	}

	tokenExpired, err := models.GenerateToken(testUser.ID, 0*time.Hour, models.ScopeAuthentication)
	if err != nil {
		t.Fatal("error generating token", err)
	}
	tokenExpired.Expiry = time.Now().Add(time.Duration(-1) * time.Hour)
	err = testApp.models.Tokens.Insert(tokenExpired)
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

		handlerToTest := authenticate(testApp)(nextHandler)
		handlerToTest.ServeHTTP(rr, req)

		if e.expectAuthorized && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got code 401, and should not have", e.name)
		}

		if !e.expectAuthorized && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get code 401, and should have", e.name)
		}

	}
}

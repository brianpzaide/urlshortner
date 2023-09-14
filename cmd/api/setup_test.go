package main

import (
	"log"
	"os"
	"testing"
	"urlshortner/models"
	"urlshortner/models/mock"
)

var app application

func TestMain(m *testing.M) {
	app.models = models.Models{
		Users:  mock.UserModel{},
		Urls:   mock.UrlModel{},
		Tokens: mock.TokenModel{},
	}
	app.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	os.Exit(m.Run())
}

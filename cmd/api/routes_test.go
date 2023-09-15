package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type tempRoute struct {
	method string
	route  string
}

func Test_app_routes(t *testing.T) {
	testApp := application{}

	var registered = []struct {
		route  string
		method string
	}{
		{"/healtz", "GET"},
		{"/register", "POST"},
		{"/login", "POST"},
		{"/v1/{url_key}", "GET"},
		{"/v1/user/urls/", "GET"},
		{"/v1/user/urls/", "POST"},
		{"/v1/user/urls/{url_key}", "GET"},
		{"/v1/user/urls/{url_key}", "DELETE"},
	}

	chiRoutes := routes(&testApp).(chi.Routes)
	for _, route := range registered {
		// check to see if the route exists
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}

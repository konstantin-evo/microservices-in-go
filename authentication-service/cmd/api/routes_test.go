package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes_exit(t *testing.T) {
	testApp := Config{}

	testRouter, ok := testApp.routes().(chi.Router)
	if !ok {
		t.Errorf("testRouter is not of type chi.Router")
		return
	}

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, testRouter, route)
	}
}

func routeExists(t *testing.T, router chi.Router, route string) {
	found := false

	_ = chi.Walk(router, func(method, currentRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == currentRoute {
			found = true
		}
		return nil
	})

	if (!found) {
		t.Errorf("didn't find %s in registered routes", route)
	}
}

package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/onlytunesradio/go-api-template/src/api/test"
)

func TestRouter() chi.Router {

	// New Chi SubRouter
	testRoute := chi.NewRouter()

	// Set up sub-routes
	testRoute.Get("/get", test.Get)
	testRoute.Post("/post", test.Post)
	testRoute.Delete("/delete", test.Delete)

	// Return the Sub-Route back to the main API Router in main.go
	return testRoute
}

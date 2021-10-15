package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// GetAPI returns an http.Handler configured with all API routes
func GetAPI() http.Handler {
	router := mux.NewRouter()

	router.
		Methods(http.MethodGet).
		Path("/users").
		HandlerFunc(GetUsers)

	router.
		Methods(http.MethodGet).
		Path("/users/{id}").
		HandlerFunc(GetUser)

	return router
}


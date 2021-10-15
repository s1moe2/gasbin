package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var usersRepo map[string]User

// GetUsers responds with the full list of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if usersRepo == nil {
		RespondJSON(w, map[string]User{}, http.StatusOK)
		return
	}
	RespondJSON(w, usersRepo, http.StatusOK)
}

// GetUser responds with a single user object if found
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		RespondError(w, NewError("invalid param id"), http.StatusBadRequest)
		return
	}

	user, found := usersRepo[userID]
	if !found {
		RespondError(w, NewError("user not found"), http.StatusNotFound)
		return
	}
	RespondJSON(w, user, http.StatusOK)
}
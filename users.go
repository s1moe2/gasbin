package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleBasic      = "basic"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type NewUserPayload struct {
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type usersHandler struct {
	usersRepo map[string]User
}

func NewUsersHandler() *usersHandler {
	users := map[string]User{}

	u1 := uuid.New().String()
	users[u1] = User{
		ID:   u1,
		Name: "Niles",
		Role: RoleAdmin,
	}

	u2 := uuid.New().String()
	users[u2] = User{
		ID:   u2,
		Name: "Mary",
		Role: RoleBasic,
	}

	return &usersHandler{
		usersRepo: users,
	}
}

// GetUsers responds with the full list of users
func (uh *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if uh.usersRepo == nil {
		RespondJSON(w, map[string]User{}, http.StatusOK)
		return
	}
	RespondJSON(w, uh.usersRepo, http.StatusOK)
}

// GetUser responds with a single user object if found
func (uh *usersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		RespondError(w, NewError("invalid param id"), http.StatusBadRequest)
		return
	}

	user, found := uh.usersRepo[userID]
	if !found {
		RespondError(w, NewError("user not found"), http.StatusNotFound)
		return
	}
	RespondJSON(w, user, http.StatusOK)
}

func (uh *usersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		RespondError(w, NewError("invalid param id"), http.StatusBadRequest)
		return
	}

	var newUser NewUserPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		RespondError(w, NewError("invalid payload"), http.StatusBadRequest)
		return
	}

	user, found := uh.usersRepo[userID]
	if !found {
		RespondError(w, NewError("user not found"), http.StatusNotFound)
		return
	}
	user.Name = newUser.Name
	uh.usersRepo[user.ID] = user

	RespondJSON(w, user, http.StatusOK)
}

func (uh *usersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser NewUserPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		RespondError(w, NewError("invalid payload"), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	user := User{
		ID:   id,
		Name: newUser.Name,
		Role: newUser.Role,
	}
	uh.usersRepo[id] = user

	RespondJSON(w, user, http.StatusCreated)
}

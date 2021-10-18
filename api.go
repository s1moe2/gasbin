package main

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func VerbMatch(req string, policy string) bool {
	if len(policy) == 1 && policy == "*" {
		return true
	}

	acceptedVerbs := strings.Split(policy, "|")
	for _, v := range acceptedVerbs {
		if v == req {
			return true
		}
	}
	return false
}

func VerbMatchFunc(args ...interface{}) (interface{}, error) {
	arg1 := args[0].(string)
	arg2 := args[1].(string)

	return VerbMatch(arg1, arg2), nil
}

func IsOwner(req []string, policy string) bool {

	return true
}

type userAccess struct {
	id string
}

func IsOwnerFunc(args ...interface{}) (interface{}, error) {
	arg1 := args[0].([]string)
	arg2 := args[1].(string)

	return IsOwner(arg1, arg2), nil
}

// GetAPI returns an http.Handler configured with all API routes
func GetAPI() http.Handler {
	router := mux.NewRouter()

	uh := NewUsersHandler()

	enforcer, _ := casbin.NewEnforcer("./model.conf", "./policy.csv")
	enforcer.AddFunction("verbMatch", VerbMatchFunc)
	enforcer.AddFunction("isOwner", IsOwnerFunc)

	router.
		Methods(http.MethodGet).
		Path("/users").
		HandlerFunc(auth(uh.GetUsers, enforcer))

	router.
		Methods(http.MethodGet).
		Path("/users/{id}").
		HandlerFunc(auth(uh.GetUser, enforcer))

	router.
		Methods(http.MethodPut).
		Path("/users/{id}").
		HandlerFunc(auth(uh.UpdateUser, enforcer))

	router.
		Methods(http.MethodPost).
		Path("/users").
		HandlerFunc(auth(uh.CreateUser, enforcer))

	return router
}


func auth(next http.HandlerFunc, e *casbin.Enforcer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			RespondError(w, NewError("unauthorized"), http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(token, ":")
		userID := tokenParts[0]
		role := tokenParts[1]
		orgs := strings.Split(tokenParts[2], ",")

		allowed, err := e.Enforce(role, r.URL.Path, r.Method, orgs)
		if err != nil {
			RespondError(w, NewError(err), http.StatusInternalServerError)
			return
		}

		if !allowed {
			RespondError(w, NewError("forbidden"), http.StatusForbidden)
			return
		}

		newReq := r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next.ServeHTTP(w, newReq)
	}
}
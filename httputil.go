package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type internalError struct {
	Message   string
	RequestID string
}

type UserError struct {
	Errors []string
}

type Errors []error

func NewError(errs interface{}) Errors {
	switch parsed := errs.(type) {
	case []string:
		var res Errors
		for _, err := range parsed {
			e := err
			res = append(res, errors.New(e))
		}
		return res
	case string:
		return []error{errors.New(parsed)}
	case error:
		return []error{parsed}
	case []error:
		return parsed
	default:
		log.Errorf("unknown error type %v", errs)
		return nil
	}
}

// RespondJSON is an helper that takes care of constructing and sending a data JSON payload
func RespondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	response := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("RespondJSON", err)
	}
}

// RespondError is an helper that takes care of the
// HTTP response part of a request handler for when the response is not on the 200 range
func RespondError(w http.ResponseWriter, errors Errors, statusCode int) {
	var errMessages []string
	for _, e := range errors {
		errMessages = append(errMessages, e.Error())
	}

	response := struct {
		Errors []string `json:"errors"`
	}{
		Errors: errMessages,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("RespondError", err)
	}
}
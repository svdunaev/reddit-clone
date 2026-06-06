package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func WriteError(w http.ResponseWriter, status int, code, msg string, details any) {
	httpError := HttpError{
		Code:    code,
		Message: msg,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(httpError)
	if err != nil {
		log.Panic(err)
	}
}

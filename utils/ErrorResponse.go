package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewErrorResponse(w http.ResponseWriter, status int, message string, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := ErrorResponse{Status: status, Message: message, Error: err}
	json.NewEncoder(w).Encode(res)
}

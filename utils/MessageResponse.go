package utils

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewMessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := MessageResponse{Status: status, Message: message}
	json.NewEncoder(w).Encode(res)
}

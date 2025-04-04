package utils

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewMessageResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := MessageResponse{Status: status, Message: message, Data: data}
	json.NewEncoder(w).Encode(res)
}
package httputil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func New(w http.ResponseWriter) *ResponseBuilder {
	return &ResponseBuilder{
		w: w,
	}
}

type ResponseBuilder struct {
	w       http.ResponseWriter
	status  int
	message string
	data    any
}

func (rb *ResponseBuilder) Status(status int) *ResponseBuilder {
	rb.status = status
	return rb
}

func (rb *ResponseBuilder) Message(msg string) *ResponseBuilder {
	rb.message = msg
	return rb
}

func (rb *ResponseBuilder) Data(data any) *ResponseBuilder {
	rb.data = data
	return rb
}

func (rb *ResponseBuilder) Send() {
	rb.w.Header().Set("Content-Type", "application/json")
	rb.w.WriteHeader(rb.status)

	json.NewEncoder(rb.w).Encode(Response{
		Status:  rb.status,
		Message: rb.message,
		Data:    rb.data,
	})
}

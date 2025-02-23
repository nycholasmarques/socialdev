package handlers

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parsear os dados da requisição
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Chamar o serviço para criar o usuário
	_, err := h.userService.CreateUser(username, email, password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
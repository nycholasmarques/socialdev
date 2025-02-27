package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserWithUsername struct {
	Username string
}

var validate *validator.Validate

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	validate = validator.New()
	// Parsear os dados da requisição
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validar os dados
	err := validate.Var(username, "required,min=3,max=30,alpha")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid username: %v", err), http.StatusBadRequest)
		return
	}

	err = validate.Var(email, "required,email")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid email: %v", err), http.StatusBadRequest)
		return
	}

	err = validate.Var(password, "required,min=6,max=12")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid password: %v", err), http.StatusBadRequest)
		return
	}

	// validar erro se email já existe

	// Chamar o serviço para criar o usuário
	_, err = h.userService.CreateUser(username, email, password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Path
	userID := strings.TrimPrefix(user_id, "/user/")

	uuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}
	
	user, err := h.userService.GetUser(uuid)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) GetUserWithUsernameHandler(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	validate = validator.New()

	err := validate.Var(username, "required,max=30")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid username: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserWithUsername(username)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		return
	}
		// Retornar a resposta
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
}
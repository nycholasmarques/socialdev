package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

type UpdateUserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Github   string `json:"github"`
	Linkedin string `json:"linkedin"`
	Website  string `json:"website"`
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

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	
	id := r.URL.Path
	userID := strings.TrimPrefix(id, "/user/update/")

	convertUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(convertUUID)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	rb.Close()

	var userUpdate *UpdateUserRequest
	err = json.Unmarshal(reader, &userUpdate)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	validate = validator.New()

	if userUpdate.Username == "" {
		userUpdate.Username = user.Username
	}

	if userUpdate.Email == "" {
		userUpdate.Email = user.Email
	}

	if userUpdate.Password == "" {
		userUpdate.Password = user.Password
	}

	err = validate.Var(userUpdate.Username, "min=3,max=30,alpha")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid username: %v", err), http.StatusBadRequest)
		return
	}

	err = validate.Var(userUpdate.Email, "email")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid email: %v", err), http.StatusBadRequest)
		return
	}

	err = validate.Var(userUpdate.Password, "min=6,max=12")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid password: %v", err), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.UpdateUser(convertUUID, userUpdate.Username, userUpdate.Email, userUpdate.Password, userUpdate.Avatar, userUpdate.Bio, userUpdate.Github, userUpdate.Linkedin, userUpdate.Website)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

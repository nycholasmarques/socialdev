package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/services"
	"github.com/nycholasmarques/socialdev/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cfg = configs.LoadConfig()

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserWithUsername struct {
	Username string
}

type LoginUserRequest struct {
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=6,max=12"`
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

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	User entity.UserResponse `json:"user"`
	Token string               `json:"token"`
}

var validate *validator.Validate

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  cfg.GoogleRedirectURL,
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	randomState = "random" 
)

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	rb.Close()

	var registerUser *RegisterUserRequest
	err = json.Unmarshal(reader, &registerUser)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	validate = validator.New()
	// Parsear os dados da requisição
	username := registerUser.Username
	email := registerUser.Email
	password := registerUser.Password

	// Validar os dados
	err = validate.Var(username, "required,min=3,max=30,alpha")
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
	
	// Chamar o serviço para criar o usuário
	userResponse, err := h.userService.CreateUser(username, email, password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := configs.CreateToken(userResponse.Username)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		User: userResponse,
		Token: token,
	}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Path
	userID := strings.TrimPrefix(user_id, "/users/")

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
	userID := strings.TrimPrefix(id, "/users/")

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

	var userUpdate UpdateUserRequest
	err = json.Unmarshal(reader, &userUpdate)
	if err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			return
	}

	validate := validator.New()

	// Atualiza apenas os campos fornecidos
	if userUpdate.Username != "" {
			err = validate.Var(userUpdate.Username, "min=3,max=30")
			if err != nil {
					http.Error(w, fmt.Sprintf("Invalid username: %v", err), http.StatusBadRequest)
					return
			}
	} else {
			userUpdate.Username = user.Username
	}

	if userUpdate.Email != "" {
			err = validate.Var(userUpdate.Email, "email")
			if err != nil {
					http.Error(w, fmt.Sprintf("Invalid email: %v", err), http.StatusBadRequest)
					return
			}
	} else {
			userUpdate.Email = user.Email
	}

	if userUpdate.Password != "" {
			err = validate.Var(userUpdate.Password, "min=6,max=12")
			if err != nil {
					http.Error(w, fmt.Sprintf("Invalid password: %v", err), http.StatusBadRequest)
					return
			}
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
				return
			}
			userUpdate.Password = string(passwordHash)
	} else {
			userUpdate.Password = user.Password
	}

	// Atualiza o usuário com os novos valores
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

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Path
	userID := strings.TrimPrefix(user_id, "/users/")

	uuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	_, err = h.userService.GetUser(uuid)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	err = h.userService.DeleteUser(uuid)
	if err != nil {
		http.Error(w, "Failed delete user", http.StatusInternalServerError)
		return
	}

	// Retornar a resposta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully!"))
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	rb.Close()
	fmt.Println("ok")
	var userLogin *LoginUserRequest
	err = json.Unmarshal(reader, &userLogin)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	validate = validator.New()

	err = validate.Var(userLogin.Email, "required,email")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid email: %v", err), http.StatusBadRequest)
		return
	}

	err = validate.Var(userLogin.Password, "required,min=6,max=12")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid password: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.userService.LoginUser(userLogin.Email)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := configs.CreateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		User: user,
		Token: token,
	}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	// Retornar a posta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *UserHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, "Método não permitido", "")
		return
	}

	if r.FormValue("state") != randomState {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Erro na autenticação", "")
		return
	}

	context := context.Background()

	token, err := googleOauthConfig.Exchange(context, r.FormValue("code"))
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao obter token", err.Error())
		return
	}

	client := googleOauthConfig.Client(context, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao obter informações do usuário", err.Error())
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao decodificar informações do usuário", "")
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro na autenticação", "")
		return
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro na autenticação", "")
		return
	}

	user, err := h.userService.LoginUser(email)
	if err != nil {
		user, err = h.userService.CreateUserWithGoogle(name, email, "", userInfo["id"].(string), userInfo["picture"].(string))
		if err != nil {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao criar usuário", err.Error())
			return
		}
	}

	jwtToken, err := configs.CreateToken(user.Username)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao criar token", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": jwtToken,
		"user":  user,
	})
}
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
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

type UserResponse struct {
	User  entity.UserResponse `json:"user"`
	Token string              `json:"token"`
}

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

	var registerUser entity.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&registerUser); err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate = validator.New()

	if err := validate.Struct(registerUser); err != nil {
		errorMessages := entity.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	userResponse, err := h.userService.CreateUser(registerUser.Username, registerUser.Email, registerUser.Password)
	if err != nil {
		switch err {
		case entity.ErrEmailAlreadyExists:
			utils.NewErrorResponse(w, http.StatusConflict, "Email already exists", err.Error())
		case entity.ErrUsernameAlreadyExists:
			utils.NewErrorResponse(w, http.StatusConflict, "Username already exists", err.Error())
		default:
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		}
		return
	}

	token, err := configs.CreateToken(userResponse.UserID, true)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	response := UserResponse{
		User:  userResponse,
		Token: token,
	}

	res, err := json.Marshal(response)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal response", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Path
	userID := strings.TrimPrefix(user_id, "/users/")

	uuid, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUser(uuid)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			utils.NewErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		default:
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get user", err.Error())
		}
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal user", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	username := r.URL.Query().Get("username")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	offset := int32(0)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	users, err := h.userService.GetAllUsers(username, limit, offset)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal users", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	convertUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUser(convertUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	var userUpdate entity.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(userUpdate)
	if err != nil {
		errorMessages := entity.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	if userUpdate.Username == nil {
		userUpdate.Username = &user.Username
	}
	if userUpdate.Email == nil {
		userUpdate.Email = &user.Email
	}
	if userUpdate.Password == nil {
		userUpdate.Password = &user.Password
	} else {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*userUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to hash password", err.Error())
			return
		}
		hashedPassword := string(passwordHash)
		userUpdate.Password = &hashedPassword
	}
	if userUpdate.Avatar == nil {
		userUpdate.Avatar = &user.Avatar
	}
	if userUpdate.Bio == nil {
		userUpdate.Bio = &user.Bio
	}
	if userUpdate.Github == nil {
		userUpdate.Github = &user.Github
	}
	if userUpdate.Linkedin == nil {
		userUpdate.Linkedin = &user.Linkedin
	}
	if userUpdate.Website == nil {
		userUpdate.Website = &user.Website
	}

	fileHandler := NewFileHandler()

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		avatarPath, err := fileHandler.UploadFile(r)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Erro no upload do avatar", err.Error())
		return
	}
	if avatarPath != "" {
		userUpdate.Avatar = &avatarPath
	}
}

	updatedUser, err := h.userService.UpdateUser(convertUUID,
		*userUpdate.Username,
		*userUpdate.Email,
		*userUpdate.Password,
		*userUpdate.Avatar,
		*userUpdate.Bio,
		*userUpdate.Github,
		*userUpdate.Linkedin,
		*userUpdate.Website,
	)

	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	res, err := json.Marshal(updatedUser)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal user", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	convertUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	_, err = h.userService.GetUser(convertUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	err = h.userService.DeleteUser(convertUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully!"))
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginUser entity.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate = validator.New()

	if err := validate.Struct(loginUser); err != nil {
		errorMessages := entity.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	user, err := h.userService.LoginUser(loginUser.Email)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			utils.NewErrorResponse(w, http.StatusNotFound, "User not found", true)
		default:
			utils.NewErrorResponse(w, http.StatusInternalServerError, "email ou senha inválido", true)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		utils.NewErrorResponse(w, http.StatusUnauthorized, "email ou senha inválido", true)
		return
	}

	println(loginUser.KeepLogin)

	token, err := configs.CreateToken(user.UserID, loginUser.KeepLogin)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	response := UserResponse{
		User:  user,
		Token: token,
	}

	res, err := json.Marshal(response)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal user", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *UserHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, "Método não permitido", "")
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

	jwtToken, err := configs.CreateToken(user.UserID, true)
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

func (h *UserHandler) VerifyTokenInSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUser(userUUID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Token valid",
		"user":    user,
	})
}

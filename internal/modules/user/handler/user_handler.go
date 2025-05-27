package user_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"

	// file_handlers "github.com/nycholasmarques/socialdev/internal/modules/file/handler"
	user_dto "github.com/nycholasmarques/socialdev/internal/modules/user/dto"
	user_errors "github.com/nycholasmarques/socialdev/internal/modules/user/errors"
	user_service "github.com/nycholasmarques/socialdev/internal/modules/user/service"
	httputil "github.com/nycholasmarques/socialdev/internal/pkg"
	pkg "github.com/nycholasmarques/socialdev/internal/pkg/auth"
	"github.com/nycholasmarques/socialdev/utils"
)

type UserHandler struct {
	userService *user_service.UserService
}

func NewUserHandler(userService *user_service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserWithUsername struct {
	Username string
}

type UserResponse struct {
	User  user_dto.UserResponse `json:"user"`
	Token string              `json:"token"`
}


func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerUser db.User

	if err := json.NewDecoder(r.Body).Decode(&registerUser); err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(registerUser); err != nil {
		errorMessages := user_errors.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	user, err := h.userService.CreateUser(registerUser.Username, registerUser.Email, registerUser.Password)
	if err != nil {
		switch err {
		case user_errors.ErrEmailAlreadyExists:
			utils.NewErrorResponse(w, http.StatusConflict, "Email already exists", err.Error())
		case user_errors.ErrUsernameAlreadyExists:
			utils.NewErrorResponse(w, http.StatusConflict, "Username already exists", err.Error())
		default:
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		}
		return
	}

	token, err := pkg.CreateToken(user.UserID.String(), true)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	userResponse := user_service.Sanitize(user)

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
		case user_errors.ErrUserNotFound:
			utils.NewErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		default:
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get user", err.Error())
		}
		return
	}

	userResponse := user_service.Sanitize(user)

	httputil.New(w).Status(200).Message("User found successfully!").Data(userResponse).Send()
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

	users, err := h.userService.GetAllUsers(username, limitStr, offsetStr)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	var usersResponse []user_dto.UserResponse

	for _, v := range users {
		user := user_service.Sanitize(v)
		usersResponse = append(usersResponse, user)
	}
	httputil.New(w).Status(200).Message("Users founds successfully!").Data(usersResponse).Send()
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middlewares .UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	convertUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var userUpdate user_dto.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(userUpdate)
	if err != nil {
		errorMessages := user_errors.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	// fileHandler := file_handlers.NewFileHandler()

	// if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
	// 	avatarPath, err := fileHandler.UploadFile(r)
	// if err != nil && !errors.Is(err, http.ErrMissingFile) {
	// 	utils.NewErrorResponse(w, http.StatusBadRequest, "Erro no upload do avatar", err.Error())
	// 	return
	// }
	// if avatarPath != "" {
	// 	userUpdate.Avatar = &avatarPath
	// }
//}

	updatedUser, err := h.userService.UpdateUser(convertUUID, userUpdate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	userResponse := user_service.Sanitize(updatedUser)
	httputil.New(w).Status(200).Message("User update successfully!").Data(userResponse).Send()
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

	err = h.userService.DeleteUser(convertUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	httputil.New(w).Status(200).Message("User Deleted successfully!").Send()
}
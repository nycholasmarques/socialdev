package auth_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nycholasmarques/socialdev/internal/db"
	services "github.com/nycholasmarques/socialdev/internal/modules/auth/service"
	user_errors "github.com/nycholasmarques/socialdev/internal/modules/user/errors"
	httputil "github.com/nycholasmarques/socialdev/internal/pkg"
	pkg "github.com/nycholasmarques/socialdev/internal/pkg/auth"
	"github.com/nycholasmarques/socialdev/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

	type UserResponse struct {
		User db.User `json:"user"`
		Token string `json:"token"`
	}

	type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
	KeepLogin bool   `json:"keep_login"`
}

func (h *AuthHandler) LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginUser LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(loginUser); err != nil {
		errorMessages := user_errors.GetValidationUserErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	user, err := h.authService.LoginUser(loginUser.Email)
	if err != nil {
		switch err {
		case user_errors.ErrUserNotFound:
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

	token, err := pkg.CreateToken(user.UserID.String(), loginUser.KeepLogin)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	response := UserResponse{
		User:  user,
		Token: token,
	}

	httputil.New(w).Status(200).Data(response).Send()
}



func (h *AuthHandler) VerifyTokenInSession(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value(middlewares.UserIDKey).(string)
	// if userID == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// userUUID, err := uuid.Parse(userID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// user, err := h.authService.GetUser(userUUID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"message": "Token valid",
	// 	"user":    user,
	// })
}

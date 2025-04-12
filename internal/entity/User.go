package entity

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type UserResponse struct {
	UserID    string `json:"user_id"`
	GoogleID  string `json:"-"`
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Password  string `json:"-" validate:"required,min=6,max=12"`
	Email     string `json:"email" validate:"required,email"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Github    string `json:"github"`
	Linkedin  string `json:"linkedin"`
	Website   string `json:"website"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserResponseSafe struct {
	UserID    string `json:"user_id"`
	GoogleID  string `json:"-"`
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Password  string `json:"-" validate:"required,min=6,max=12"`
	Email     string `json:"-" validate:"required,email"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Github    string `json:"-"`
	Linkedin  string `json:"-"`
	Website   string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=30"`
	Password *string `json:"password" validate:"omitempty,min=6,max=12"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Avatar   *string `json:"avatar"`
	Bio      *string `json:"bio"`
	Github   *string `json:"github"`
	Linkedin *string `json:"linkedin"`
	Website  *string `json:"website"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
	KeepLogin bool   `json:"keep_login"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

// definindo erros
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserExists            = errors.New("user already exists")
	ErrInternalServer        = errors.New("internal server error")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

// Função para retornar mensagens personalizadas de erro
func GetValidationUserErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()

			switch field {
			case "Username":
				errors[field] = "O nome de usuário deve ter entre 3 e 30 caracteres."
			case "Password":
				errors[field] = "A senha deve ter entre 6 e 12 caracteres."
			case "Email":
				errors[field] = "O e-mail deve ser válido."
			default:
				errors[field] = "Campo inválido."
			}
		}
	}

	return errors
}

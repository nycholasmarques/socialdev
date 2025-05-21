package entity

import (
	"errors"

	"github.com/go-playground/validator/v10"
)


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

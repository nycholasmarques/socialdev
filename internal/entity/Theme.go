package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Theme struct {
	ThemeID      string          `json:"theme_id"`
	Name         string          `json:"name"`
	Description  string          `json:"description" `
	IsPremium    bool            `json:"is_premium" `
	PreviewImage string          `json:"preview_image" `
	ThemeConfig  json.RawMessage `json:"theme_config" `
	Status       bool            `json:"status"`
	Active       bool            `json:"active"`
	ExpiresIn    time.Time       `json:"expires_in"`
}

type ListUserThemesRow struct {
	Theme 	 Theme `json:"theme"`	
	UserID   string `json:"user_id"`
	Useractive   bool   `json:"user_active_theme"`
	UserExpiresIn time.Time `json:"user_expires_in_theme"`
}

type ThemeRequest struct {
	Name         string          `json:"name" validate:"required,min=3,max=30"`
	Description  string          `json:"description" validate:"required,min=5,max=255"`
	IsPremium    bool            `json:"is_premium" `
	PreviewImage string          `json:"preview_image" validate:"required,url"`
	ThemeConfig  json.RawMessage `json:"theme_config" validate:"required"`
	Status       bool            `json:"status"`
	Active       bool            `json:"active"`
	ExpiresIn    time.Time       `json:"expires_in" validate:"required"`
}

type UserTheme struct {
	UserID   string `json:"user_id"`
	ThemeID  string `json:"theme_id"`
	Active   bool   `json:"active"`
	ExpiresIn time.Time `json:"expires_in"`
}

// definindo erros
var (
	ErrThemeNotFound = errors.New("theme not found")
)

// Função para retornar mensagens personalizadas de erro
func GetValidationThemeErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()

			switch field {
			case "Name":
				errors[field] = "O nome precisa ter no mínimo 3 caracteres e no máximo 30."
			case "Description":
				errors[field] = "A descrição precisa ter no mínimo 5 caracteres e no máximo 255."
			case "PreviewImage":
				errors[field] = "Campo requirido."
			case "ThemeConfig":
				errors[field] = "Campo requirido."
			case "Status":
				errors[field] = "Campo requirido."
			case "Active":
				errors[field] = "Campo requirido."
			case "ExpiresIn":
				errors[field] = "Campo requirido."
			default:
				errors[field] = "Campo inválido."
			}
		}
	}

	return errors
}

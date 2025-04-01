package entity

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	PostID    string        `json:"post_id" validate:"omitempty,uuid4"`
	UserID    string        `json:"user_id" validate:"omitempty,uuid4"`
	Title     string        `json:"title" validate:"omitempty,min=5,max=100"`
	Content   string        `json:"content" validate:"required,min=5"`
	Photo     string        `json:"photo" validate:"omitempty,url"`
	Video     string        `json:"video" validate:"omitempty,url"`
	LikedBy   bool        `json:"liked_by" validate:"omitempty"`
	CreatedAt string        `json:"created_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	UpdatedAt string        `json:"updated_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	User      UserPost      `json:"user" validate:"omitempty"`
	Comments  []CommentPost `json:"comments" validate:"omitempty"`
	Like      int           `json:"like" validate:"omitempty"`
	Views     int           `json:"views" validate:"omitempty"`
}

type UserPost struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type CommentPost struct {
	CommentID string   `json:"comment_id"`
	User      UserPost `json:"user"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
}

// Response representa a estrutura da resposta JSON
type Response struct {
	Posts  []Post `json:"posts"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Total  int64  `json:"total"`
}

// definindo erros
var (
	ErrPostNotFound          = errors.New("post not found")
)

// Função para retornar mensagens personalizadas de erro
func GetValidationPostErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()

			switch field {
			case "Content":
				errors[field] = "O contéudo deve ter pelo menos 5 caracteres."
			default:
				errors[field] = "Campo inválido."
			}
		}
	}

	return errors
}
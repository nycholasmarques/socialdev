package entity

import (

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Comment struct {
	CommentID uuid.UUID `json:"comment_id"`
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt string    `json:"created_at"`
	User      User      `json:"user"`
}

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar,omitempty"`
}

// var (
// 	errCommentNotFound          = errors.New("post not found")
// )

// Função para retornar mensagens personalizadas de erro
func GetValidationCommentErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()

			switch field {
			case "Content":
				errors[field] = "O contéudo deve ser preenchido."
			default:
				errors[field] = "Campo inválido."
			}
		}
	}

	return errors
}
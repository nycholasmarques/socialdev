package user_dto

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	CoverImage string    `json:"cover_image"`
	Bio        string    `json:"bio"`
	Github     string    `json:"github"`
	Linkedin   string    `json:"linkedin"`
	Website    string    `json:"website"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

	type UserResponseJSON struct {
		User UserResponse `json:"user"`
		Token string `json:"token"`
	}
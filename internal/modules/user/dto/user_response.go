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

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=30"`
	Password *string `json:"password" validate:"omitempty,min=6,max=12"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Avatar   *string `json:"avatar" validate:"omitempty"`
	Bio      *string `json:"bio" validate:"omitempty,min=5,max=255"`
	Github   *string `json:"github" validate:"omitempty,url"`
	Linkedin *string `json:"linkedin" validate:"omitempty,url"`
	Website  *string `json:"website" validate:"omitempty,url"`
}

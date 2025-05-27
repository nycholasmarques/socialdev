package user_service

import (
	"github.com/nycholasmarques/socialdev/internal/db"
	user_dto "github.com/nycholasmarques/socialdev/internal/modules/user/dto"
)

func Sanitize(user db.User) user_dto.UserResponse {
	return user_dto.UserResponse{
		UserID:     user.UserID,
		Username:   user.Username,
		Email:      user.Email,
		Avatar:     user.Avatar.String,
		CoverImage: user.CoverImage.String,
		Bio:        user.Bio.String,
		Github:     user.Github.String,
		Linkedin:   user.Linkedin.String,
		Website:    user.Website.String,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}

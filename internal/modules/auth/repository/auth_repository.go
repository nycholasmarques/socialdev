package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
)

type AuthRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) *AuthRepository {
	return &AuthRepository{queries: queries}
}

func (r *AuthRepository) LoginUser(ctx context.Context, email string) (entity.UserResponse, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return entity.UserResponse{}, entity.ErrUserNotFound
	} else if err != nil {
		log.Printf("%v", err)
	}

	var userResponse entity.UserResponse

	userResponse.UserID = user.UserID.String()
	userResponse.Username = user.Username
	userResponse.Password = user.Password
	userResponse.Email = user.Email
	userResponse.Avatar = user.Avatar.String
	userResponse.Bio = user.Bio.String
	userResponse.Github = user.Github.String
	userResponse.Linkedin = user.Linkedin.String
	userResponse.Website = user.Website.String
	userResponse.CreatedAt = user.CreatedAt.Time.String()
	userResponse.UpdatedAt = user.UpdatedAt.Time.String()

	return userResponse, err
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (entity.UserResponse, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("%v", err)
	}
	var userResponse entity.UserResponse

	userResponse.UserID = user.UserID.String()
	userResponse.Username = user.Username
	userResponse.Email = user.Email
	userResponse.Avatar = user.Avatar.String
	userResponse.Bio = user.Bio.String
	userResponse.Github = user.Github.String
	userResponse.Linkedin = user.Linkedin.String
	userResponse.Website = user.Website.String

	return userResponse, err
}

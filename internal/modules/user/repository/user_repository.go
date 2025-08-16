package user_repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (entity.UserResponse, error) {
	user, err := r.queries.CreateUser(ctx, params)
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

func (r *UserRepository) GetUser(ctx context.Context, user_id uuid.UUID) (db.User, error) {
	user, err := r.queries.GetUser(ctx, user_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}

func (r *UserRepository) GetAllUsers(ctx context.Context, params db.GetAllUsersParams) ([]db.User, error) {
	allUsers, err := r.queries.GetAllUsers(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return allUsers, err
}

func (r *UserRepository) GetUserWithUsername(ctx context.Context, username string) ([]db.GetUserWithUsernameRow, error) {
	usernamePattern := username + "%"
	getUserWithUsername, err := r.queries.GetUserWithUsername(ctx, usernamePattern)

	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if len(getUserWithUsername) == 0 {
		return nil, sql.ErrNoRows
	}

	return getUserWithUsername, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.UserResponse, error) {
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

func (r *UserRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, user_id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, user_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}

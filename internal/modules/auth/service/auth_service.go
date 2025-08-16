package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo *repository.AuthRepository
	queries  *db.Queries
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{AuthRepo: authRepo}
}

func (s *AuthService) LoginUser(email string) (entity.UserResponse, error) {
	ctx := context.Background()
	user, err := s.AuthRepo.LoginUser(ctx, email)
	if err == sql.ErrNoRows {
		return entity.UserResponse{}, entity.ErrUserNotFound
	} else if err != nil {
		return entity.UserResponse{}, entity.ErrInternalServer
	}
	return user, nil
}

func (s *AuthService) CreateUserWithGoogle(username, email, password, googleID, avatar string) (entity.UserResponse, error) {
	ctx := context.Background()

	existingUser, err := s.AuthRepo.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting user by email: %v", err)
		return entity.UserResponse{}, fmt.Errorf("erro ao verificar usuário existente: %v", err)
	}

	if existingUser.UserID != "" {
		if existingUser.GoogleID == googleID {
			return entity.UserResponse{
				UserID:   existingUser.UserID,
				Username: existingUser.Username,
				Email:    existingUser.Email,
			}, nil
		}
	}

	existingUserByUsername, err := s.queries.GetUserWithUsername(ctx, username)
	if err != nil && err != sql.ErrNoRows {
		return entity.UserResponse{}, fmt.Errorf("erro ao verificar nome de usuário: %v", err)
	}

	if len(existingUserByUsername) > 0 {
		username = fmt.Sprintf("%s_%s", username, uuid.New().String()[:8])
	}

	var passwordHash string
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error encrypting password: %v", err)
			return entity.UserResponse{}, fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		passwordHash = string(hash)
	}

	userParams := db.CreateUserParams{
		UserID:   uuid.New(),
		Username: username,
		Password: passwordHash,
		Email:    email,
		GoogleID: sql.NullString{String: googleID, Valid: googleID != ""},
		Avatar:   sql.NullString{String: avatar, Valid: avatar != ""},
	}

	createdUser, err := s.queries.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return entity.UserResponse{}, fmt.Errorf("{erro ao criar usuário: %v}", err)
	}

	return entity.UserResponse{
		// UserID:   createdUser.UserID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		// Avatar:   createdUser.Avatar,
	}, nil
}

package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(username, email, password string) (entity.UserResponse, error) {
	ctx := context.Background()

	// Gerar hash da senha
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error encrypting password: %v", err)
		return entity.UserResponse{}, err
	}

	// Criar parâmetros para o repositório
	userParams := db.CreateUserParams{
		UserID:   uuid.New(),
		Username: username,
		Password: string(passwordHash),
		Email:    email,
	}
	// Chamar o repositório para criar o usuário
	return s.userRepo.CreateUser(ctx, userParams)
}

func (s *UserService) CreateUserWithGoogle(username, email, password, googleID, avatar string) (entity.UserResponse, error) {
	ctx := context.Background()

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
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

	existingUserByUsername, err := s.userRepo.GetUserWithUsername(ctx, username)
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

	createdUser, err := s.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return entity.UserResponse{}, fmt.Errorf("{erro ao criar usuário: %v}", err)
	}

	return entity.UserResponse{
		UserID:   createdUser.UserID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Avatar:   createdUser.Avatar,
	}, nil
}

func (s *UserService) GetUser(user_id uuid.UUID) (db.User, error) {
	ctx := context.Background()
	return s.userRepo.GetUser(ctx, user_id)
}

func (s *UserService) GetAllUsers() ([]db.User, error) {
	ctx := context.Background()
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUserWithUsername(username string) ([]db.GetUserWithUsernameRow, error) {
	ctx := context.Background()
	return s.userRepo.GetUserWithUsername(ctx, username)
}

func (s *UserService) UpdateUser(user_id uuid.UUID, username, email, password, avatar, bio, github, linkedin, website string) (db.User, error) {
	ctx := context.Background()

	return s.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		UserID:   user_id,
		Username: username,
		Email:    email,
		Password: password,
		Avatar:   utils.StringToSqlNull(avatar),
		Bio:      utils.StringToSqlNull(bio),
		Github:   utils.StringToSqlNull(github),
		Linkedin: utils.StringToSqlNull(linkedin),
		Website:  utils.StringToSqlNull(website),
	})
}

func (s *UserService) DeleteUser(user_id uuid.UUID) error {
	ctx := context.Background()
	return s.userRepo.DeleteUser(ctx, user_id)
}

func (s *UserService) LoginUser(email string) (entity.UserResponse, error) {
	ctx := context.Background()
	return s.userRepo.LoginUser(ctx, email)
}

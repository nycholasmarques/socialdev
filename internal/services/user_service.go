package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
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

func (s *UserService) CreateUser(username, email, password string) (db.User, error) {
	ctx := context.Background()

	// Gerar hash da senha
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error encrypting password: %v", err)
		return db.User{}, err
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

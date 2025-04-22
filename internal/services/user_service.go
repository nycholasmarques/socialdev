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

	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return entity.UserResponse{}, entity.ErrEmailAlreadyExists
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking email existence: %v", err)
		return entity.UserResponse{}, entity.ErrInternalServer
	}

	_, err = s.userRepo.GetUserWithUsername(ctx, username)
	if err == nil {
		return entity.UserResponse{}, entity.ErrUsernameAlreadyExists
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking username existence: %v", err)
		return entity.UserResponse{}, entity.ErrInternalServer
	}

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

	userResponse, err := s.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return entity.UserResponse{}, entity.ErrInternalServer
	}

	return userResponse, nil
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

func (s *UserService) GetUser(user_id uuid.UUID) (entity.UserResponse, error) {
	ctx := context.Background()
	user, err := s.userRepo.GetUser(ctx, user_id)
	if err == sql.ErrNoRows {
		return entity.UserResponse{}, entity.ErrUserNotFound
	} else if err != nil {
		log.Printf("Error getting user: %v", err)
		return entity.UserResponse{}, entity.ErrInternalServer
	}

	entityUser := entity.UserResponse{
		UserID:    user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Avatar:    user.Avatar.String,
		Bio:       user.Bio.String,
		Github:    user.Github.String,
		Linkedin:  user.Linkedin.String,
		Website:   user.Website.String,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	}

	return entityUser, nil
}

func (s *UserService) GetAllUsers(username string, limit, offset int32) ([]entity.UserResponseSafe, error) {
	ctx := context.Background()

	params := db.GetAllUsersParams{
		Username: sql.NullString{String: username, Valid: username != ""},
		Limit:    sql.NullInt32{Int32: limit, Valid: true},
		Offset:   sql.NullInt32{Int32: offset, Valid: true},
	}

	users, err := s.userRepo.GetAllUsers(ctx, params)

	var usersResponse []entity.UserResponseSafe

	for _, user := range users {
		var userResponse entity.UserResponseSafe

		userResponse.UserID = user.UserID.String()
		userResponse.Username = user.Username
		userResponse.Email = user.Email
		userResponse.Avatar = user.Avatar.String
		userResponse.Bio = user.Bio.String
		userResponse.Github = user.Github.String
		userResponse.Linkedin = user.Linkedin.String
		userResponse.Website = user.Website.String
		userResponse.CreatedAt = user.CreatedAt.Time.String()
		userResponse.UpdatedAt = user.UpdatedAt.Time.String()

		usersResponse = append(usersResponse, userResponse)
	}

	return usersResponse, err
}

func (s *UserService) GetUserWithUsername(username string) ([]db.GetUserWithUsernameRow, error) {
	ctx := context.Background()
	return s.userRepo.GetUserWithUsername(ctx, username)
}

func (s *UserService) UpdateUser(user_id uuid.UUID, username, email, password, avatar, bio, github, linkedin, website string) (entity.UserResponse, error) {
	ctx := context.Background()

	user, err := s.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		UserID:   user_id,
		Username: username,
		Email:    email,
		Password: password,
		Avatar:   sql.NullString{String: avatar, Valid: avatar != ""},
		Bio:      sql.NullString{String: bio, Valid: bio != ""},
		Github:   sql.NullString{String: github, Valid: github != ""},
		Linkedin: sql.NullString{String: linkedin, Valid: linkedin != ""},
		Website:  sql.NullString{String: website, Valid: website != ""},
	})
	if err != nil {

		return entity.UserResponse{}, entity.ErrInternalServer
	}

	return entity.UserResponse{
		UserID:    user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar.String,
		Bio:       user.Bio.String,
		Github:    user.Github.String,
		Linkedin:  user.Linkedin.String,
		Website:   user.Website.String,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	}, nil
}

func (s *UserService) DeleteUser(user_id uuid.UUID) error {
	ctx := context.Background()
	return s.userRepo.DeleteUser(ctx, user_id)
}

func (s *UserService) LoginUser(email string) (entity.UserResponse, error) {
	ctx := context.Background()
	user, err := s.userRepo.LoginUser(ctx, email)
	if err == sql.ErrNoRows {
		return entity.UserResponse{}, entity.ErrUserNotFound
	} else if err != nil {
		return entity.UserResponse{}, entity.ErrInternalServer
	}
	return user, nil
}

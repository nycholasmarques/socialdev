package user_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	user_dto "github.com/nycholasmarques/socialdev/internal/modules/user/dto"
	user_errors "github.com/nycholasmarques/socialdev/internal/modules/user/errors"
	user_repository "github.com/nycholasmarques/socialdev/internal/modules/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *user_repository.UserRepository
}

func NewUserService(userRepo *user_repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(username, email, password string) (db.User, error) {
	ctx := context.Background()

	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return db.User{}, user_errors.ErrEmailAlreadyExists
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking email existence: %v", err)
		return db.User{}, user_errors.ErrInternalServer
	}
	
	_, err = s.userRepo.GetUserWithUsername(ctx, username)
	if err == nil {
		return db.User{}, user_errors.ErrUsernameAlreadyExists
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking username existence: %v", err)
		return db.User{}, user_errors.ErrInternalServer
	}

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

	userResponse, err := s.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return db.User{}, user_errors.ErrInternalServer
	}

	return userResponse, nil
}

func (s *UserService) CreateUserWithGoogle(username, email, password, googleID, avatar string) (db.User, error) {
	ctx := context.Background()

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting user by email: %v", err)
		return db.User{}, fmt.Errorf("erro ao verificar usuário existente: %v", err)
	}


	if existingUser.UserID.String() != "" {
		if existingUser.GoogleID.String == googleID {
			return db.User{
				UserID:   existingUser.UserID,
				Username: existingUser.Username,
				Email:    existingUser.Email,
			}, nil
		}
	}

	existingUserByUsername, err := s.userRepo.GetUserWithUsername(ctx, username)
	if err != nil && err != sql.ErrNoRows {
		return db.User{}, fmt.Errorf("erro ao verificar nome de usuário: %v", err)
	}

	if len(existingUserByUsername) > 0 {
		username = fmt.Sprintf("%s_%s", username, uuid.New().String()[:8])
	}

	var passwordHash string
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error encrypting password: %v", err)
			return db.User{}, fmt.Errorf("erro ao gerar hash da senha: %v", err)
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
		return db.User{}, fmt.Errorf("{erro ao criar usuário: %v}", err)
	}

	return db.User{
		UserID:   createdUser.UserID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Avatar:   createdUser.Avatar,
	}, nil
}

func (s *UserService) GetUser(user_id uuid.UUID) (db.User, error) {
	ctx := context.Background()
	user, err := s.userRepo.GetUser(ctx, user_id)
	if err == sql.ErrNoRows {
		return db.User{}, user_errors.ErrUserNotFound
	} else if err != nil {
		log.Printf("Error getting user: %v", err)
		return db.User{}, user_errors.ErrInternalServer
	}

	return user, nil
}

func (s *UserService) GetAllUsers(username, limitStr, offsetStr string) ([]db.User, error) {
	ctx := context.Background()

	limit := int32(10)
	offset := int32(0)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	params := db.GetAllUsersParams{
		Username: sql.NullString{String: username, Valid: username != ""},
		Limit:    sql.NullInt32{Int32: limit, Valid: true},
		Offset:   sql.NullInt32{Int32: offset, Valid: true},
	}

	users, err := s.userRepo.GetAllUsers(ctx, params)

	return users, err
}

func (s *UserService) GetUserWithUsername(username string) ([]db.GetUserWithUsernameRow, error) {
	ctx := context.Background()
	return s.userRepo.GetUserWithUsername(ctx, username)
}

func (s *UserService) UpdateUser(user_id uuid.UUID, userUpdate user_dto.UpdateUserRequest) (db.User, error) {
	ctx := context.Background()
	
	userExist, err := s.GetUser(user_id)
	if err == sql.ErrNoRows {
		return db.User{}, user_errors.ErrUserNotFound
	} else if err != nil {
		return db.User{}, user_errors.ErrInternalServer
	}

	if userUpdate.Username == nil {
		userUpdate.Username = &userExist.Username
	}
	if userUpdate.Email == nil {
		userUpdate.Email = &userExist.Email
	}
	if userUpdate.Password == nil {
		userUpdate.Password = &userExist.Password
	} else {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*userUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			return db.User{}, user_errors.ErrInternalServer
		}
		hashedPassword := string(passwordHash)
		userUpdate.Password = &hashedPassword
	}
	if userUpdate.Avatar == nil {
		userUpdate.Avatar = &userExist.Avatar.String
	}
	if userUpdate.Bio == nil {
		userUpdate.Bio = &userExist.Bio.String
	}
	if userUpdate.Github == nil {
		userUpdate.Github = &userExist.Github.String
	}
	if userUpdate.Linkedin == nil {
		userUpdate.Linkedin = &userExist.Linkedin.String
	}
	if userUpdate.Website == nil {
		userUpdate.Website = &userExist.Website.String
	}

	user, err := s.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		UserID:   user_id,
		Username: *userUpdate.Username,
		Email:    *userUpdate.Email,
		Password: *userUpdate.Password,
		Avatar:   sql.NullString{String: *userUpdate.Avatar, Valid: *userUpdate.Avatar != ""},
		Bio:      sql.NullString{String: *userUpdate.Bio, Valid: *userUpdate.Bio != ""},
		Github:   sql.NullString{String: *userUpdate.Github, Valid: *userUpdate.Github != ""},
		Linkedin: sql.NullString{String: *userUpdate.Linkedin, Valid: *userUpdate.Linkedin != ""},
		Website:  sql.NullString{String: *userUpdate.Website, Valid: *userUpdate.Website != ""},
	})
	if err != nil {
		return db.User{}, user_errors.ErrInternalServer
	}

	return user, nil
}

func (s *UserService) DeleteUser(user_id uuid.UUID) error {
	ctx := context.Background()

	_, err := s.GetUser(user_id)
	if err != nil {
		return user_errors.ErrInternalServer
	}

	return s.userRepo.DeleteUser(ctx, user_id)
}

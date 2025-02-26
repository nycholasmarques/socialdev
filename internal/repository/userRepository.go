package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func GetUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}

func (r *UserRepository) GetUser(ctx context.Context, user_id uuid.UUID) (db.GetUserRow, error) {
	user, err := r.queries.GetUser(ctx, user_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}
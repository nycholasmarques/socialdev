package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/nycholasmarques/socialdev/internal/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Println(user)
	return user, err
}
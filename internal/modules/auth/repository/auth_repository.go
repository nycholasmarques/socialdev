package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/nycholasmarques/socialdev/internal/db"
	user_errors "github.com/nycholasmarques/socialdev/internal/modules/user/errors"
)

type AuthRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) *AuthRepository {
	return &AuthRepository{queries: queries}
}

func (r *AuthRepository) LoginUser(ctx context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return db.User{}, user_errors.ErrUserNotFound
	} else if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}


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

func (r *AuthRepository) LoginUser(ctx context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return db.User{}, entity.ErrUserNotFound
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


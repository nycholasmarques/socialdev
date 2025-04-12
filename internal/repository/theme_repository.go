package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
)

type ThemeRepository struct {
	queries *db.Queries
}

func NewThemeRepository(queries *db.Queries) *ThemeRepository {
	return &ThemeRepository{queries: queries}
}

func GetThemeRepository(queries *db.Queries) *ThemeRepository {
	return &ThemeRepository{queries: queries}
}

func (r *ThemeRepository) CreateTheme(ctx context.Context, params db.CreateThemeParams) (db.Theme, error) {
	theme, err := r.queries.CreateTheme(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return theme, err
}

func (r *ThemeRepository) ListAllThemes(ctx context.Context) ([]db.Theme, error) {
	allThemes, err := r.queries.ListAllThemes(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
	return allThemes, err
}

func (r *ThemeRepository) GetThemeByID(ctx context.Context, theme_id uuid.UUID) (db.Theme, error) {
	theme, err := r.queries.GetThemeByID(ctx, theme_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return theme, err
}

func (r *ThemeRepository) UpdateTheme(ctx context.Context, params db.UpdateThemeParams) (db.Theme, error) {
	theme, err := r.queries.UpdateTheme(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return theme, err
}

func (r *ThemeRepository) DeleteTheme(ctx context.Context, theme_id uuid.UUID) error {
	err := r.queries.DeleteTheme(ctx, theme_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}

func (r *ThemeRepository) AddUserTheme(ctx context.Context, params db.AddUserThemeParams) (db.UserTheme, error) {
	userTheme, err := r.queries.AddUserTheme(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return userTheme, err
}

func (r *ThemeRepository) ListUserThemes(ctx context.Context, user_id uuid.UUID) ([]db.Theme, error) {
	listUserTheme, err := r.queries.ListUserThemes(ctx, user_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return listUserTheme, err
}

func (r *ThemeRepository) UserHasTheme(ctx context.Context, params db.UserHasThemeParams) (bool, error) {
	userHasTheme, err := r.queries.UserHasTheme(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return userHasTheme, err
}

func (r *ThemeRepository) DeleteUserTheme(ctx context.Context, params db.DeleteUserThemeParams) error {
	err := r.queries.DeleteUserTheme(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}
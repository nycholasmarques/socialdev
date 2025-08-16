package theme_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	theme_repository "github.com/nycholasmarques/socialdev/internal/modules/theme/repository"
	"github.com/sqlc-dev/pqtype"
)

type ThemeService struct {
	themeRepo *theme_repository.ThemeRepository
}

func NewThemeService(themeRepo *theme_repository.ThemeRepository) *ThemeService {
	return &ThemeService{themeRepo: themeRepo}
}

func (r *ThemeService) CreateTheme(name, description, previewImage string, isPremium, status, active bool, expiresIn time.Time, themeConfig json.RawMessage) (db.Theme, error) {
	ctx := context.Background()

	params := db.CreateThemeParams{
		ThemeID:      uuid.New(),
		Name:         sql.NullString{String: name, Valid: name != ""},
		Description:  sql.NullString{String: description, Valid: description != ""},
		IsPremium:    sql.NullBool{Bool: isPremium, Valid: true},
		PreviewImage: sql.NullString{String: previewImage, Valid: previewImage != ""},
		ThemeConfig:  pqtype.NullRawMessage{RawMessage: themeConfig, Valid: themeConfig != nil},
		Status:       sql.NullBool{Bool: status, Valid: true},
		Active:       sql.NullBool{Bool: active, Valid: true},
		ExpiresIn:    sql.NullTime{Time: expiresIn, Valid: expiresIn != time.Time{}},
	}

	return r.themeRepo.CreateTheme(ctx, params)
}

func (r *ThemeService) ListAllThemes() ([]db.Theme, error) {
	ctx := context.Background()

	return r.themeRepo.ListAllThemes(ctx)
}

func (r *ThemeService) GetThemeByID(theme_id uuid.UUID) (db.Theme, error) {
	ctx := context.Background()

	return r.themeRepo.GetThemeByID(ctx, theme_id)
}

func (r *ThemeService) UpdateTheme(theme_id uuid.UUID, name, description, previewImage string, isPremium, status, active bool, expiresIn time.Time, themeConfig json.RawMessage) (db.Theme, error) {
	ctx := context.Background()

	params := db.UpdateThemeParams{
		ThemeID:      theme_id,
		Name:         sql.NullString{String: name, Valid: name != ""},
		Description:  sql.NullString{String: description, Valid: description != ""},
		IsPremium:    sql.NullBool{Bool: isPremium, Valid: true},
		PreviewImage: sql.NullString{String: previewImage, Valid: previewImage != ""},
		ThemeConfig:  pqtype.NullRawMessage{RawMessage: themeConfig, Valid: themeConfig != nil},
		Status:       sql.NullBool{Bool: status, Valid: true},
		Active:       sql.NullBool{Bool: active, Valid: true},
		ExpiresIn:    sql.NullTime{Time: expiresIn, Valid: expiresIn != time.Time{}},
	}

	return r.themeRepo.UpdateTheme(ctx, params)
}

func (r *ThemeService) DeleteTheme(theme_id uuid.UUID) error {
	ctx := context.Background()

	return r.themeRepo.DeleteTheme(ctx, theme_id)
}

func (r *ThemeService) AddUserTheme(user_id uuid.UUID, theme_id uuid.UUID, active bool) (db.UserTheme, error) {
	ctx := context.Background()

	params := db.AddUserThemeParams{UserID: user_id, ThemeID: theme_id, Active: sql.NullBool{Bool: active, Valid: true}}
	return r.themeRepo.AddUserTheme(ctx, params)
}

func (r *ThemeService) UpdateUserTheme(user_id uuid.UUID, theme_id uuid.UUID, active bool) (db.UserTheme, error) {
	ctx := context.Background()

	userThemeList, err := r.themeRepo.ListUserThemes(ctx, user_id)
	if err != nil {
		return db.UserTheme{}, err
	}

	for _, userTheme := range userThemeList {
		if userTheme.UserActive.Bool {
			params := db.UpdateUserThemeParams{UserID: user_id, ThemeID: userTheme.ThemeID, Active: sql.NullBool{Bool: false, Valid: true}}
			r.themeRepo.UpdateUserTheme(ctx, params)
		}
	}

	params := db.UpdateUserThemeParams{UserID: user_id, ThemeID: theme_id, Active: sql.NullBool{Bool: active, Valid: true}}
	return r.themeRepo.UpdateUserTheme(ctx, params)
}

func (r *ThemeService) ListUserThemes(user_id uuid.UUID) ([]db.ListUserThemesRow, error) {
	ctx := context.Background()

	return r.themeRepo.ListUserThemes(ctx, user_id)
}

func (r *ThemeService) UserHasTheme(user_id uuid.UUID, theme_id uuid.UUID) (bool, error) {
	ctx := context.Background()

	params := db.UserHasThemeParams{UserID: user_id, ThemeID: theme_id}
	return r.themeRepo.UserHasTheme(ctx, params)
}

func (r *ThemeService) UserHasActiveTheme(user_id uuid.UUID) (*db.UserHasActiveThemeRow, error) {
	ctx := context.Background()

	return r.themeRepo.UserHasActiveTheme(ctx, user_id)
}

func (r *ThemeService) DeleteUserTheme(user_id uuid.UUID, theme_id uuid.UUID) error {
	ctx := context.Background()

	userHasTheme, err := r.themeRepo.UserHasTheme(ctx, db.UserHasThemeParams{
		UserID:  user_id,
		ThemeID: theme_id,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	if !userHasTheme {
		return errors.New("user does not have this theme")
	}

	params := db.DeleteUserThemeParams{UserID: user_id, ThemeID: theme_id}
	return r.themeRepo.DeleteUserTheme(ctx, params)
}

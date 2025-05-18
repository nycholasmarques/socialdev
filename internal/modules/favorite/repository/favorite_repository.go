package favorite_repository

import (
	"context"
	"log"

	"github.com/nycholasmarques/socialdev/internal/db"
)

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, params db.AddFavoriteParams) (db.Favorite, error)
	RemoveFavorite(ctx context.Context, params db.RemoveFavoriteParams) error
	GetUserFavorites(ctx context.Context, params db.GetUserFavoritesParams) ([]db.GetUserFavoritesRow, error)
	CheckFavoriteExists(ctx context.Context, params db.CheckFavoriteExistsParams) (bool, error)
}

type favoriteRepository struct {
	queries *db.Queries
}

func NewFavoriteRepository(queries *db.Queries) FavoriteRepository {
	return &favoriteRepository{queries: queries}
}

func (r *favoriteRepository) AddFavorite(ctx context.Context, params db.AddFavoriteParams) (db.Favorite, error) {
	favorite, err := r.queries.AddFavorite(ctx, params)
	if err != nil {
		log.Printf("Erro ao adicionar favorito: %v", err)
		return db.Favorite{}, err
	}
	return favorite, nil
}

func (r *favoriteRepository) RemoveFavorite(ctx context.Context, params db.RemoveFavoriteParams) error {
	err := r.queries.RemoveFavorite(ctx, params)
	if err != nil {
		log.Printf("Erro ao remover favorito: %v", err)
	}
	return err
}

func (r *favoriteRepository) GetUserFavorites(ctx context.Context, params db.GetUserFavoritesParams) ([]db.GetUserFavoritesRow, error) {
	favorites, err := r.queries.GetUserFavorites(ctx, params)
	if err != nil {
		log.Printf("Erro ao listar favoritos: %v", err)
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) CheckFavoriteExists(ctx context.Context, params db.CheckFavoriteExistsParams) (bool, error) {
	exists, err := r.queries.CheckFavoriteExists(ctx, params)
	if err != nil {
		log.Printf("Erro ao verificar favorito existente: %v", err)
		return false, err
	}
	return exists, nil
}

package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/repository"
)

type FavoriteService interface {
	AddFavorite(userID, postID string) error
	RemoveFavorite(userID, postID string) error
	GetUserFavorites(userID string, limit, offset int32) ([]entity.FavoriteResponse, error)
}

type favoriteService struct {
	favoriteRepo repository.FavoriteRepository
}

func NewFavoriteService(favoriteRepo repository.FavoriteRepository) FavoriteService {
	return &favoriteService{favoriteRepo: favoriteRepo}
}

func (s *favoriteService) AddFavorite(userID, postID string) error {
	ctx := context.Background()

	// Parse dos IDs para UUID
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		return err
	}

	favoriteExists, err := s.favoriteRepo.CheckFavoriteExists(ctx, db.CheckFavoriteExistsParams{
		UserID: userIDUUID,
		PostID: postIDUUID,
	})
	if err != nil {
		return err
	}
	if favoriteExists {
		return errors.New("favorite already exists")
	}

	params := db.AddFavoriteParams{
		UserID: userIDUUID,
		PostID: postIDUUID,
	}
	_, err = s.favoriteRepo.AddFavorite(ctx, params)
	return err
}

func (s *favoriteService) RemoveFavorite(userID, postID string) error {
	ctx := context.Background()

	post_id_UUID, err := uuid.Parse(postID)
	if err != nil {
		return err
	}

	user_id_UUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	params := db.RemoveFavoriteParams{
		UserID: user_id_UUID,
		PostID: post_id_UUID,
	}
	return s.favoriteRepo.RemoveFavorite(ctx, params)
}

func (s *favoriteService) GetUserFavorites(userID string, limit, offset int32) ([]entity.FavoriteResponse, error) {
	ctx := context.Background()

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
			return nil, err
	}

	params := db.GetUserFavoritesParams{
			UserID: userIDUUID,
			Limit:  limit,
			Offset: offset,
	}

	favorites, err := s.favoriteRepo.GetUserFavorites(ctx, params)
	if err != nil {
			return nil, err
	}

	var favoritesResponse []entity.FavoriteResponse
	for _, fav := range favorites {
			favResponse := entity.FavoriteResponse{
					PostID:            fav.PostID.String(),
					UserID:            fav.UserID.String(),
					Title:             fav.Title.String,
					Content:           fav.Content,
					Photo:             fav.Photo.String,
					Video:             fav.Video.String,
					PostCreatedAt:     fav.PostCreatedAt.Time.UTC().Format(time.RFC3339),
					UpdatedAt:         fav.UpdatedAt.Time.UTC().Format(time.RFC3339),
					FavoriteCreatedAt: fav.FavoriteCreatedAt.Time.UTC().Format(time.RFC3339),
					UserName:          fav.UserName,
					UserAvatar:        fav.UserAvatar.String,
					LikeCount:         fav.LikeCount,
					ViewsCount:        fav.ViewsCount,
					CommentCount:      fav.CommentCount,
					LikedByUser:       fav.LikedByUser,
			}
			favoritesResponse = append(favoritesResponse, favResponse)
	}

	return favoritesResponse, nil
}

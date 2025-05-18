package like_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	like_repository "github.com/nycholasmarques/socialdev/internal/modules/like/repository"
)

type LikeService struct {
	likeRepo *like_repository.LikeRepository
}

func NewLikeService(likeRepo *like_repository.LikeRepository) *LikeService {
	return &LikeService{likeRepo: likeRepo}
}

func (r *LikeService) CreateLike(user_id uuid.UUID, post_id uuid.UUID) (db.Like, error) {
	ctx := context.Background()

	params := db.CreateLikeParams{UserID: user_id, PostID: post_id}
	return r.likeRepo.CreateLike(ctx, params)
}

func (r *LikeService) GetLike(user_id uuid.UUID, post_id uuid.UUID) (db.Like, error) {
	ctx := context.Background()

	params := db.GetLikeParams{UserID: user_id, PostID: post_id}
	return r.likeRepo.GetLike(ctx, params)
}

func (r *LikeService) DeleteLike(user_id uuid.UUID, post_id uuid.UUID) error {
	ctx := context.Background()

	params := db.DeleteLikeParams{UserID: user_id, PostID: post_id}
	return r.likeRepo.DeleteLike(ctx, params)
}

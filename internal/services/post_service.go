package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostSerivice(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (p *PostService) CreatePost(user_id, title, content, photo, video string) (entity.Post, error) {
	ctx := context.Background()

	user_id_UUID, err := uuid.Parse(user_id)
	if err != nil {
		return entity.Post{}, err
	}

	post_id_UUID := uuid.New()

	params := db.CreatePostParams{PostID: post_id_UUID, UserID: user_id_UUID, Title: sql.NullString{String: title, Valid: title != ""}, Content: content, Photo: sql.NullString{String: title, Valid: title != ""}, Video: sql.NullString{String: title, Valid: title != ""}}

	return p.postRepo.CreatePost(ctx, params)
}
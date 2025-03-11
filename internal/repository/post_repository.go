package repository

import (
	"context"
	"log"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
)

type PostRepository struct {
	queries *db.Queries
}

func NewPostRepository(queries *db.Queries) *PostRepository {
	return &PostRepository{queries: queries}
}

func GetPostRepository(queries *db.Queries) *PostRepository {
	return &PostRepository{queries: queries}
}

func (r *PostRepository) CreatePost(ctx context.Context, params db.CreatePostParams) (entity.Post, error) {
	post, err := r.queries.CreatePost(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}

	var postResponse entity.Post

	postResponse.PostID = post.PostID.String()
	postResponse.UserID = post.UserID.String()
	postResponse.Title = post.Title.String
	postResponse.Content = post.Content
	postResponse.Photo = post.Photo.String
	postResponse.Video = post.Video.String
	postResponse.CreatedAt = post.CreatedAt.Time.String()
	postResponse.UpdatedAt = post.UpdatedAt.Time.String()

	return postResponse, err
}
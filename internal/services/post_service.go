package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
		fmt.Println(user_id_UUID)
		return entity.Post{}, err
	}

	post_id_UUID := uuid.New()

	params := db.CreatePostParams{PostID: post_id_UUID, UserID: user_id_UUID, Title: sql.NullString{String: title, Valid: title != ""}, Content: content, Photo: sql.NullString{String: photo, Valid: photo != ""}, Video: sql.NullString{String: video, Valid: video != ""}}

	return p.postRepo.CreatePost(ctx, params)
}

func (p *PostService) GetPost(post_id string) (entity.Post, error) {
	ctx := context.Background()

	post_id_UUID, err := uuid.Parse(post_id)
	if err != nil {
		return entity.Post{}, err
	}

	return p.postRepo.GetPost(ctx, post_id_UUID)
}


func (p *PostService) GetAllPosts(title, content *string, orderBy string, limit, offset int32, user_id string) ([]entity.Post, int64, error) {
	ctx := context.Background()
	return p.postRepo.GetAllPosts(ctx, title, content, orderBy, limit, offset, user_id)
}

func (p *PostService) GetAllUserPosts(user_id string) ([]entity.Post, error) {
	ctx := context.Background()

	user_id_UUID, err := uuid.Parse(user_id)
	if err != nil {
		return nil, err
	}

	return p.postRepo.GetAllUserPosts(ctx, user_id_UUID)
}

func (p *PostService) GetUserPost(post_id, user_id string) (entity.Post, error) {
	ctx := context.Background()

	post_id_UUID, err := uuid.Parse(post_id)
	if err != nil {
		return entity.Post{}, err
	}

	user_id_UUID, err := uuid.Parse(user_id)
	if err != nil {
		return entity.Post{}, err
	}

	return p.postRepo.GetUserPost(ctx, db.GetUserPostParams{PostID: post_id_UUID, UserID: user_id_UUID})
}

func (p *PostService) UpdatePost(post_id, user_id, title, content, photo, video string) (entity.Post, error) {
	ctx := context.Background()

	post_id_UUID, err := uuid.Parse(post_id)
	if err != nil {
		fmt.Println("ok")
		return entity.Post{}, err
	}

	user_id_UUID, err := uuid.Parse(user_id)
	if err != nil {
		return entity.Post{}, err
	}

	_, err = p.postRepo.GetUserPost(ctx, db.GetUserPostParams{PostID: post_id_UUID, UserID: user_id_UUID})
	if err != nil {
		log.Printf("Usuário não autorizado: %v", err)
		return entity.Post{}, err
	}

	params := db.UpdatePostParams{PostID: post_id_UUID, UserID: user_id_UUID, Title: sql.NullString{String: title, Valid: title != ""}, Content: content, Photo: sql.NullString{String: photo, Valid: title != ""}, Video: sql.NullString{String: video, Valid: title != ""}}

	return p.postRepo.UpdatePost(ctx, params)
}

func (p *PostService) DeletePost(post_id, user_id string) error {
	ctx := context.Background()

	post_id_UUID, err := uuid.Parse(post_id)
	if err != nil {
		return err
	}

	user_id_UUID, err := uuid.Parse(user_id)
	if err != nil {
		return err
	}

	_, err = p.postRepo.GetUserPost(ctx, db.GetUserPostParams{PostID: post_id_UUID, UserID: user_id_UUID})
	if err != nil {
		log.Printf("Usuário não autorizado: %v", err)
		return err
	}

	return p.postRepo.DeletePost(ctx, post_id_UUID)
}

package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(userID uuid.UUID, postID uuid.UUID, content string) (db.Comment, error) {
	ctx := context.Background()

	params := db.CreateCommentParams{
		CommentID: uuid.New(),
		UserID:    userID,
		PostID:    postID,
		Content:   content,
	}

	return s.commentRepo.CreateComment(ctx, params)
}

func (s *CommentService) GetComment(commentID uuid.UUID) (db.Comment, error) {
	ctx := context.Background()
	return s.commentRepo.GetComment(ctx, commentID)
}

func (s *CommentService) UpdateComment(userID uuid.UUID, commentID uuid.UUID, content string) (db.Comment, error) {
	ctx := context.Background()
	return s.commentRepo.UpdateComment(ctx, userID, commentID, content)
}

func (s *CommentService) DeleteComment(userID, commentID uuid.UUID) error {
	ctx := context.Background()

	params := db.DeleteCommentParams{
		CommentID: commentID,
		UserID:    userID,
	}

	return s.commentRepo.DeleteComment(ctx, params)
}

func (s *CommentService) GetAllComments(post_id uuid.UUID) ([]db.Comment, error) {
	ctx := context.Background()
	return s.commentRepo.GetAllComments(ctx, post_id)
}

package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
)

type CommentRepository struct {
	queries *db.Queries
}

func NewCommentRepository(queries *db.Queries) *CommentRepository {
	return &CommentRepository{queries: queries}
}

func (r *CommentRepository) CreateComment(ctx context.Context, params db.CreateCommentParams) (db.Comment, error) {
	comment, err := r.queries.CreateComment(ctx, params)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		return db.Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepository) GetComment(ctx context.Context, commentID uuid.UUID) (db.Comment, error) {
	comment, err := r.queries.GetComment(ctx, commentID)
	if err != nil {
		log.Printf("Error retrieving comment: %v", err)
		return db.Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, content string) (db.Comment, error) {
	params := db.UpdateCommentParams{
		CommentID: commentID,
		Content:   content,
		UserID:    userID,
	}
	comment, err := r.queries.UpdateComment(ctx, params)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		return db.Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, params db.DeleteCommentParams) error {
	err := r.queries.DeleteComment(ctx, params)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		return err
	}
	return nil
}

func (r *CommentRepository) GetAllComments(ctx context.Context, post_id uuid.UUID) ([]db.Comment, error) {
	comments, err := r.queries.GetAllComments(ctx, post_id)
	if err != nil {
		log.Printf("Error retrieving all comments: %v", err)
		return nil, err
	}
	return comments, nil
}

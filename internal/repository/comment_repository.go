package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
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

func (r *CommentRepository) GetComment(ctx context.Context, commentID uuid.UUID) (entity.Comment, error) {
	commentRow, err := r.queries.GetComment(ctx, commentID)
	if err != nil {
			log.Printf("Error retrieving comment: %v", err)
			return entity.Comment{}, err
	}
	
	return entity.Comment{
			CommentID: commentRow.CommentID,
			UserID:    commentRow.UserID,
			PostID:    commentRow.PostID,
			Content:   commentRow.Content,
			CreatedAt: commentRow.CreatedAt.Time.UTC().Format(time.RFC3339),
			User: entity.User{
					UserID:   commentRow.UserID,
					Username: commentRow.Username,
					Avatar:   commentRow.Avatar.String,
			},
	}, nil
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

func (r *CommentRepository) GetAllComments(ctx context.Context, post_id uuid.UUID) ([]entity.Comment, error) {
	commentRows, err := r.queries.GetAllComments(ctx, post_id)
	if err != nil {
			log.Printf("Error retrieving all comments: %v", err)
			return nil, err
	}
	
	var comments []entity.Comment
	for _, row := range commentRows {
			comments = append(comments, entity.Comment{
					CommentID: row.CommentID,
					UserID:    row.UserID,
					PostID:    row.PostID,
					Content:   row.Content,
					CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
					UpdatedAt: row.UpdatedAt.Time.UTC().Format(time.RFC3339),
					User: entity.User{
							UserID:   row.UserID,
							Username: row.Username,
							Avatar:   row.Avatar.String,
					},
			})
	}
	
	return comments, nil
}

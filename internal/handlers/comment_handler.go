package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/services"
	"github.com/nycholasmarques/socialdev/utils"
)

type CommentHandler struct {
	commentService *services.CommentService
}
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "User ID not found in context", nil)
			return
	}

	postID := r.PathValue("post_id")
	if postID == "" {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Missing post_id in URL", nil)
			return
	}

	var req struct {
			Content string `json:"content" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
	}

	validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
			errorMessages := entity.GetValidationCommentErrors(err)
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload", errorMessages)
			return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID format", nil)
			return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid post ID format", nil)
			return
	}

	comment, err := h.commentService.CreateComment(userIdUUID, postIDUUID, req.Content)
	if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.NewErrorResponse(w, http.StatusNotFound, "Post not found", "")
				return
			}
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create comment", "")
			return
	}

	commentResponse := map[string]any{
			"comment_id": comment.CommentID,
			"user_id":    comment.UserID,
			"post_id":    comment.PostID,
			"content":    comment.Content,
			"created_at": comment.CreatedAt.Time.Format(time.RFC3339),
	}

	utils.NewMessageResponse(w, http.StatusCreated, "Comment created successfully", commentResponse)
}

func (h *CommentHandler) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	commentID := r.PathValue("comment_id")
	if commentID == "" {
			http.Error(w, `{"error": "Comment ID is required"}`, http.StatusBadRequest)
			return
	}

	commentIDUUID, err := uuid.Parse(commentID)
	if err != nil {
			http.Error(w, `{"error": "Invalid comment ID format"}`, http.StatusBadRequest)
			return
	}

	comment, err := h.commentService.GetComment(commentIDUUID)
	if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.NewErrorResponse(w, http.StatusNotFound, "Comment not found", "")
				return
			}
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get comment", "")
			return
	}

	response := map[string]any{
			"comment_id": comment.CommentID,
			"user_id":    comment.UserID,
			"post_id":    comment.PostID,
			"content":    comment.Content,
			"created_at": comment.CreatedAt,
			"user": map[string]any{
					"user_id":   comment.User.UserID,
					"username":  comment.User.Username,
					"avatar":    comment.User.Avatar,
			},
	}
	utils.NewMessageResponse(w, http.StatusOK, "Comment found successfully", response)
}

func (h *CommentHandler) GetAllCommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	postID := r.PathValue("post_id")
	if postID == "" {
			http.Error(w, `{"error": "Post ID is required"}`, http.StatusBadRequest)
			return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
			http.Error(w, `{"error": "Invalid post ID format"}`, http.StatusBadRequest)
			return
	}

	comments, err := h.commentService.GetAllComments(postIDUUID)
	if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
	}

	var response []map[string]interface{}
	for _, comment := range comments {
			response = append(response, map[string]interface{}{
					"comment_id": comment.CommentID,
					"user_id":    comment.UserID,
					"post_id":    comment.PostID,
					"content":    comment.Content,
					"created_at": comment.CreatedAt,
					"updated_at": comment.UpdatedAt,
					"user": map[string]interface{}{
							"user_id":   comment.User.UserID,
							"username":  comment.User.Username,
							"avatar":    comment.User.Avatar,
					},
			})
	}

	utils.NewMessageResponse(w, http.StatusOK, "Comments found successfully", response)
}

func (h *CommentHandler) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "User ID not found in context", nil)
			return
	}

	commentID := r.PathValue("comment_id")
	if commentID == "" {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Missing comment_id in URL", nil)
			return
	}

	commentIDUUID, err := uuid.Parse(commentID)
	if err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid comment ID format", nil)
			return
	}

	var req struct {
			Content string `json:"content" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
			errorMessages := entity.GetValidationCommentErrors(err)
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload", errorMessages)
			return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID format", nil)
			return
	}

	updatedComment, err := h.commentService.UpdateComment(userIdUUID, commentIDUUID, req.Content)
	if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
					utils.NewErrorResponse(w, http.StatusNotFound, "Comment not found or you don't have permission", nil)
					return
			}
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update comment: "+err.Error(), nil)
			return
	}

	commentResponse := map[string]interface{}{
			"comment_id": updatedComment.CommentID,
			"user_id":    updatedComment.UserID,
			"post_id":    updatedComment.PostID,
			"content":    updatedComment.Content,
			"created_at": updatedComment.CreatedAt.Time.Format(time.RFC3339),
			"updated_at": updatedComment.UpdatedAt.Time.Format(time.RFC3339),
	}

	utils.NewMessageResponse(w, http.StatusOK, "Comment updated successfully", commentResponse)
}

func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	commentID := r.PathValue("comment_id")

	commentIDUUID, err := uuid.Parse(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.commentService.DeleteComment(userIdUUID, commentIDUUID)
	if err != nil {
		http.Error(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.NewMessageResponse(w, http.StatusOK, "Comment deleted successfully", nil)
}

package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/services"
	"github.com/nycholasmarques/socialdev/utils"
)

type LikeHandler struct {
	likeService *services.LikeService
}

type LikeResponse struct {
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt string    `json:"created_at"`
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

func (h *LikeHandler) CreateLikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "User ID not found in context", "")
		return
	}

	postID := r.PathValue("post_id")

	if userID == "" || postID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id or like_id in URL", "")
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid post ID", "")
		return
	}

	isLiked, _ := h.likeService.GetLike(userIdUUID, postIDUUID)
	if isLiked != (db.Like{}) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Like already exists", "")
		return
	}

	_, err = h.likeService.CreateLike(userIdUUID, postIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create like", "")
		return
	}

	responseData := map[string]interface{}{
		"liked": true,
	}

	utils.NewMessageResponse(w, http.StatusOK, "Like created successfully!", responseData)
}

func (h *LikeHandler) GetLikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "User ID not found in context", "")
		return
	}

	postID := r.PathValue("post_id")

	if userID == "" || postID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id or like_id in URL", "")
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid post ID", "")
		return
	}

	isLiked, _ := h.likeService.GetLike(userIdUUID, postIDUUID)
	if isLiked == (db.Like{}) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Like does not exist", "")
		return
	}

	responseData := map[string]interface{}{
		"liked": true,
	}

	utils.NewMessageResponse(w, http.StatusOK, "Like found successfully!", responseData)
}

func (h *LikeHandler) DeleteLikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "User ID not found in context", "")
		return
	}

	postID := r.PathValue("post_id")

	if userID == "" || postID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id or like_id in URL", "")
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid post ID", "")
		return
	}

	isLiked, _ := h.likeService.GetLike(userIdUUID, postIDUUID)
	if isLiked == (db.Like{}) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Like does not exist", "")
		return
	}

	if err := h.likeService.DeleteLike(userIdUUID, postIDUUID); err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete like", "")
		return
	}

	responseData := map[string]interface{}{
		"liked": false,
	}

	utils.NewMessageResponse(w, http.StatusOK, "Like deleted successfully", responseData)
}

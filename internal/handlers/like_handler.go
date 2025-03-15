package handlers

import (
	"encoding/json"
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
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	postID := r.PathValue("post_id")

	if userID == "" || postID == "" {
		http.Error(w, "Missing user_id or like_id in URL", http.StatusBadRequest)
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	isLiked, _ := h.likeService.GetLike(userIdUUID, postIDUUID)

	if isLiked != (db.Like{}) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Like already exists", "")
		return
	}

	like, err := h.likeService.CreateLike(userIdUUID, postIDUUID)
	if err != nil {
		http.Error(w, "Failed to create user with like: "+err.Error(), http.StatusInternalServerError)
		return
	}

	likeResponse := LikeResponse{
		PostID:    like.PostID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt.Time.String(),
	}

	res, err := json.Marshal(likeResponse)
	if err != nil {
		http.Error(w, "Failed to marshal like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))

}

func (h *LikeHandler) DeleteLikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	postID := r.PathValue("post_id")

	if userID == "" || postID == "" {
		http.Error(w, "Missing user_id or like_id in URL", http.StatusBadRequest)
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	postIDUUID, err := uuid.Parse(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	isLiked, _ := h.likeService.GetLike(userIdUUID, postIDUUID)
	if isLiked == (db.Like{}) {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Like does not exist", "")
		return
	}

	if err := h.likeService.DeleteLike(userIdUUID, postIDUUID); err != nil {
		http.Error(w, "Failed to delete like: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.NewMessageResponse(w, http.StatusOK, "Like deleted successfully")
}

package favorite_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	favorite_services "github.com/nycholasmarques/socialdev/internal/modules/favorite/service"
	"github.com/nycholasmarques/socialdev/utils"
)

type FavoriteHandler struct {
	favoriteService favorite_services.FavoriteService
}

func NewFavoriteHandler(favoriteService favorite_services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func (h *FavoriteHandler) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Post ID is required", "")
		return
	}

	err := h.favoriteService.AddFavorite(userID, postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to add favorite", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Favorite added successfully"})
}

func (h *FavoriteHandler) RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Post ID is required", "")
		return
	}

	err := h.favoriteService.RemoveFavorite(userID, postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to remove favorite", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Favorite removed successfully"})
}

func (h *FavoriteHandler) GetUserFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	offset := int32(0)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	favorites, err := h.favoriteService.GetUserFavorites(userID, limit, offset)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get favorites", err.Error())
		return
	}

	res, err := json.Marshal(favorites)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to marshal favorites", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

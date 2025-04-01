package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	// configs "github.com/nycholasmarques/socialdev/config"
	"github.com/go-playground/validator/v10"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/services"
	"github.com/nycholasmarques/socialdev/utils"
)

// var cfg = configs.LoadConfig()

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

var validate *validator.Validate

func (ph *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", "")
		return
	}

	var postCreate entity.Post
	err = json.Unmarshal(reader, &postCreate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to unmarshal request body", "")
		return
	}

	validate = validator.New()
	err = validate.Struct(postCreate)
	if err != nil {
		errorMessages := entity.GetValidationPostErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	createPost, err := ph.postService.CreatePost(userID, postCreate.Title, postCreate.Content, postCreate.Photo, postCreate.Video)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Post created successfully",
		"post":    createPost,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	id := r.URL.Path
	postID := strings.TrimPrefix(id, "/posts/")

	post, err := ph.postService.GetPost(postID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "Post not found", "")
		return
	}

	response := map[string]interface{}{
		"post": post,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	// Pegar parâmetros da query string
	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")
	orderBy := r.URL.Query().Get("order_by") // "likes" ou "comments"
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	// Definir valores padrão
	limit := int32(10)  // default limit
	offset := int32(0)  // default offset
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

	var titlePtr, contentPtr *string
	if title != "" {
			titlePtr = &title
	}
	if content != "" {
			contentPtr = &content
	}

	posts, totalCount, err := ph.postService.GetAllPosts(titlePtr, contentPtr, orderBy, limit, offset, userID)
	if err != nil {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get posts", "")
			return
	}

	response := entity.Response{
			Posts:  posts,
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	accessToken, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || accessToken == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	// Extrai o user_id da URL
	vars := r.PathValue("user_id")
	if vars == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	userID := vars

	// Busca as postagens do usuário
	posts, err := ph.postService.GetAllUserPosts(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to get user posts", "")
		return
	}

	// Retorna as postagens
	response := map[string]interface{}{
		"user_id": userID,
		"posts":   posts,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	id := r.URL.Path
	postID := strings.TrimPrefix(id, "/posts/")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", "")
		return
	}

	var postUpdate entity.Post
	err = json.Unmarshal(reader, &postUpdate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to unmarshal request body", "")
		return
	}

	validate = validator.New()
	err = validate.Struct(postUpdate)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		var validationErrors []string
		for _, e := range errors {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
		}
		response, _ := json.Marshal(map[string]interface{}{"errors": validationErrors})
		http.Error(w, string(response), http.StatusBadRequest)
		return
	}

	createPost, err := ph.postService.UpdatePost(postID, userID, postUpdate.Title, postUpdate.Content, postUpdate.Photo, postUpdate.Video)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update post", "")
		return
	}

	response := map[string]interface{}{
		"message": "Post updated successfully",
		"post":    createPost,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	id := r.URL.Path
	postID := strings.TrimPrefix(id, "/posts/")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	err := ph.postService.DeletePost(postID, userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete post", "")
		return
	}

	utils.NewMessageResponse(w, http.StatusOK, "Post deleted successfully")
}

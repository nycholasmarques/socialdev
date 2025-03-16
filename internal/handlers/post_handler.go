package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var postCreate entity.Post
	err = json.Unmarshal(reader, &postCreate)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	validate = validator.New()
	err = validate.Struct(postCreate)
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
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
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

	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")

	var titleFilter, contentFilter *string
	if title != "" {
		titleFilter = &title
	}
	if content != "" {
		contentFilter = &content
	}

	posts, err := ph.postService.GetAllPosts(titleFilter, contentFilter)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "No posts found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"posts": posts,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ph *PostHandler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	accessToken, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || accessToken == "" {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Extrai o user_id da URL
	vars := r.PathValue("user_id")
	if vars == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID := vars

	// Busca as postagens do usuário
	posts, err := ph.postService.GetAllUserPosts(userID)
	if err != nil {
		http.Error(w, "Failed to get user posts", http.StatusInternalServerError)
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
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var postUpdate entity.Post
	err = json.Unmarshal(reader, &postUpdate)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
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

func (ph *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	id := r.URL.Path
	postID := strings.TrimPrefix(id, "/posts/")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	err := ph.postService.DeletePost(postID, userID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	utils.NewMessageResponse(w, http.StatusOK, "Post deleted successfully")
}
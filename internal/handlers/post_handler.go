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
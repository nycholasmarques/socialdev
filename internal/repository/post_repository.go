package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/entity"
)

type PostRepository struct {
	queries *db.Queries
}

func NewPostRepository(queries *db.Queries) *PostRepository {
	return &PostRepository{queries: queries}
}

func GetPostRepository(queries *db.Queries) *PostRepository {
	return &PostRepository{queries: queries}
}

func (r *PostRepository) CreatePost(ctx context.Context, params db.CreatePostParams) (entity.Post, error) {
	post, err := r.queries.CreatePost(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}

	var postResponse entity.Post

	postResponse.PostID = post.PostID.String()
	postResponse.UserID = post.UserID.String()
	postResponse.Title = post.Title.String
	postResponse.Content = post.Content
	postResponse.Photo = post.Photo.String
	postResponse.Video = post.Video.String
	postResponse.CreatedAt = post.CreatedAt.Time.UTC().Format(time.RFC3339)
	postResponse.UpdatedAt = post.UpdatedAt.Time.UTC().Format(time.RFC3339)

	return postResponse, err
}

func (r *PostRepository) GetPost(ctx context.Context, post_id uuid.UUID) (entity.Post, error) {
	post, err := r.queries.GetPost(ctx, post_id)
	if err != nil {
		log.Printf("%v", err)
		return entity.Post{}, err
	}

	var postResponse entity.Post

	postResponse.PostID = post.PostID.String()
	postResponse.UserID = post.UserID.String()
	postResponse.Title = post.Title.String
	postResponse.Content = post.Content
	postResponse.Photo = post.Photo.String
	postResponse.Video = post.Video.String
	postResponse.CreatedAt = post.CreatedAt.Time.String()
	postResponse.UpdatedAt = post.UpdatedAt.Time.String()

	return postResponse, err
}

func (r *PostRepository) GetAllPosts(ctx context.Context, title, content *string, orderBy string, limit, offset int32) ([]entity.Post, int64, error) {
	var titleFilter, contentFilter sql.NullString
	if title != nil {
			titleFilter = sql.NullString{String: *title, Valid: true}
	}
	if content != nil {
			contentFilter = sql.NullString{String: *content, Valid: true}
	}

	posts, err := r.queries.GetAllPosts(ctx, db.GetAllPostsParams{
			Title:   titleFilter,
			Content: contentFilter,
			OrderBy: sql.NullString{String: orderBy, Valid: true},
			Limit:   sql.NullInt32{Int32: limit, Valid: true},
			Offset:  sql.NullInt32{Int32: offset, Valid: true},
	})
	if err != nil {
			log.Printf("%v", err)
			return nil, 0, err
	}

	var postsResponse []entity.Post
	var totalCount int64
	for i, post := range posts {
			if i == 0 { // Pegamos o total_count da primeira linha
				totalCount = post.TotalCount
		}
			postsResponse = append(postsResponse, entity.Post{
					PostID:      post.PostID.String(),
					UserID:      post.UserID.String(),
					Title:       post.Title.String,
					Content:     post.Content,
					Photo:       post.Photo.String,
					Video:       post.Video.String,
					CreatedAt:   post.CreatedAt.Time.UTC().Format(time.RFC3339),
					UpdatedAt:   post.UpdatedAt.Time.UTC().Format(time.RFC3339),
					Like:        int(post.LikeCount),
					Views:       int(post.ViewsCount),
					Comments:    []entity.CommentPost{},
					User: entity.UserPost{
							UserID:   post.UserID.String(),
							Username: post.UserName.String,
							Avatar:   post.UserAvatar.String,
					},
			})
	}

	return postsResponse, totalCount, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, params db.UpdatePostParams) (entity.Post, error) {
	post, err := r.queries.UpdatePost(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}

	var postResponse entity.Post

	postResponse.PostID = post.PostID.String()
	postResponse.UserID = post.UserID.String()
	postResponse.Title = post.Title.String
	postResponse.Content = post.Content
	postResponse.Photo = post.Photo.String
	postResponse.Video = post.Video.String
	postResponse.CreatedAt = post.CreatedAt.Time.String()
	postResponse.UpdatedAt = post.UpdatedAt.Time.String()

	return postResponse, err
}

func (r *PostRepository) DeletePost(ctx context.Context, post_id uuid.UUID) error {
	err := r.queries.DeletePost(ctx, post_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}

func (r *PostRepository) GetAllUserPosts(ctx context.Context, user_id uuid.UUID) ([]entity.Post, error) {
	posts, err := r.queries.GetAllUserPosts(ctx, user_id)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var postsResponse []entity.Post

	for _, post := range posts {
		postID := post.PostID.String()

		postExists := false
		for _, existingPost := range postsResponse {
			if existingPost.PostID == postID {
				postExists = true
				break
			}
		}

		if !postExists {
			postsResponse = append(postsResponse, entity.Post{
				PostID:    postID,
				UserID:    post.UserID.String(),
				Title:     post.Title.String,
				Content:   post.Content,
				Photo:     post.Photo.String,
				Video:     post.Video.String,
				CreatedAt: post.CreatedAt.Time.UTC().Format(time.RFC3339),
				UpdatedAt: post.UpdatedAt.Time.UTC().Format(time.RFC3339),
				Like:      int(post.LikeCount),
				Views:     int(post.ViewsCount),
				User: entity.UserPost{
					UserID:   post.UserID.String(),
					Username: post.UserName.String,
					Avatar:   post.UserAvatar.String,
				},
				Comments: []entity.CommentPost{},
			})
		}

		if post.CommentID.Valid {
			for i, p := range postsResponse {
				if p.PostID == postID {
					postsResponse[i].Comments = append(postsResponse[i].Comments, entity.CommentPost{
						CommentID: post.CommentID.UUID.String(),
						Content:   post.CommentContent.String,
						CreatedAt: post.CommentCreatedAt.Time.String(),
						User: entity.UserPost{
							UserID:   post.CommentUserID.UUID.String(),
							Username: post.CommentUserName.String,
							Avatar:   post.CommentUserAvatar.String,
						},
					})
				}
			}
		}
	}

	return postsResponse, nil
}

func (r *PostRepository) GetUserPost(ctx context.Context, params db.GetUserPostParams) (entity.Post, error) {
	post, err := r.queries.GetUserPost(ctx, params)
	if err != nil {
		log.Printf("%v", err)
		return entity.Post{}, err
	}

	var postResponse entity.Post
	postResponse.PostID = post.PostID.String()
	postResponse.UserID = post.UserID.String()
	postResponse.Title = post.Title.String
	postResponse.Content = post.Content
	postResponse.Photo = post.Photo.String
	postResponse.Video = post.Video.String
	postResponse.CreatedAt = post.CreatedAt.Time.String()
	postResponse.UpdatedAt = post.UpdatedAt.Time.String()

	return postResponse, err
}

package entity

type FavoriteResponse struct {
    PostID            string `json:"post_id"`
    UserID            string `json:"user_id"`
    Title             string `json:"title"`
    Content           string `json:"content"`
    Photo             string `json:"photo,omitempty"`
    Video             string `json:"video,omitempty"`
    PostCreatedAt     string `json:"post_created_at"`
    UpdatedAt         string `json:"updated_at"`
    FavoriteCreatedAt string `json:"favorite_created_at"`
    UserName          string `json:"user_name"`
    UserAvatar        string `json:"user_avatar"`
    LikeCount         int64  `json:"like_count"`
    ViewsCount        int64  `json:"views_count"`
    CommentCount      int64  `json:"comment_count"`
    LikedByUser       bool   `json:"liked_by_user"`
}
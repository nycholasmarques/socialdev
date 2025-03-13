package entity

type Post struct {
	PostID    string `json:"post_id" validate:"omitempty,uuid4"`
	UserID    string `json:"user_id" validate:"omitempty,uuid4"`
	Title     string `json:"title" validate:"omitempty,min=5,max=100"`
	Content   string `json:"content" validate:"required,min=10"`
	Photo     string `json:"photo" validate:"omitempty,url"` 
	Video     string `json:"video" validate:"omitempty,url"`
	CreatedAt string `json:"created_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	UpdatedAt string `json:"updated_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

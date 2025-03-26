// internal/models/user.go
package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	GoogleID  string    `gorm:"unique"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Avatar    string
	Bio       string
	Github    string
	Linkedin  string
	Website   string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Roles     []Role    `gorm:"many2many:user_roles"`
	Posts     []Post
	Comments  []Comment
	Likes     []Like
}
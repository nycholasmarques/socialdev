package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Title     string
	Content   string `gorm:"not null"`
	Photo     string
	Video     string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Comments  []Comment
	Likes     []Like
}
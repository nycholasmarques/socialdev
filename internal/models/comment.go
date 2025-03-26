package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	CommentID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	PostID    uuid.UUID `gorm:"not null"`
	Post      Post      `gorm:"foreignKey:PostID"`
	Content   string    `gorm:"not null"`
	Visible   bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
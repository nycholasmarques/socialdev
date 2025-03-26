package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	PostID    uuid.UUID `gorm:"primaryKey"`
	User      User      `gorm:"foreignKey:UserID"`
	Post      Post      `gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
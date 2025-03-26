package models

import (
	"time"
	"github.com/google/uuid"
)

type Role struct {
	RoleID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Users     []User    `gorm:"many2many:user_roles"`
}
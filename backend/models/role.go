package models

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;unique;not null" json:"name"`
	AccessLevel int       `gorm:"not null" json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

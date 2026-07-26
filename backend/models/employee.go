package models

import (
	"time"
)

type Employee struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	FirstName    string     `gorm:"size:100" json:"first_name"`
	LastName     string     `gorm:"size:100" json:"last_name"`
	RoleID       uint       `gorm:"not null" json:"role_id"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"role"`
	ManagerID    *uint      `gorm:"index" json:"manager_id"`
	Manager      *Employee  `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	DirectReports []Employee `gorm:"foreignKey:ManagerID" json:"direct_reports,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

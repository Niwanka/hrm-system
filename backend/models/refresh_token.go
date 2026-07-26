package models

import (
	"time"
)

type RefreshToken struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"-"`
	TokenHash  string    `gorm:"size:255;not null;index" json:"-"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
	Revoked    bool      `gorm:"default:false" json:"revoked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

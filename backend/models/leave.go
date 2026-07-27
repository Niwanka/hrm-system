package models

import (
	"time"
)

type LeaveType struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:50;unique;not null" json:"name"` // Annual, Casual, Sick, Maternity
	MaxDaysPerYear int       `gorm:"not null" json:"max_days_per_year"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LeaveRequest struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   uint       `gorm:"not null;index" json:"employee_id"`
	Employee     Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
	LeaveTypeID  uint       `gorm:"not null" json:"leave_type_id"`
	LeaveType    LeaveType  `gorm:"foreignKey:LeaveTypeID" json:"leave_type"`
	StartDate    time.Time  `gorm:"not null" json:"start_date"`
	EndDate      time.Time  `gorm:"not null" json:"end_date"`
	DaysCount    int        `gorm:"not null" json:"days_count"`
	Reason       string     `gorm:"size:255" json:"reason"`
	Status       string     `gorm:"size:30;default:'Pending'" json:"status"` // Pending, Approved, Rejected
	ApprovedByID *uint      `gorm:"index" json:"approved_by_id"`
	ApprovedBy   *Employee  `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LeaveBalance struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     uint      `gorm:"not null;index" json:"employee_id"`
	LeaveTypeID    uint      `gorm:"not null" json:"leave_type_id"`
	LeaveType      LeaveType `gorm:"foreignKey:LeaveTypeID" json:"leave_type"`
	AllocatedDays  int       `gorm:"not null" json:"allocated_days"`
	UsedDays       int       `gorm:"default:0" json:"used_days"`
	RemainingDays  int       `gorm:"not null" json:"remaining_days"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

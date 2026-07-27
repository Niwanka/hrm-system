package models

import (
	"time"
)

type AttendanceLog struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EmployeeID  uint       `gorm:"not null;index" json:"employee_id"`
	Employee    Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
	Date        time.Time  `gorm:"type:date;not null;index" json:"date"`
	ClockIn     time.Time  `gorm:"not null" json:"clock_in"`
	ClockOut    *time.Time `json:"clock_out,omitempty"`
	HoursWorked float64    `gorm:"default:0" json:"hours_worked"`
	Status      string     `gorm:"size:30;default:'Present'" json:"status"` // Present, Late, Half Day, Absent
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

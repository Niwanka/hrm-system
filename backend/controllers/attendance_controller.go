package controllers

import (
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"

	"hrm-backend/config"
	"hrm-backend/middleware"
	"hrm-backend/models"
)

func GetTodayAttendance(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	today := time.Now().Format("2006-01-02")
	var logRecord models.AttendanceLog
	err := config.DB.Where("employee_id = ? AND DATE(date) = ?", claims.EmployeeID, today).First(&logRecord).Error

	if err != nil {
		return c.JSON(fiber.Map{
			"is_clocked_in":  false,
			"is_clocked_out": false,
			"log":            nil,
		})
	}

	isClockedIn := logRecord.ClockOut == nil
	isClockedOut := logRecord.ClockOut != nil

	return c.JSON(fiber.Map{
		"is_clocked_in":  isClockedIn,
		"is_clocked_out": isClockedOut,
		"log":            logRecord,
	})
}

func ClockIn(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	now := time.Now()
	todayStr := now.Format("2006-01-02")

	var existing models.AttendanceLog
	if err := config.DB.Where("employee_id = ? AND DATE(date) = ?", claims.EmployeeID, todayStr).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You have already clocked in for today!",
		})
	}

	// Status calculation: Late if after 09:30 AM
	status := "Present"
	if now.Hour() > 9 || (now.Hour() == 9 && now.Minute() > 30) {
		status = "Late"
	}

	todayDate, _ := time.Parse("2006-01-02", todayStr)

	logRecord := models.AttendanceLog{
		EmployeeID: claims.EmployeeID,
		Date:       todayDate,
		ClockIn:    now,
		Status:     status,
	}

	if err := config.DB.Create(&logRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record clock-in"})
	}

	empID := claims.EmployeeID
	middleware.LogAudit(c, &empID, "CLOCK_IN", fmt.Sprintf("Clocked in at %s (Status: %s)", now.Format("15:04:05"), status))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("Clocked in successfully at %s (%s)", now.Format("15:04:05"), status),
		"log":     logRecord,
	})
}

func ClockOut(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	now := time.Now()
	todayStr := now.Format("2006-01-02")

	var logRecord models.AttendanceLog
	if err := config.DB.Where("employee_id = ? AND DATE(date) = ? AND clock_out IS NULL", claims.EmployeeID, todayStr).First(&logRecord).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No active clock-in session found for today.",
		})
	}

	logRecord.ClockOut = &now
	hours := math.Round(now.Sub(logRecord.ClockIn).Hours()*100) / 100
	logRecord.HoursWorked = hours

	config.DB.Save(&logRecord)

	empID := claims.EmployeeID
	middleware.LogAudit(c, &empID, "CLOCK_OUT", fmt.Sprintf("Clocked out at %s (Total: %.2f hrs)", now.Format("15:04:05"), hours))

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Clocked out successfully at %s (Worked: %.2f hrs)", now.Format("15:04:05"), hours),
		"log":     logRecord,
	})
}

func GetMyAttendanceLogs(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var logs []models.AttendanceLog
	config.DB.Where("employee_id = ?", claims.EmployeeID).
		Order("date desc").
		Limit(30).
		Find(&logs)

	return c.JSON(fiber.Map{
		"logs": logs,
	})
}

func GetTeamAttendanceLogs(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var logs []models.AttendanceLog
	query := config.DB.Preload("Employee.Role").Preload("Employee.Department")

	if claims.AccessLevel < 80 {
		var directReportIDs []uint
		config.DB.Model(&models.Employee{}).Where("manager_id = ?", claims.EmployeeID).Pluck("id", &directReportIDs)
		directReportIDs = append(directReportIDs, claims.EmployeeID)
		query = query.Where("employee_id IN ?", directReportIDs)
	}

	if err := query.Order("date desc").Limit(100).Find(&logs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch team logs"})
	}

	return c.JSON(fiber.Map{
		"team_logs": logs,
	})
}

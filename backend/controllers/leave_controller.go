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

type SubmitLeaveRequestPayload struct {
	LeaveTypeID uint   `json:"leave_type_id"`
	StartDate   string `json:"start_date"` // YYYY-MM-DD
	EndDate     string `json:"end_date"`   // YYYY-MM-DD
	Reason      string `json:"reason"`
}

func GetMyLeaveData(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Ensure LeaveBalance is initialized for this user for all leave types
	ensureLeaveBalancesInitialized(claims.EmployeeID)

	var balances []models.LeaveBalance
	config.DB.Preload("LeaveType").Where("employee_id = ?", claims.EmployeeID).Find(&balances)

	var requests []models.LeaveRequest
	config.DB.Preload("LeaveType").Preload("ApprovedBy").
		Where("employee_id = ?", claims.EmployeeID).
		Order("created_at desc").
		Find(&requests)

	var leaveTypes []models.LeaveType
	config.DB.Find(&leaveTypes)

	return c.JSON(fiber.Map{
		"balances":    balances,
		"requests":    requests,
		"leave_types": leaveTypes,
	})
}

func SubmitLeaveRequest(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req SubmitLeaveRequestPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if req.LeaveTypeID == 0 || req.StartDate == "" || req.EndDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Leave type, start date, and end date are required"})
	}

	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil || endDate.Before(startDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start or end date selection"})
	}

	// Calculate total calendar days requested
	days := int(math.Ceil(endDate.Sub(startDate).Hours()/24)) + 1
	if days <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Leave duration must be at least 1 day"})
	}

	// Check Leave Balance
	ensureLeaveBalancesInitialized(claims.EmployeeID)
	var balance models.LeaveBalance
	err := config.DB.Where("employee_id = ? AND leave_type_id = ?", claims.EmployeeID, req.LeaveTypeID).First(&balance).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Leave balance record not found"})
	}

	if balance.RemainingDays < days {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Insufficient leave balance. Requested %d day(s), but only %d day(s) remaining.", days, balance.RemainingDays),
		})
	}

	leaveReq := models.LeaveRequest{
		EmployeeID:  claims.EmployeeID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   startDate,
		EndDate:     endDate,
		DaysCount:   days,
		Reason:      req.Reason,
		Status:      "Pending",
	}

	if err := config.DB.Create(&leaveReq).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to submit leave request"})
	}

	empID := claims.EmployeeID
	middleware.LogAudit(c, &empID, "SUBMIT_LEAVE_REQUEST", fmt.Sprintf("Submitted %d-day leave request (%s to %s)", days, req.StartDate, req.EndDate))

	config.DB.Preload("LeaveType").First(&leaveReq, leaveReq.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Leave request submitted successfully",
		"request": leaveReq,
	})
}

func GetPendingLeaveRequests(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var requests []models.LeaveRequest
	query := config.DB.Preload("Employee.Role").Preload("Employee.Department").Preload("LeaveType").Where("status = ?", "Pending")

	// If Manager (Level 50), filter by direct reports unless HR/Admin (Level 80+)
	if claims.AccessLevel < 80 {
		var directReportIDs []uint
		config.DB.Model(&models.Employee{}).Where("manager_id = ?", claims.EmployeeID).Pluck("id", &directReportIDs)
		if len(directReportIDs) == 0 {
			return c.JSON(fiber.Map{"pending_requests": []models.LeaveRequest{}})
		}
		query = query.Where("employee_id IN ?", directReportIDs)
	}

	if err := query.Order("created_at asc").Find(&requests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch pending requests"})
	}

	return c.JSON(fiber.Map{
		"pending_requests": requests,
	})
}

func ApproveLeaveRequest(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var leaveReq models.LeaveRequest
	if err := config.DB.Preload("LeaveType").First(&leaveReq, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Leave request not found"})
	}

	if leaveReq.Status != "Pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Leave request is already processed"})
	}

	// Update Leave Request Status
	leaveReq.Status = "Approved"
	leaveReq.ApprovedByID = &claims.EmployeeID
	config.DB.Save(&leaveReq)

	// Update Employee Leave Balance
	var balance models.LeaveBalance
	if err := config.DB.Where("employee_id = ? AND leave_type_id = ?", leaveReq.EmployeeID, leaveReq.LeaveTypeID).First(&balance).Error; err == nil {
		balance.UsedDays += leaveReq.DaysCount
		balance.RemainingDays = balance.AllocatedDays - balance.UsedDays
		if balance.RemainingDays < 0 {
			balance.RemainingDays = 0
		}
		config.DB.Save(&balance)
	}

	approverID := claims.EmployeeID
	middleware.LogAudit(c, &approverID, "APPROVE_LEAVE_REQUEST", fmt.Sprintf("Approved %d-day leave request ID %d for employee ID %d", leaveReq.DaysCount, leaveReq.ID, leaveReq.EmployeeID))

	return c.JSON(fiber.Map{
		"message": "Leave request approved successfully",
		"request": leaveReq,
	})
}

func RejectLeaveRequest(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var leaveReq models.LeaveRequest
	if err := config.DB.First(&leaveReq, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Leave request not found"})
	}

	if leaveReq.Status != "Pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Leave request is already processed"})
	}

	leaveReq.Status = "Rejected"
	leaveReq.ApprovedByID = &claims.EmployeeID
	config.DB.Save(&leaveReq)

	approverID := claims.EmployeeID
	middleware.LogAudit(c, &approverID, "REJECT_LEAVE_REQUEST", fmt.Sprintf("Rejected leave request ID %d for employee ID %d", leaveReq.ID, leaveReq.EmployeeID))

	return c.JSON(fiber.Map{
		"message": "Leave request rejected",
		"request": leaveReq,
	})
}

func ensureLeaveBalancesInitialized(employeeID uint) {
	var leaveTypes []models.LeaveType
	config.DB.Find(&leaveTypes)

	for _, lt := range leaveTypes {
		var count int64
		config.DB.Model(&models.LeaveBalance{}).Where("employee_id = ? AND leave_type_id = ?", employeeID, lt.ID).Count(&count)
		if count == 0 {
			config.DB.Create(&models.LeaveBalance{
				EmployeeID:    employeeID,
				LeaveTypeID:   lt.ID,
				AllocatedDays: lt.MaxDaysPerYear,
				UsedDays:      0,
				RemainingDays: lt.MaxDaysPerYear,
			})
		}
	}
}

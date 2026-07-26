package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"hrm-backend/config"
	"hrm-backend/middleware"
	"hrm-backend/models"
)

type CreateEmployeeRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	RoleID       uint   `json:"role_id"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
	ManagerID    *uint  `json:"manager_id"`
}

type UpdateEmployeeRequest struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	RoleID       uint   `json:"role_id"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
	ManagerID    *uint  `json:"manager_id"`
}

func GetEmployeeHierarchy(c *fiber.Ctx) error {
	var topLevelManagers []models.Employee

	err := config.DB.Preload("Role").Preload("Department").
		Preload("DirectReports.Role").Preload("DirectReports.Department").
		Preload("DirectReports.DirectReports.Role").Preload("DirectReports.DirectReports.Department").
		Where("manager_id IS NULL").
		Find(&topLevelManagers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch employee hierarchy",
		})
	}

	return c.JSON(fiber.Map{
		"hierarchy": topLevelManagers,
	})
}

func GetEmployees(c *fiber.Ctx) error {
	var employees []models.Employee

	err := config.DB.Preload("Role").
		Preload("Department").
		Preload("Manager.Role").
		Find(&employees).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch employees",
		})
	}

	return c.JSON(fiber.Map{
		"employees": employees,
	})
}

func GetEmployeeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.Employee

	err := config.DB.Preload("Role").
		Preload("Department").
		Preload("Manager.Role").
		Preload("DirectReports.Role").
		First(&employee, id).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employee not found",
		})
	}

	return c.JSON(fiber.Map{
		"employee": employee,
	})
}

func CreateEmployee(c *fiber.Ctx) error {
	var req CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.RoleID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First name, last name, email, password, and role are required",
		})
	}

	status := req.Status
	if status == "" {
		status = "Active"
	}

	// Check if email already exists
	var existing models.Employee
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "An employee with this email address already exists",
		})
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	newEmployee := models.Employee{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
		Status:       status,
		ManagerID:    req.ManagerID,
	}

	if err := config.DB.Create(&newEmployee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert employee into database",
		})
	}

	var creatorID *uint
	if claims, ok := c.Locals("user").(*middleware.UserClaims); ok {
		creatorID = &claims.EmployeeID
	}
	middleware.LogAudit(c, creatorID, "CREATE_EMPLOYEE", fmt.Sprintf("Created employee %s %s (%s)", newEmployee.FirstName, newEmployee.LastName, newEmployee.Email))

	// Preload relations for response
	config.DB.Preload("Role").Preload("Department").Preload("Manager.Role").First(&newEmployee, newEmployee.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Employee inserted successfully",
		"employee": newEmployee,
	})
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employee record not found",
		})
	}

	// Email uniqueness check if modified
	if req.Email != "" && req.Email != employee.Email {
		var existing models.Employee
		if err := config.DB.Where("email = ? AND id != ?", req.Email, employee.ID).First(&existing).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email address is already in use by another employee",
			})
		}
		employee.Email = req.Email
	}

	if req.FirstName != "" {
		employee.FirstName = req.FirstName
	}
	if req.LastName != "" {
		employee.LastName = req.LastName
	}
	if req.RoleID != 0 {
		employee.RoleID = req.RoleID
	}
	employee.DepartmentID = req.DepartmentID
	if req.Status != "" {
		employee.Status = req.Status
	}
	employee.ManagerID = req.ManagerID

	if err := config.DB.Save(&employee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update employee details",
		})
	}

	var editorID *uint
	if claims, ok := c.Locals("user").(*middleware.UserClaims); ok {
		editorID = &claims.EmployeeID
	}
	middleware.LogAudit(c, editorID, "UPDATE_EMPLOYEE", fmt.Sprintf("Updated employee ID %d (%s %s)", employee.ID, employee.FirstName, employee.LastName))

	config.DB.Preload("Role").Preload("Department").Preload("Manager.Role").First(&employee, employee.ID)

	return c.JSON(fiber.Map{
		"message":  "Employee details updated successfully",
		"employee": employee,
	})
}

func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employee record not found",
		})
	}

	// GORM Soft Delete
	if err := config.DB.Delete(&employee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete employee record",
		})
	}

	var deleterID *uint
	if claims, ok := c.Locals("user").(*middleware.UserClaims); ok {
		deleterID = &claims.EmployeeID
	}
	middleware.LogAudit(c, deleterID, "DELETE_EMPLOYEE", fmt.Sprintf("Soft-deleted employee ID %d (%s %s)", employee.ID, employee.FirstName, employee.LastName))

	return c.JSON(fiber.Map{
		"message": "Employee soft-deleted successfully",
	})
}

func GetRoles(c *fiber.Ctx) error {
	var roles []models.Role
	if err := config.DB.Order("access_level desc").Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}
	return c.JSON(fiber.Map{
		"roles": roles,
	})
}

func GetDepartments(c *fiber.Ctx) error {
	var departments []models.Department
	if err := config.DB.Order("name asc").Find(&departments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch departments",
		})
	}
	return c.JSON(fiber.Map{
		"departments": departments,
	})
}

package controllers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"hrm-backend/config"
	"hrm-backend/models"
)

type CreateEmployeeRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    uint   `json:"role_id"`
	ManagerID *uint  `json:"manager_id"`
}

func GetEmployeeHierarchy(c *fiber.Ctx) error {
	var topLevelManagers []models.Employee

	// Fetch top-level employees (where ManagerID is NULL) and recursively preload DirectReports
	err := config.DB.Preload("Role").
		Preload("DirectReports.Role").
		Preload("DirectReports.DirectReports.Role").
		Preload("DirectReports.DirectReports.DirectReports.Role").
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
		ManagerID:    req.ManagerID,
	}

	if err := config.DB.Create(&newEmployee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert employee into database",
		})
	}

	// Preload relations for response
	config.DB.Preload("Role").Preload("Manager.Role").First(&newEmployee, newEmployee.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Employee inserted successfully",
		"employee": newEmployee,
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

package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"hrm-backend/config"
	"hrm-backend/middleware"
	"hrm-backend/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	var employee models.Employee
	err := config.DB.Preload("Role").Where("email = ?", req.Email).First(&employee).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	tokenString, err := middleware.GenerateJWT(
		employee.ID,
		employee.Role.AccessLevel,
		employee.Email,
		employee.FirstName,
		employee.LastName,
		employee.Role.Name,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Set HTTP-Only Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
		"user": fiber.Map{
			"id":           employee.ID,
			"email":        employee.Email,
			"first_name":   employee.FirstName,
			"last_name":    employee.LastName,
			"role_id":      employee.RoleID,
			"role_name":    employee.Role.Name,
			"access_level": employee.Role.AccessLevel,
		},
	})
}

func Me(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.UserClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var employee models.Employee
	err := config.DB.Preload("Role").Preload("Manager.Role").First(&employee, claims.EmployeeID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employee profile not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":           employee.ID,
			"email":        employee.Email,
			"first_name":   employee.FirstName,
			"last_name":    employee.LastName,
			"role_id":      employee.RoleID,
			"role_name":    employee.Role.Name,
			"access_level": employee.Role.AccessLevel,
			"manager_id":   employee.ManagerID,
			"manager_name": getManagerName(employee.Manager),
		},
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func getManagerName(manager *models.Employee) string {
	if manager == nil {
		return "None (Top Level)"
	}
	return manager.FirstName + " " + manager.LastName
}

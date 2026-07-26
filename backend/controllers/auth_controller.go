package controllers

import (
	"fmt"
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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		middleware.LogAudit(c, nil, "LOGIN_FAILED", "Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Email == "" || req.Password == "" {
		middleware.LogAudit(c, nil, "LOGIN_FAILED", "Missing email or password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	var employee models.Employee
	err := config.DB.Preload("Role").Where("email = ?", req.Email).First(&employee).Error
	if err != nil {
		middleware.LogAudit(c, nil, "LOGIN_FAILED", fmt.Sprintf("Email not found: %s", req.Email))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(req.Password))
	if err != nil {
		middleware.LogAudit(c, &employee.ID, "LOGIN_FAILED", "Invalid password provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// 1. Generate 15-Minute Short-Lived Access Token
	accessToken, err := middleware.GenerateAccessToken(
		employee.ID,
		employee.Role.AccessLevel,
		employee.Email,
		employee.FirstName,
		employee.LastName,
		employee.Role.Name,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	// 2. Generate 7-Day Long-Lived Refresh Token
	refreshTokenString, err := middleware.GenerateRefreshTokenString()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// 3. Store Refresh Token Hash in Database
	refreshTokenHash := middleware.HashToken(refreshTokenString)
	tokenRecord := models.RefreshToken{
		EmployeeID: employee.ID,
		TokenHash:  refreshTokenHash,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 Days
		Revoked:    false,
	}

	if err := config.DB.Create(&tokenRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save refresh token session",
		})
	}

	// 4. Set HTTP-Only Cookies
	accessCookie := new(fiber.Cookie)
	accessCookie.Name = "jwt"
	accessCookie.Value = accessToken
	accessCookie.Expires = time.Now().Add(15 * time.Minute)
	accessCookie.HTTPOnly = true
	accessCookie.SameSite = "Lax"
	c.Cookie(accessCookie)

	refreshCookie := new(fiber.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = refreshTokenString
	refreshCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	refreshCookie.HTTPOnly = true
	refreshCookie.SameSite = "Lax"
	c.Cookie(refreshCookie)

	// Log Audit Event
	middleware.LogAudit(c, &employee.ID, "LOGIN_SUCCESS", fmt.Sprintf("User logged in successfully (Role: %s)", employee.Role.Name))

	return c.JSON(fiber.Map{
		"message":       "Login successful",
		"token":         accessToken,
		"refresh_token": refreshTokenString,
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

func RefreshToken(c *fiber.Ctx) error {
	var refreshTokenString string

	// 1. Read from Cookie
	refreshTokenString = c.Cookies("refresh_token")

	// 2. Read from Request Body fallback
	if refreshTokenString == "" {
		var req RefreshRequest
		if err := c.BodyParser(&req); err == nil {
			refreshTokenString = req.RefreshToken
		}
	}

	if refreshTokenString == "" {
		middleware.LogAudit(c, nil, "REFRESH_TOKEN_FAILED", "No refresh token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token required",
		})
	}

	tokenHash := middleware.HashToken(refreshTokenString)

	var tokenRecord models.RefreshToken
	err := config.DB.Preload("Employee.Role").
		Where("token_hash = ? AND revoked = ? AND expires_at > ?", tokenHash, false, time.Now()).
		First(&tokenRecord).Error

	if err != nil {
		middleware.LogAudit(c, nil, "REFRESH_TOKEN_FAILED", "Invalid, revoked, or expired refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid, revoked, or expired refresh token. Please login again.",
		})
	}

	employee := tokenRecord.Employee

	// Issue New 15-Minute Access Token
	newAccessToken, err := middleware.GenerateAccessToken(
		employee.ID,
		employee.Role.AccessLevel,
		employee.Email,
		employee.FirstName,
		employee.LastName,
		employee.Role.Name,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new access token",
		})
	}

	// Update Access Cookie
	accessCookie := new(fiber.Cookie)
	accessCookie.Name = "jwt"
	accessCookie.Value = newAccessToken
	accessCookie.Expires = time.Now().Add(15 * time.Minute)
	accessCookie.HTTPOnly = true
	accessCookie.SameSite = "Lax"
	c.Cookie(accessCookie)

	middleware.LogAudit(c, &employee.ID, "TOKEN_REFRESH_SUCCESS", "Access token renewed successfully via refresh token")

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"token":   newAccessToken,
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
	var employeeID *uint

	// Extract claims if present
	if claims, ok := c.Locals("user").(*middleware.UserClaims); ok {
		employeeID = &claims.EmployeeID
	}

	// Revoke Refresh Token in Database
	refreshTokenString := c.Cookies("refresh_token")
	if refreshTokenString != "" {
		tokenHash := middleware.HashToken(refreshTokenString)
		config.DB.Model(&models.RefreshToken{}).
			Where("token_hash = ?", tokenHash).
			Update("revoked", true)
	}

	// Clear Cookies
	accessCookie := new(fiber.Cookie)
	accessCookie.Name = "jwt"
	accessCookie.Value = ""
	accessCookie.Expires = time.Now().Add(-1 * time.Hour)
	accessCookie.HTTPOnly = true
	c.Cookie(accessCookie)

	refreshCookie := new(fiber.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = ""
	refreshCookie.Expires = time.Now().Add(-1 * time.Hour)
	refreshCookie.HTTPOnly = true
	c.Cookie(refreshCookie)

	middleware.LogAudit(c, employeeID, "LOGOUT", "User logged out successfully")

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func GetAuditLogs(c *fiber.Ctx) error {
	var logs []models.AuditLog

	err := config.DB.Preload("Employee.Role").
		Order("created_at desc").
		Limit(50).
		Find(&logs).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch audit logs",
		})
	}

	return c.JSON(fiber.Map{
		"logs": logs,
	})
}

func getManagerName(manager *models.Employee) string {
	if manager == nil {
		return "None (Top Level)"
	}
	return manager.FirstName + " " + manager.LastName
}

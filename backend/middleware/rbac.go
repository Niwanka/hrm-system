package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"hrm-backend/config"
	"hrm-backend/models"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "hrm-super-secret-jwt-key-2026-cyan-mode"
}

type UserClaims struct {
	EmployeeID  uint   `json:"employee_id"`
	AccessLevel int    `json:"access_level"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	RoleName    string `json:"role_name"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a short-lived 15-minute Access Token
func GenerateAccessToken(employeeID uint, accessLevel int, email, firstName, lastName, roleName string) (string, error) {
	claims := UserClaims{
		EmployeeID:  employeeID,
		AccessLevel: accessLevel,
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		RoleName:    roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15-minute short-lived access
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshTokenString generates a secure 32-byte crypto-random string
func GenerateRefreshTokenString() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken generates SHA-256 hash of a token string for safe DB storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// LogAudit records security and system audit trail events
func LogAudit(c *fiber.Ctx, employeeID *uint, action, details string) {
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	auditEntry := models.AuditLog{
		EmployeeID: employeeID,
		Action:     action,
		IPAddress:  ip,
		UserAgent:  userAgent,
		Details:    details,
		CreatedAt:  time.Now(),
	}

	if err := config.DB.Create(&auditEntry).Error; err != nil {
		log.Printf("Failed to record audit log: %v", err)
	}
}

func AuthenticateJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// 1. Check Cookie
		tokenString = c.Cookies("jwt")

		// 2. Fallback to Authorization Header
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required. No valid access token found.",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Access token expired or invalid",
				"code":  "TOKEN_EXPIRED",
			})
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func RequireAccessLevel(minLevel int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userLocal := c.Locals("user")
		if userLocal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized. User context missing.",
			})
		}

		claims, ok := userLocal.(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse user claims",
			})
		}

		if claims.AccessLevel < minLevel {
			LogAudit(c, &claims.EmployeeID, "ACCESS_DENIED", fmt.Sprintf("Attempted access to route requiring level %d", minLevel))
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Forbidden: Requires minimum access level %d (your level: %d)", minLevel, claims.AccessLevel),
			})
		}

		return c.Next()
	}
}

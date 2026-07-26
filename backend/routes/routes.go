package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"hrm-backend/controllers"
	"hrm-backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rate Limiter for Login Endpoint (Max 5 attempts per minute per IP)
	loginLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			middleware.LogAudit(c, nil, "RATE_LIMIT_EXCEEDED", "Too many login attempts from IP")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many login attempts. Please wait 1 minute before trying again.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
		},
	})

	// Public Auth Endpoints
	api.Post("/login", loginLimiter, controllers.Login)
	api.Post("/refresh", controllers.RefreshToken)

	// Protected Endpoints
	protected := api.Group("", middleware.AuthenticateJWT())

	protected.Post("/logout", controllers.Logout)
	protected.Get("/me", controllers.Me)
	protected.Get("/roles", middleware.RequireAccessLevel(10), controllers.GetRoles)
	protected.Get("/departments", middleware.RequireAccessLevel(10), controllers.GetDepartments)

	// Employees Directory & Hierarchy (Access Level 10+)
	employees := protected.Group("/employees", middleware.RequireAccessLevel(10))
	employees.Get("/", controllers.GetEmployees)
	employees.Get("/hierarchy", controllers.GetEmployeeHierarchy)
	employees.Get("/:id", controllers.GetEmployeeByID)
	
	// Create & Update Employee Endpoints (Requires Access Level 50+ Manager/HR/Admin)
	employees.Post("/", middleware.RequireAccessLevel(50), controllers.CreateEmployee)
	employees.Put("/:id", middleware.RequireAccessLevel(50), controllers.UpdateEmployee)

	// Delete/Terminate Employee Endpoint (Requires Access Level 80+ HR/Admin)
	employees.Delete("/:id", middleware.RequireAccessLevel(80), controllers.DeleteEmployee)

	// Audit Trail Logs Endpoint (Requires Access Level 80+ for HR/Admin)
	auditLogs := protected.Group("/audit-logs", middleware.RequireAccessLevel(80))
	auditLogs.Get("/", controllers.GetAuditLogs)

	// Payroll Endpoint (Access Level 80+ for HR/Admin)
	payroll := protected.Group("/payroll", middleware.RequireAccessLevel(80))
	payroll.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Company Payroll Access Granted",
			"payroll_records": []fiber.Map{
				{"period": "2026-Q2", "total_disbursed": "$1,450,000", "status": "Processed"},
				{"period": "2026-Q3", "total_disbursed": "$1,480,000", "status": "Pending Approval"},
			},
		})
	})
}

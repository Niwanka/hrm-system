package routes

import (
	"github.com/gofiber/fiber/v2"

	"hrm-backend/controllers"
	"hrm-backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Public Auth Endpoints
	api.Post("/login", controllers.Login)

	// Protected Endpoints
	protected := api.Group("", middleware.AuthenticateJWT())

	protected.Post("/logout", controllers.Logout)
	protected.Get("/me", controllers.Me)
	protected.Get("/roles", middleware.RequireAccessLevel(10), controllers.GetRoles)

	// Employees Directory & Hierarchy (Access Level 10+)
	employees := protected.Group("/employees", middleware.RequireAccessLevel(10))
	employees.Get("/", controllers.GetEmployees)
	employees.Get("/hierarchy", controllers.GetEmployeeHierarchy)
	employees.Get("/:id", controllers.GetEmployeeByID)
	
	// Create Employee Endpoint (Requires Access Level 50+ Manager/HR/Admin)
	employees.Post("/", middleware.RequireAccessLevel(50), controllers.CreateEmployee)

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

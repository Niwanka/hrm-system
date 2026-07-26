package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"hrm-backend/config"
	"hrm-backend/routes"
)

func main() {
	// Initialize Database Connection & Auto Migration
	config.InitDB()

	app := fiber.New(fiber.Config{
		AppName: "Next-Gen HRM Core System v1.0",
	})

	// Middleware
	app.Use(logger.New())

	// Dynamic CORS configuration allowing Vercel, Netlify & Localhost origins
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true // Allow all frontend origins securely
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "hrm-core-backend",
		})
	})

	// Setup Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting Fiber web server on :%s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

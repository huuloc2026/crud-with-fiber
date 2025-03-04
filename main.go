package main

import (
	"fmt"
	"log"
	"myapp/config"
	"myapp/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Connect to database
	db := config.ConnectDB(&cfg.Database)

	// Setup routes
	routes.SetupRoutes(app, db)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package routes

import (
	"myapp/config"
	"myapp/handlers"
	"myapp/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, config *config.Config) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(
		db,
		config.JWT.Secret,
		config.JWT.ExpiresIn,
	)
	productHandler := handlers.NewProductHandler(db)

	// Public routes
	auth := app.Group("/api/auth")
	auth.Post("/login", userHandler.Login)

	// User routes
	users := app.Group("/api/users")
	users.Post("/", userHandler.CreateUser)
	users.Use(middleware.JWTProtected(config.JWT.Secret))
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Product routes
	products := app.Group("/api/products")
	products.Use(middleware.JWTProtected(config.JWT.Secret))
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)
}

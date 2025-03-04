package routes

import (
	"myapp/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)

	// User routes
	users := app.Group("/api/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Product routes
	products := app.Group("/api/products")
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)
}

package routes

import (
	"log"
	"myapp/config"
	"myapp/handlers"
	"myapp/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client, rabbitmqChannel *amqp.Channel, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(
		db,
		redisClient,
		rabbitmqChannel,
		cfg.JWT.Secret,
		24*time.Hour,
	)
	duration, err := time.ParseDuration(cfg.JWT.ExpiresIn)
	if err != nil {
		log.Fatalf("Invalid JWT expiration duration: %v", err)
	}
	userHandler := handlers.NewUserHandler(
		db,
		cfg.JWT.Secret,
		duration,
	)
	productHandler := handlers.NewProductHandler(db)

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	api := app.Group("/api", middleware.Protected(cfg.JWT.Secret))

	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Product routes
	products := api.Group("/products")
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)
}

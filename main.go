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

	// Initialize Redis
	redisCfg := &config.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	}
	redisClient := config.NewRedisClient(redisCfg)
	defer redisClient.Close()

	// Initialize RabbitMQ
	rabbitmqCfg := &config.RabbitMQConfig{
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
		User:     cfg.RabbitMQ.User,
		Password: cfg.RabbitMQ.Password,
	}
	rabbitmqConn, rabbitmqChannel := config.NewRabbitMQConnection(rabbitmqCfg)
	defer rabbitmqConn.Close()
	defer rabbitmqChannel.Close()

	// Setup routes with Redis and RabbitMQ
	routes.SetupRoutes(app, db, redisClient, rabbitmqChannel, cfg)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

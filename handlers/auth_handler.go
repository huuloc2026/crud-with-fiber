package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/models"
	"myapp/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	Redis     *redis.Client
	RabbitMQ  *amqp.Channel
	JWTSecret string
	ExpiresIn time.Duration
}

func NewAuthHandler(db *gorm.DB, redis *redis.Client, rabbitmq *amqp.Channel, jwtSecret string, expiresIn time.Duration) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Redis:     redis,
		RabbitMQ:  rabbitmq,
		JWTSecret: jwtSecret,
		ExpiresIn: expiresIn,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	user.Password = hashedPassword

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(h.ExpiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   tokenString,
	})
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	input := new(models.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check Redis cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:email:%s", input.Email)
	cachedUser, err := h.Redis.Get(ctx, cacheKey).Result()

	var user models.User
	if err == nil {
		// User found in cache
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Cache error"})
		}
	} else {
		// User not in cache, query DB
		if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Cache user data
		userData, _ := json.Marshal(user)
		h.Redis.Set(ctx, cacheKey, userData, 24*time.Hour)
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	passCheckpoint := utils.CheckPassword(hashPassword, input.Password)

	if !passCheckpoint {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(h.ExpiresIn).Unix()

	t, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not login"})
	}

	// Send login notification to RabbitMQ
	notification := map[string]interface{}{
		"event":     "user_login",
		"user_id":   user.ID,
		"timestamp": time.Now(),
	}
	notificationBytes, _ := json.Marshal(notification)

	err = h.RabbitMQ.Publish(
		"",              // exchange
		"notifications", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        notificationBytes,
		},
	)
	if err != nil {
		// Log the error but don't fail the login
		fmt.Printf("Failed to publish login notification: %v\n", err)
	}

	return c.JSON(fiber.Map{"token": t})
}

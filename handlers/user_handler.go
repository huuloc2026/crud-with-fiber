package handlers

import (
	"myapp/models"
	"myapp/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB        *gorm.DB
	JWTSecret string
	ExpiresIn time.Duration
}

func NewUserHandler(db *gorm.DB, jwtSecret string, expiresIn time.Duration) *UserHandler {
	return &UserHandler{
		DB:        db,
		JWTSecret: jwtSecret,
		ExpiresIn: expiresIn,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(h.ExpiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create token"})
	}
	return c.Status(201).JSON(fiber.Map{"user": user, "token": tokenString})
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := utils.Paginate(h.DB.Find(&users), &users, page, limit)
	return c.JSON(pagination)
}

// GetUser returns a specific user
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := h.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.DB.First(&models.User{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	h.DB.Model(&models.User{}).Where("id = ?", id).Updates(user)
	return c.JSON(user)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return c.SendStatus(204)
}

// Login authenticates a user and returns a JWT token
func (h *UserHandler) Login(c *fiber.Ctx) error {
	input := new(models.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// In production, use proper password hashing comparison
	if input.Password != user.Password {
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

	return c.JSON(fiber.Map{"token": t})
}

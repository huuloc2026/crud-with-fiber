package handlers

import (
	"myapp/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
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

	return c.Status(201).JSON(user)
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User
	h.DB.Find(&users)
	return c.JSON(users)
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

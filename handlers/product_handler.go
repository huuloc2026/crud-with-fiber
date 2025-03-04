package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myapp/models"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create product"})
	}

	return c.Status(201).JSON(product)
}

// GetProducts returns all products
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	h.DB.Find(&products)
	return c.JSON(products)
}

// GetProduct returns a specific product
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := h.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(product)
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.DB.First(&models.Product{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	h.DB.Model(&models.Product{}).Where("id = ?", id).Updates(product)
	return c.JSON(product)
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.DB.Delete(&models.Product{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete product"})
	}

	return c.SendStatus(204)
}

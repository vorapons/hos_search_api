package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HelloHandler handles GET /hello
func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

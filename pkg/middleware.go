package pkg

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		userID, err := ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
func RecoverFromError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)

				if err := c.Status(500).JSON(fiber.Map{"error": "Internal server error"}); err != nil {
					log.Printf("Error sending response: %v", err)
				}
			}
		}()
		return c.Next()
	}
}

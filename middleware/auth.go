package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

type Claims struct {
	Login      string `json:"login"`
	HospitalID string `json:"hospital_id"`
	Hospital   string `json:"hospital"`
	jwt.RegisteredClaims
}

// JWTProtected validates the Bearer token and injects claims into c.Locals.
// isBlacklisted is called to reject logged-out tokens.
func JWTProtected(jwtSecret string, isBlacklisted func(string) bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid Authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if isBlacklisted(tokenString) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Token has been invalidated",
			})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Invalid or expired token",
			})
		}

		c.Locals("login", claims.Login)
		c.Locals("hospital_id", claims.HospitalID)
		c.Locals("hospital", claims.Hospital)
		c.Locals("token", tokenString)
		c.Locals("exp", claims.ExpiresAt.Time)
		return c.Next()
	}
}

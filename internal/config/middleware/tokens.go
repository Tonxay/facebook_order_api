package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("SECRET_GENTOKEN"))

func GenerateToken(userID string, route string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    route,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

func JWTProtected(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Expecting: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization format"})
	}
	tokenStr := parts[1]

	// Validate token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Extract user_id and save in context
	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", claims["user_id"])

	return c.Next()

}

func GetJwtSecret() []byte {
	return jwtSecret
}

func GenerateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%d", now.Format("20060102"), now.UnixNano()%100000)
}

func GetUserID(c *fiber.Ctx) (string, bool) {
	user := c.Locals("user_id")
	userID, ok := user.(string)
	return userID, ok
}

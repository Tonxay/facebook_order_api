package middleware

import (
	"fmt"
	"math/rand"
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

var generated = make(map[string]bool) // Track generated order numbers

func GenerateOrderNumber() string {
	now := time.Now()
	prefix := now.Format("0601") // "yyMM", e.g., "2406"

	for {
		suffix := rand.Intn(10000) // 0000 to 9999
		orderNo := fmt.Sprintf("%s%04d", prefix, suffix)

		if !generated[orderNo] {
			generated[orderNo] = true
			return orderNo
		}
		// Retry until a unique number is found
	}
}
func GetUserID(c *fiber.Ctx) (string, bool) {
	user := c.Locals("user_id")
	userID, ok := user.(string)
	return userID, ok
}

func GenerateFacebookID() string {
	rand.Seed(time.Now().UnixNano())
	// First digit must not be 0
	id := fmt.Sprintf("%d", rand.Intn(9)+1)
	// Add remaining 16 digits
	for i := 0; i < 16; i++ {
		id += fmt.Sprintf("%d", rand.Intn(10))
	}

	return id
}

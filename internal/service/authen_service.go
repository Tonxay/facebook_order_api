package service

import (
	"go-api/internal/config/middleware"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/service/db_service"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var requestData request.User

	if err := c.BodyParser(&requestData); err != nil {
		return fiber.NewError(400, "Invalid input")

	}

	db := gormpkg.GetDB()
	userAuthen, _ := dbservice.GetUserForUserName(db, requestData.UserName)
	if userAuthen.UserName != "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}

	pageID := os.Getenv("PAGE_ID")

	user := models.User{
		UserName: requestData.UserName,
		Password: requestData.Password,
		Status:   "active",
		Rolo:     "admin",
		PageID:   pageID,
	}

	err := dbservice.CreateUser(db, &user, c.Context())

	if err != nil {
		return fiber.NewError(500, "Server error")
	}

	return c.JSON(fiber.Map{"message": "User registered", "user_id": user.ID})
}

func Login(c *fiber.Ctx) error {
	var requestData request.User

	if err := middleware.ParseAndValidateBody(c, &requestData); err != nil {
		return fiber.NewError(400, err.Error())
	}

	user, err := dbservice.GetUserNamePassword(gormpkg.GetDB(), requestData.UserName)
	if err != nil {
		return fiber.NewError(401, "Invalid credentials")
	}

	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))

	if err1 != nil {
		return fiber.NewError(401, "Invalid credentials")
	}

	accessToken, _ := middleware.GenerateToken(user.ID, user.Rolo, time.Hour+7)
	refreshToken, _ := middleware.GenerateToken(user.ID, user.Rolo, time.Hour+8)

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh token handler
func Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Parse and validate the token
	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return middleware.GetJwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	role := claims["role"].(string)

	// Issue new access token
	newAccessToken, _ := middleware.GenerateToken(userID, role, time.Minute*15)
	refreshToken, _ := middleware.GenerateToken(userID, role, time.Hour*24)
	return c.JSON(fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": refreshToken,
	})
}

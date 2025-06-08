package middlewares

import (
	"DeallsJobsTest/config"
	"DeallsJobsTest/models"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JWTMiddlewareHandler struct {
	DB *gorm.DB
}

func NewJWTMiddlewareHandler(db *gorm.DB) *JWTMiddlewareHandler {
	return &JWTMiddlewareHandler{
		DB: db,
	}
}

func (h *JWTMiddlewareHandler) JWTMiddleware(c *fiber.Ctx) error {
	var envCfg config.EnvConfig
	if errcnf := env.Parse(&envCfg); errcnf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse environment configuration",
			"error":   errcnf.Error(),
		})
	}
	jwtSecret := []byte(envCfg.SecretKEY)
	authHeader := c.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	tokenStr := authHeader[7:]
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	claims := &models.JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("userID", claims.ID)
	c.Locals("role", claims.Role)
	c.Locals("username", claims.Username)

	return c.Next()
}

func (h *JWTMiddlewareHandler) ValidateUserRole(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Role not found"})
	}

	if role != "admin" && role != "employee" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if role == "admin" {
		c.Locals("isAdmin", true)
	} else {
		c.Locals("isAdmin", false)
	}

	return c.Next()
}

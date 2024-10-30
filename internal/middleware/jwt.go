package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("very-secret-key") // Убедитесь, что этот ключ хранится в безопасном месте

// JWTMiddleware проверяет наличие и валидность JWT токена
func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Отсутствует токен",
		})
	}

	token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный токен",
		})
	}

	// Сохраняем токен в контексте, если необходимо
	c.Locals("user", token.Claims)
	return c.Next() // Продолжаем выполнение следующего обработчика
}

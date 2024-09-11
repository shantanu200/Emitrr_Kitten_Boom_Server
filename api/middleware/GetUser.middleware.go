package middleware

import (

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserId(c *fiber.Ctx) (string, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	_id := claims["username"]

	if _id == nil {
		return "", fiber.ErrUnauthorized
	}

	return _id.(string), nil
}

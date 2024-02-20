package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"hotelapi.com/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"Type": "error",
			"Msg":  "unauthorized",
		})
	}
	if !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"Type": "error",
			"Msg":  "unauthorized",
		})
	}
	return c.Next()
}

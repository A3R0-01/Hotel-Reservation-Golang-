package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelapi.com/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedNormal(c)
	}
	if !user.IsAdmin {
		return UnauthorizedNormal(c)
	}
	return c.Next()
}

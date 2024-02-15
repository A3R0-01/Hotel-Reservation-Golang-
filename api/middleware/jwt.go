package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorized user")
	}
	fmt.Println("token:", token)
	return nil
}

func parseJWTToken(tokenStr string) error {
	token, err := jwt.Pa
}

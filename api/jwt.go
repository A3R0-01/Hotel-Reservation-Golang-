package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelapi.com/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{
				"Type": "unauthorized",
				"Msg":  "please login ",
			})
		}
		claim, err := validateJWTToken(transform(token))
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{
				"Type": "error",
				"Msg":  "invalid user or token",
			})
		}
		// check token expiration
		expiresFloat := claim["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > (expires) {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{
				"Type": "unauthorized",
				"Msg":  "please login again",
			})
		}
		userID := claim["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{
				"Type": "unauthorized",
				"Msg":  "you are not registered",
			})
		}
		//  set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}

}

func validateJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, parseFunc)
	if err != nil {
		return nil, fmt.Errorf("unauthorized access")
	}
	if !token.Valid {
		fmt.Println("invalid token: ")
		return nil, fmt.Errorf("unauthorized access")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
func parseFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		fmt.Println("Invalid signing method", t.Header["alg"])
		return nil, fmt.Errorf("Unauthorised access")
	}

	secret := os.Getenv("JWT_SECRET")
	return []byte(secret), nil
}
func transform(array []string) string {
	var storage string = ""
	for _, str := range array {
		storage += str
	}
	return storage
}

// func JWTAuthentication(c *fiber.Ctx) error {
// 	fmt.Println("--")

// 	token, ok := c.GetReqHeaders()["X-Api-Token"]
// 	if !ok {
// 		fmt.Println("token not present in the header")
// 		return fmt.Errorf("unauthorized user")
// 	}
// 	claim, err := validateJWTToken(transform(token))
// 	if err != nil {
// 		return err
// 	}
// 	// check token expiration
// 	expiresFloat := claim["expires"].(float64)
// 	expires := int64(expiresFloat)
// 	if time.Now().Unix() > (expires) {
// 		return fmt.Errorf("toke expired")
// 	}
// 	userID := claim["id"]
// 	return c.Next()
// }

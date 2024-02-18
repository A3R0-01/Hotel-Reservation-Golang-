package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return invalidCredentials(c)
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		return invalidCredentials(c)
	}
	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return invalidCredentials(c)
	}

	fmt.Println("authenticated user: ", user)
	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}
func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}
func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * types.HoursAuth).Unix()
	claims := jwt.MapClaims{
		"email":   user.Email,
		"id":      user.ID,
		"expires": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(secret)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr

}

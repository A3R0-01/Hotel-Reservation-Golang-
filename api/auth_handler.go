package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
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
		return BadRequest(c)
	}
	user, err := h.store.User.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		return InvalidCredentials(c)
	}
	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return InvalidCredentials(c)
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * types.HoursAuth).Unix()
	claims := jwt.MapClaims{
		"email":   user.Email,
		"id":      user.ID,
		"expires": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr

}

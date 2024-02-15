package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"msg": "not found"})
		}
		return err
	}
	return c.JSON(user)
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	createdUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(createdUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	err := h.userStore.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"msg": "user deleted",
	})
}
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return c.JSON(map[string]string{"msg": "invalid fields"})
	}
	userId := c.Params("id")
	if err := h.userStore.UpdateUser(c.Context(), userId, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "user updated"})
}

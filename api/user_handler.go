package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}
func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedNormal(c)
	}
	if (user.ID.Hex() != id) && !user.IsAdmin {
		return UnauthorizedNormal(c)
	}
	user, err := h.store.User.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(c, "user")
		}
		return InternalServerError(c, "could not fetch user please try again")
	}
	return c.JSON(user)
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return InternalServerError(c, "could not fetch users please try again")
	}
	return c.JSON(users)
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return InternalServerError(c, "failed to encrypt password")
	}
	createdUser, err := h.store.User.CreateUser(c.Context(), user)
	if err != nil {
		return InternalServerError(c, "failed to create user")
	}
	return c.JSON(createdUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	userId := c.Params("id")
	if !ok {
		return UnauthorizedNormal(c)
	}
	if (user.ID.Hex() != userId) && !user.IsAdmin {
		return UnauthorizedNormal(c)
	}
	err := h.store.User.DeleteUser(c.Context(), userId)
	if err != nil {
		return InternalServerError(c, "failed to delete user")

	}
	return c.JSON(map[string]string{
		"msg": "user deleted",
	})
}
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return BadRequest(c)
	}
	userId := c.Params("id")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedNormal(c)
	}
	if (user.ID.Hex() != userId) && !user.IsAdmin {
		return UnauthorizedNormal(c)
	}
	if err := h.store.User.UpdateUser(c.Context(), userId, params); err != nil {
		return InternalServerError(c, "failed to update")
	}
	return c.JSON(map[string]string{"msg": "user updated"})
}

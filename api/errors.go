package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func ErrRespond(c *fiber.Ctx, status int, typeMsg, msg string) error {
	return c.Status(status).JSON(genericResp{
		Type: typeMsg,
		Msg:  msg,
	})
}

func BadRequest(c *fiber.Ctx) error {
	return ErrRespond(c, http.StatusBadRequest, "error", "invalid parameters")
}
func InvalidID(c *fiber.Ctx, name string) error {
	return ErrRespond(c, http.StatusBadRequest, "error", fmt.Sprint("invalid ", name))
}
func InternalServerError(c *fiber.Ctx, msg string) error {
	return ErrRespond(c, http.StatusInternalServerError, "error", "invalid parameters")
}
func UnauthorizedNormal(c *fiber.Ctx) error {
	return ErrRespond(c, http.StatusUnauthorized, "error", "you are unauthorized to access this info")
}
func UnauthorizedSpec(c *fiber.Ctx, msg string) error {
	return ErrRespond(c, http.StatusUnauthorized, "error", msg)
}
func NotFound(c *fiber.Ctx, name string) error {
	return ErrRespond(c, http.StatusNotFound, "error", fmt.Sprint(name, " not found"))
}
func Conflict(c *fiber.Ctx, msg string) error {
	return ErrRespond(c, http.StatusConflict, "error", msg)
}
func InvalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	errNew := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(errNew.Code).JSON(errNew)
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}
func ErrInvalidID(name string) Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  fmt.Sprint("invalid ", name),
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Msg:  "you are unauthorized to access this info",
	}
}
func ErrNotFound(name string) Error {
	return Error{
		Code: http.StatusNotFound,
		Msg:  fmt.Sprint(name, " not found"),
	}
}

package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "invalid booking",
		})
	}
	booking, err := h.store.Booking.GetBookingByID(c.Context(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internal Server Error",
		})
	}
	if user.ID != booking.UserID && !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "you are unauthorized to access this info",
		})
	}
	return c.JSON(booking)
}

// needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internal Server Error",
		})
	}
	var bookings []*types.Booking
	if user.IsAdmin {
		booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
		if err != nil {
			return err
		}
		bookings = booking
	} else {
		booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{"userID": user.ID})
		if err != nil {
			return err
		}
		bookings = booking
	}

	return c.JSON(bookings)
}
func (h *BookingHandler) HandleUpdateBooking(c *fiber.Ctx) error {
	return nil
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internal Server Error",
		})
	}
	if err := h.store.Booking.CancelBooking(c.Context(), id, user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  err.Error(),
		})
	}
	return c.JSON(genericResp{
		Type: "success",
		Msg:  "booking cancelled",
	})
}

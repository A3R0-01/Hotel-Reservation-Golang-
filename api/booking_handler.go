package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return BadRequest(c)
	}
	booking, err := h.store.Booking.GetBookingByID(c.Context(), db.Map{"_id": oid})
	if err != nil {
		return NotFound(c, "booking")
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return InternalServerError(c, "could not identify user")
	}
	if user.ID != booking.UserID && !user.IsAdmin {
		return UnauthorizedNormal(c)
	}
	return c.JSON(booking)
}

// needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedNormal(c)
	}
	var bookings []*types.Booking
	if user.IsAdmin {
		booking, err := h.store.Booking.GetBookings(c.Context(), db.Map{})
		if err != nil {
			return NotFound(c, "bookings")
		}
		bookings = booking
	} else {
		booking, err := h.store.Booking.GetBookings(c.Context(), db.Map{"userID": user.ID})
		if err != nil {
			return NotFound(c, "bookings")
		}
		bookings = booking
	}

	return c.JSON(bookings)
}
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedSpec(c, "please sign in again")
	}
	if err := h.store.Booking.CancelBooking(c.Context(), id, user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(c, "booking")
		}
		return InternalServerError(c, "failed to cancel booking")
	}
	return c.JSON(genericResp{
		Type: "success",
		Msg:  "booking cancelled",
	})
}

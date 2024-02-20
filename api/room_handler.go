package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type RoomHandler struct {
	store *db.Store
}
type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p *BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	if p.FromDate.After(p.TillDate) {
		return fmt.Errorf("cannot set expiration date after the start date")
	}
	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	room, err := h.store.Room.GetRoomByIDOne(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(room)
}
func (h *RoomHandler) HandleGetBookingsPerRoom(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "invalid room",
		})
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{
		"roomID": oid,
		"fromDate": bson.M{
			"$gte": time.Now(),
		},
		"cancelled": false,
	})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}
func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	var params types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return fmt.Errorf("failed to create room")
	}
	room, err := params.CreateRoom()
	if err != nil {
		return fmt.Errorf("failed to create room")
	}
	createdRoom, err := h.store.Room.InsertRoom(c.Context(), &room)
	if err != nil {
		return fmt.Errorf("failed to create room")
	}
	return c.JSON(createdRoom)

}
func (h *RoomHandler) HandleDeleteRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.store.Room.DeleteRoom(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNoContent).JSON(genericResp{
			Type: "error",
			Msg:  err.Error(),
		})
	}
	return c.JSON(genericResp{
		Type: "success",
		Msg:  "room deleted",
	})
}
func (store *RoomHandler) HandlePutRoom(c *fiber.Ctx) error {
	return nil
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))

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
	roomAvail, err := h.isRoomAvailable(c.Context(), roomOID, params)
	if err != nil {
		return err
	}
	if !roomAvail {
		return c.Status(http.StatusConflict).JSON(genericResp{
			Type: "error",
			Msg:  "room is occupied",
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)

}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, id primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": id,
		"tillDate": bson.M{
			"$gte": params.FromDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil

}

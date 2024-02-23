package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return NotFound(c, "room")
	}
	return c.JSON(room)
}
func (h *RoomHandler) HandleGetBookingsPerRoom(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return InvalidID(c, "booking")
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), db.Map{
		"roomID": oid,
		"fromDate": db.Map{
			"$gte": time.Now(),
		},
		"cancelled": false,
	})
	if err != nil {
		return InternalServerError(c, "could not get bookings")
	}
	return c.JSON(bookings)
}
func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	var params types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return BadRequest(c)

	}
	room, err := params.CreateRoom()
	if err != nil {
		return InternalServerError(c, "failed to create room")
	}
	createdRoom, err := h.store.Room.InsertRoom(c.Context(), &room, h.store)
	if err != nil {
		return InternalServerError(c, "failed to create room")
	}
	return c.JSON(createdRoom)

}
func (h *RoomHandler) HandleDeleteRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.store.Room.DeleteRoom(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(c, "room")
		}
		return InternalServerError(c, "failed to delete room")
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
		return BadRequest(c)
	}
	if err := params.validate(); err != nil {
		return BadRequest(c)
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return InvalidID(c, "room")
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return UnauthorizedNormal(c)
	}
	roomAvail, err := h.isRoomAvailable(c.Context(), roomOID, params)
	if err != nil {
		return Conflict(c, "room is booked")
	}
	if !roomAvail {
		return Conflict(c, "room is booked")
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
		return InternalServerError(c, "failed to book room")
	}
	return c.JSON(inserted)

}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.Map{})
	if err != nil {
		return NotFound(c, "room")
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, id primitive.ObjectID, params BookRoomParams) (bool, error) {
	// unreliable filter
	// filter := bson.M{
	// 	"roomID": id,
	// 	"fromDate": bson.M{
	// 		"$gte": params.FromDate,
	// 	},
	// 	"tillDate": bson.M{
	// 		"$lte": params.TillDate,
	// 	},
	// }
	filter2 := db.Map{
		"roomID": id,
		"fromDate": db.Map{
			"$lte": params.FromDate,
		},
		"tillDate": db.Map{
			"$gte": params.FromDate,
		},
	}
	filter3 := db.Map{
		"roomID": id,
		"fromDate": db.Map{
			"$gte": params.FromDate,
		},
		"FromDate": db.Map{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter3)
	if err != nil {
		return false, err
	}
	bookings2, err := h.store.Booking.GetBookings(ctx, filter2)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	ok2 := len(bookings2) == 0
	return ok && ok2, nil

}

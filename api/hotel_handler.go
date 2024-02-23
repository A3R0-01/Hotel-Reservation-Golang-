package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params types.CreateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return BadRequest(c)
	}
	if err := params.Validate(); len(err) > 0 {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	hotel := params.CreateNewHotelFromParams()
	insertedHotel, err := h.store.Hotel.InsertHotel(c.Context(), hotel)
	if err != nil {
		return InternalServerError(c, "failed to create hotel")
	}
	return c.JSON(insertedHotel)
}
func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return NotFound(c, "an error occurred, hotels")
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return InvalidID(c, "hotel")
	}
	filter := db.Map{"_id": oid}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), filter)
	if err != nil {
		return NotFound(c, "hotel")
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return InvalidID(c, "hotel")
	}
	filter := db.Map{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return NotFound(c, "rooms")
	}
	return c.JSON(rooms)
}

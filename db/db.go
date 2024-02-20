package db

import (
	"github.com/gofiber/fiber/v2"
)

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(map[string]string{"error": err.Error()})
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

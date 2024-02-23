package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

func AddBooking(store *db.Store, roomID, userID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal("failed to book a room\n", err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, price float64, size string, seaSideRoom bool, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Price:       price,
		Size:        size,
		SeaSideRoom: seaSideRoom,
		HotelID:     hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room, store)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}
func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}
func AddUser(store *db.Store, isAdmin bool, fname, lname string, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@gmail.com", fname),
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

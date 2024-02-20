package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedBooking(userID, roomID primitive.ObjectID, fromDate, tillDate time.Time, numPersons int) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: fromDate,
		TillDate: tillDate,
	}
	insertedBooking, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal("failed to book a room\n", err)
	}
	fmt.Println(insertedBooking.ID, insertedBooking.Cancelled)
}
func seedUser(isAdmin bool, fname, lname, email, password string) primitive.ObjectID {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser.ID

}
func seedHotel(name, location string, rating int) []*types.Room {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:        "small",
			SeaSideRoom: true,
			Price:       99.9,
		},
		{
			Size:        "kingsize",
			SeaSideRoom: true,
			Price:       199.9,
		},
		{
			Size:        "normal",
			SeaSideRoom: false,
			Price:       19.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Roomms for hotel: ", insertedHotel.Name, "\t", insertedHotel.ID)
	insertedRooms := []*types.Room{}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		insertedRooms = append(insertedRooms, insertedRoom)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
	return insertedRooms
}
func main() {
	userID := seedUser(false, "Esther", "Mutane", "esther@gmail.com", types.DefaultUserPassword)
	userID2 := seedUser(false, "Tilly", "Monkey", "tilly@gmail.com", types.DefaultUserPassword)
	_ = seedUser(true, "Harvey", "Specter", "harvey@gmail.com", "1234bsrvnt")
	hotel1 := seedHotel("Bellucia", "France", 5)
	_ = seedHotel("The Cozy Hotel", "Netherlands", 3)
	_ = seedHotel("Meikles Hotel", "Zimbabwe", 4)
	for key, id := range []primitive.ObjectID{userID, userID2} {
		for key2, room := range hotel1 {
			number2 := (key2 + 1) * (key + 1) * 24 * 2
			from := time.Now().Add(time.Hour * time.Duration(number2))
			till := from.Add(time.Hour * (24 * 2))
			seedBooking(id, room.ID, from, till, 2)
		}
	}
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

}

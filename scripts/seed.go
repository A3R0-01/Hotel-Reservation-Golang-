package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/api"
	"hotelapi.com/db"
	"hotelapi.com/db/fixtures"
	"hotelapi.com/types"
)

var (
	client *mongo.Client
	store  *db.Store
	ctx    = context.Background()
)

func main() {
	admin := fixtures.AddUser(store, true, "harvey", "Specter", "1234bsrvnt")
	user := fixtures.AddUser(store, false, "tilly", "Monkey", types.DefaultUserPassword)
	hotel := fixtures.AddHotel(store, "Njeke Hotel", "Harare, Zimbabwe", 4, []primitive.ObjectID{})
	room := fixtures.AddRoom(store, 88.50, "large", true, hotel.ID)
	booking := fixtures.AddBooking(store, room.ID, user.ID, time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 5))
	fmt.Println("Admin:->\t", admin, "\n Token: ->\t", api.CreateTokenFromUser(admin))
	fmt.Println("User:->\t", user, "\n Token: ->\t", api.CreateTokenFromUser(user))

	for i := 0; i < 100; i++ {
		name := fmt.Sprint("random hotel ", i)
		location := fmt.Sprint("location ", i)
		_ = fixtures.AddHotel(store, name, location, rand.Intn(5)+1, []primitive.ObjectID{})
	}
	fmt.Println("Hotel:->\t", hotel)
	fmt.Println("Room:->\t", room)
	fmt.Println("Booking:->\t", booking)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store = &db.Store{
		Room:    db.NewMongoRoomStore(client),
		User:    db.NewMongoUserStore(client),
		Hotel:   db.NewMongoHotelStore(client),
		Booking: db.NewMongoBookingStore(client),
	}

}

package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname, lname, email, password string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

}
func seedHotel(name, location string, rating int) {
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
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(insertedRoom)
	}
}
func main() {
	seedUser("James", "Kimmy", "james@gmail.com", "hellomanu")
	seedUser("Toby", "Money", "toby@gmail.com", "hellomanu")
	seedUser("Takudzwa", "Jaja", "takudzwa@gmail.com", "hellomanu")
	seedHotel("Bellucia", "France", 5)
	seedHotel("The Cozy Hotel", "Netherlands", 3)
	seedHotel("Meikles Hotel", "Zimbabwe", 4)

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

}

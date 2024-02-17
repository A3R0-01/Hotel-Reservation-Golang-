package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/api"
	"hotelapi.com/api/middleware"
	"hotelapi.com/db"
)

func main() {
	listenAddress := flag.String("listenAddress", ":5000", "The listen address for the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	config := fiber.Config{ErrorHandler: db.ErrorHandler}

	// Declarations
	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		// apiV1        = app.Group("/api/v1")
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUserById)

	// Hotel Handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// authentication

	auth.Post("/auth", authHandler.HandleAuthenticate)

	app.Listen(*listenAddress)

}

func HandleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Server is running"})
}

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
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		// bookingHandler = api.NewBookingHandler(store)
		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", middleware.AdminAuth)
	)

	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUserById)

	// Hotel Handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Post("/hotel", hotelHandler.HandlePostHotel)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	// Booking Handlers
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/room/:id/book", roomHandler.HandleGetBookingsPerRoom)
	apiV1.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// cancel booking

	// Room Handlers
	apiV1.Post("/room", roomHandler.HandlePostRoom)
	apiV1.Get("/room/:id", roomHandler.HandleGetRoom)
	apiV1.Delete("/room/:id", roomHandler.HandleDeleteRoom)
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	// authentication
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// admin

	admin.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddress)

}

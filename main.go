package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/api"
	"hotelapi.com/db"
)

func main() {
	// mongoEndPoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	config := fiber.Config{ErrorHandler: api.ErrorHandler}

	// Declarations
	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(store)
		authHandler    = api.NewAuthHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		// bookingHandler = api.NewBookingHandler(store)
		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", api.AdminAuth)
	)

	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUserById)

	// Hotel Handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
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
	apiV1.Get("/room/:id", roomHandler.HandleGetRoom)
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	// authentication
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// admin
	admin.Get("/user", userHandler.HandleGetUsers)
	admin.Post("/hotel", hotelHandler.HandlePostHotel)
	admin.Post("/room", roomHandler.HandlePostRoom)
	admin.Delete("/room/:id", roomHandler.HandleDeleteRoom)
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	listenAddress := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddress)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

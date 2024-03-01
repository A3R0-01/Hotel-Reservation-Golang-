package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DBNAME string
var DBURI string
var TestDBNAME string

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Map map[string]any

type Pagination struct {
	Limit int64
	Page  int64
}

func init() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal(err)
		}
	}
	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURI = os.Getenv("MONGO_DB_URL")
	TestDBNAME = os.Getenv("TEST_DB_NAME")
}

package api

import (
	"context"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/db"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		client: client,
		store: &db.Store{
			User:    db.NewMongoUserStoreTest(client),
			Room:    db.NewMongoRoomStoreTest(client),
			Hotel:   db.NewMongoHotelStore(client),
			Booking: db.NewMongoBookingStoreTest(client),
		},
	}
}

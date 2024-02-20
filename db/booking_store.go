package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelapi.com/types"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookingByID(context.Context, bson.M) (*types.Booking, error)
	UpdateBooking(context.Context, bson.M, bson.M) error
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	CancelBooking(ctx context.Context, id string, user *types.User) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(bookingColl),
	}
}

func (store *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {

	resp, err := store.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (store *MongoBookingStore) GetBookingByID(ctx context.Context, filter bson.M) (*types.Booking, error) {
	var booking types.Booking
	err := store.coll.FindOne(ctx, filter).Decode(&booking)
	return &booking, err
}

func (store *MongoBookingStore) UpdateBooking(ctx context.Context, filter bson.M, values bson.M) error {
	_, err := store.coll.UpdateOne(ctx, filter, values)
	return err
}

func (store *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking
	resp, err := store.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := resp.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (store *MongoBookingStore) CancelBooking(ctx context.Context, id string, user *types.User) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid booking")
	}
	booking, err := store.GetBookingByID(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("invalid booking")
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return fmt.Errorf("unauthorized")
	}
	update := bson.M{"$set": bson.M{"cancelled": true}}
	_, err = store.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return fmt.Errorf("failed to update")
	}
	return nil
}

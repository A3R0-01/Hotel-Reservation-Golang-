package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelapi.com/types"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	GetHotelByID(context.Context, Map) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}
func NewMongoHotelStoreTest(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(TestDBNAME).Collection(hotelColl),
	}
}

func (store *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := store.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (store *MongoHotelStore) Update(ctx context.Context, filter Map, values Map) error {
	_, err := store.coll.UpdateOne(ctx, filter, values)
	return err
}
func (store *MongoHotelStore) GetHotelByID(ctx context.Context, filter Map) (*types.Hotel, error) {
	var hotel types.Hotel
	err := store.coll.FindOne(ctx, filter).Decode(&hotel)
	return &hotel, err
}
func (store *MongoHotelStore) GetHotels(ctx context.Context, filter Map) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	resp, err := store.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

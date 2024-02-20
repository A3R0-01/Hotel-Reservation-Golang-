package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelapi.com/types"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error)
	GetRoomByIDOne(ctx context.Context, id string) (*types.Room, error)
	DeleteRoom(ctx context.Context, id string) error
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (store *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	// check if the hotelID is valid and if it exists
	_, err := store.HotelStore.GetHotelByID(ctx, bson.M{"_id": room.HotelID})
	if err != nil {
		return nil, fmt.Errorf("invalid hotel")
	}
	//
	resp, err := store.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := store.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (store *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {

	resp, err := store.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
func (store *MongoRoomStore) GetRoomByIDOne(ctx context.Context, id string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	var room types.Room
	if err := store.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		return nil, fmt.Errorf("failed to find room")
	}
	return &room, nil
}

func (store *MongoRoomStore) DeleteRoom(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	res, err := store.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("could not find room")
	}
	return nil
}

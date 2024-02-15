package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type RoomType int

const (
	_ RoomType = iota
	SinglePersonRoom
	DoubleRoom
	SeaSideRoom
	DeluxeRoom
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price       float64            `bson:"price" json:"price"`
	Size        string             `bson:"size" json:"size"`
	SeaSideRoom bool               `bson:"SeaSideRoom" json:"SeaSideRoom"`
	HotelID     primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

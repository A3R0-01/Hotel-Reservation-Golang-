package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price       float64            `bson:"price" json:"price"`
	Size        string             `bson:"size" json:"size"`
	SeaSideRoom bool               `bson:"SeaSideRoom" json:"SeaSideRoom"`
	HotelID     primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
type CreateRoomParams struct {
	Price       float64 `bson:"price" json:"price"`
	Size        string  `bson:"size" json:"size"`
	SeaSideRoom bool    `bson:"SeaSideRoom" json:"SeaSideRoom"`
	HotelID     string  `bson:"hotelID" json:"hotelID"`
}

func (c *CreateRoomParams) CreateRoom() (Room, error) {
	oid, err := primitive.ObjectIDFromHex(c.HotelID)

	if err != nil {
		return Room{}, err
	}

	return Room{
		Price:       c.Price,
		Size:        c.Size,
		SeaSideRoom: c.SeaSideRoom,
		HotelID:     oid,
	}, nil
}

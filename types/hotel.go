package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minHotelLen    = 5
	minHotelRating = 0
	minLocationLen = 7
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type CreateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

func (params *CreateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Name) < minHotelLen {
		errors["Name"] = fmt.Sprintf("name should be at least %d characters", minHotelLen)
	}
	if len(params.Location) < minLocationLen {
		errors["Location"] = fmt.Sprintf("location should be at least %d characters", minLocationLen)
	}
	if params.Rating < minHotelRating {
		errors["Rating"] = fmt.Sprintf("rating should be at least %d", minHotelRating)
	}
	return errors
}
func (params *CreateHotelParams) CreateNewHotelFromParams() *Hotel {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
		Rooms:    []primitive.ObjectID{},
	}
}

type UpdateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

func (params *UpdateHotelParams) Validate() bson.M {
	m := bson.M{}
	if len(params.Name) > minHotelLen {
		m["name"] = params.Name
	}
	if len(params.Location) < minLocationLen {
		m["location"] = params.Location
	}
	if params.Rating < minHotelRating {
		m["rating"] = params.Rating
	}
	return m
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID             primitive.ObjectID   `bson:"_id"`
	Name           string   `bson:"name"`
	Description    string   `bson:"description"`
	Photos         []string `bson:"photos"`
	Amenities      []string `bson:"amenities"`
	RoomCount      int      `bson:"room_count"`
	City           string   `bson:"city"`
	AvailableRooms int      `bson:"available_rooms"`
}

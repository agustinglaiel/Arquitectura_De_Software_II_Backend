package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	Id      primitive.ObjectID `bson:"_id"`
	Imagen  []byte             `bson:"imagen;type:binary"`
	HotelID primitive.ObjectID `bson:"hotel_id"`
	Hotel   *Hotel             `bson:"hotel,inline"` // Si deseas embeber los detalles del hotel
}

type Images []Image
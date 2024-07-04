package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type ImageDto struct {
	HotelId primitive.ObjectID `json:"hotel_id"`
	Data    []byte
}

type ImagesDto struct {
	Images []ImageDto `json:"images"`
}
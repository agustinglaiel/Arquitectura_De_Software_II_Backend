package daos

import (
	"context"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/db"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetHotelById(id string) models.Hotel {
	var hotel models.Hotel
	db := db.MongoDb
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	err = db.Collection("Hotels").FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(&hotel)

	if err != nil{
		fmt.Println(err)
		return hotel
	}

	return hotel
}

func InsertHotel(hotel models.Hotel) models.Hotel {
	db := db.MongoDb
	insertHotel := hotel 
	insertHotel.ID = primitive.NewObjectID()
	_, err := db.Collection("Hotels").InsertOne(context.TODO(), & insertHotel)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	hotel.ID = insertHotel.ID
	return hotel
}

func UpdateHotel(hotel models.Hotel) models.Hotel {
	db := db.MongoDb
	filter := bson.M{"_id": hotel.ID}
	update := bson.M{
		"$set": bson.M{
			"name":            hotel.Name,
			"description":     hotel.Description,
			"photos":          hotel.Photos,
			"amenities":       hotel.Amenities,
			"room_count":      hotel.RoomCount,
			"city":            hotel.City,
			"available_rooms": hotel.AvailableRooms,
		},
	}

	_, err := db.Collection("Hotels").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return hotel
	}

	return hotel
}



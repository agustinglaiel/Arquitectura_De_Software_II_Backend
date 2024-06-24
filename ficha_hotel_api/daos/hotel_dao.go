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

	err = db.Collection("hotels").FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(&hotel)

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
	_, err := db.Collection("hotels").InsertOne(context.TODO(), & insertHotel)

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

	_, err := db.Collection("hotels").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return hotel
	}

	return hotel
}

func GetHotels() ([]models.Hotel, error) {
	db := db.MongoDb
	var hotels []models.Hotel

	cursor, err := db.Collection("hotels").Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var hotel models.Hotel
		if err := cursor.Decode(&hotel); err != nil {
			fmt.Println(err)
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}

func DeleteHotelById(id primitive.ObjectID) error {
	db := db.MongoDb
	_, err := db.Collection("hotels").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
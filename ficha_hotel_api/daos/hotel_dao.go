package daos

import (
	"context"
	"ficha_hotel_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelDAO struct {
	collection *mongo.Collection
}

func NewHotelDAO(db *mongo.Database) *HotelDAO {
	return &HotelDAO{
		collection: db.Collection("hotels"),
	}
}

func (dao *HotelDAO) InsertHotel(ctx context.Context, hotel models.Hotel) error {
	hotel.ID = primitive.NewObjectID().Hex()
	_, err := dao.collection.InsertOne(ctx, hotel)
	return err
}

func (dao *HotelDAO) DeleteHotel(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (dao *HotelDAO) UpdateHotel(ctx context.Context, id primitive.ObjectID, hotel models.Hotel) error {
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
	_, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (dao *HotelDAO) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*models.Hotel, error) {
	var hotel models.Hotel
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

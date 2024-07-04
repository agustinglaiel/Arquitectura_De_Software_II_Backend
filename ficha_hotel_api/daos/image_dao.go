package daos

import (
	"context"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/db"
	"ficha_hotel_api/utils/errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertImage(image models.Image) (models.Image, errors.ApiError) {
	db := db.MongoDb
	insertImage := image 
	insertImage.Id = primitive.NewObjectID()
	_, err := db.Collection("images").InsertOne(context.TODO(), insertImage)

	if err != nil {
		fmt.Println(err)
		return image, errors.NewBadRequestApiError("Error la insertar la imágen")
	}

	return insertImage, nil
}

// GetImagesByHotelId recupera todas las imágenes asociadas con un hotel específico.
func GetImagesByHotelId(hotelId string) ([]models.Image, errors.ApiError) {
	db := db.MongoDb
	var images []models.Image

	hotelObjectId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return nil, errors.NewBadRequestApiError("Invalid hotel ID format: " + err.Error())
	}

	// Construye el filtro de búsqueda usando el HotelID
	filter := bson.M{"hotel_id": hotelObjectId}
	cursor, err := db.Collection("images").Find(context.TODO(), filter)
	if err != nil {
		return nil, errors.NewBadRequestApiError("Error retrieving images: ")
	}
	defer cursor.Close(context.TODO())

	// Itera sobre el cursor para decodificar cada documento
	for cursor.Next(context.TODO()) {
		var img models.Image
		if err := cursor.Decode(&img); err != nil {
			continue // o maneja el error como creas necesario
		}
		images = append(images, img)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.NewBadRequestApiError("Cursor error")
	}

	return images, nil
}
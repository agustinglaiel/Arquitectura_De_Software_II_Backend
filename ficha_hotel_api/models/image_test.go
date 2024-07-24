package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateImage(t *testing.T) {
	assert := assert.New(t)
	imageID := primitive.NewObjectID()
	hotelID := primitive.NewObjectID()

	image := Image{
		Id:      imageID,
		Imagen:  []byte("imagen de prueba"),
		HotelID: hotelID,
	}

	assert.Equal(imageID, image.Id, "El ID de la imagen no coincide")
	assert.Equal([]byte("imagen de prueba"), image.Imagen, "Los datos de la imagen no coinciden")
	assert.Equal(hotelID, image.HotelID, "El ID del hotel asociado no coincide")
}
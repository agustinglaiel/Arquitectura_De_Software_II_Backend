package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestImageDtoSerialization(t *testing.T) {
	assert := assert.New(t)
	hotelID := primitive.NewObjectID()

	imageDto := ImageDto{
		HotelId: hotelID,
		Data:    []byte("data de imagen"),
	}

	// Serialización a JSON
	jsonData, err := json.Marshal(imageDto)
	assert.Nil(err, "Error al serializar ImageDto")

	// Deserialización de JSON
	var decodedImageDto ImageDto
	err = json.Unmarshal(jsonData, &decodedImageDto)
	assert.Nil(err, "Error al deserializar JSON a ImageDto")

	// Comprobar que los datos son iguales después de la serialización y deserialización
	assert.Equal(imageDto.HotelId, decodedImageDto.HotelId, "El HotelId de ImageDto no coincide")
	assert.Equal(imageDto.Data, decodedImageDto.Data, "Los datos de la imagen no coinciden")
}

func TestImagesDtoSerialization(t *testing.T) {
	assert := assert.New(t)
	hotelID := primitive.NewObjectID()

	imagesDto := ImagesDto{
		Images: []ImageDto{
			{
				HotelId: hotelID,
				Data:    []byte("data de imagen"),
			},
		},
	}

	// Serialización a JSON
	jsonData, err := json.Marshal(imagesDto)
	assert.Nil(err, "Error al serializar ImagesDto")

	// Deserialización de JSON
	var decodedImagesDto ImagesDto
	err = json.Unmarshal(jsonData, &decodedImagesDto)
	assert.Nil(err, "Error al deserializar JSON a ImagesDto")

	// Comprobar que los datos son iguales después de la serialización y deserialización
	assert.Equal(imagesDto.Images[0].HotelId, decodedImagesDto.Images[0].HotelId, "El HotelId de ImagesDto no coincide")
	assert.Equal(imagesDto.Images[0].Data, decodedImagesDto.Images[0].Data, "Los datos de la imagen no coinciden")
}
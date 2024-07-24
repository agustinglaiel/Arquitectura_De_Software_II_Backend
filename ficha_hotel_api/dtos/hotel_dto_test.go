package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotelDtoSerialization(t *testing.T) {
	assert := assert.New(t)
	hotelDto := HotelDto{
		ID:             "123",
		Name:           "Hotel Test",
		Description:    "Descripción de prueba",
		Amenities:      []string{"WiFi", "Pool"},
		RoomCount:      100,
		City:           "Ciudad Test",
		AvailableRooms: 80,
	}

	// Serialización a JSON
	jsonData, err := json.Marshal(hotelDto)
	assert.Nil(err, "Error al serializar HotelDto")

	// Deserialización de JSON
	var decodedHotelDto HotelDto
	err = json.Unmarshal(jsonData, &decodedHotelDto)
	assert.Nil(err, "Error al deserializar JSON a HotelDto")

	// Comprobar que los datos son iguales después de la serialización y deserialización
	assert.Equal(hotelDto, decodedHotelDto, "El HotelDto original y el deserializado no son iguales")
}

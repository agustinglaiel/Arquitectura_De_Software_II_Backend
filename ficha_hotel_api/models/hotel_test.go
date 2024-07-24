package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateHotel(t *testing.T){
	assert := assert.New(t)
	id := primitive.NewObjectID()

	hotel := Hotel{
		ID:          id,
		Name:        "test",
		Description: "test",
		City:      "test",
	}

	assert.Equal(id, hotel.ID, "El ID del hotel no coincide")
	assert.Equal("test", hotel.Name, "El nombre del hotel no coincide")
	assert.Equal("test", hotel.Description, "La descripcion del hotel no coincide")
	assert.Equal("test", hotel.City, "La City del hotel no coincide")
}
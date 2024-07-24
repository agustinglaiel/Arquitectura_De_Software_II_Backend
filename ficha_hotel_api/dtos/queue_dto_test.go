package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueDtoSerialization(t *testing.T) {
	assert := assert.New(t)

	queueDto := QueueDto{
		Id:     "queue123",
		Action: "update",
	}

	// Serialización a JSON
	jsonData, err := json.Marshal(queueDto)
	assert.Nil(err, "Error al serializar QueueDto")

	// Deserialización de JSON
	var decodedQueueDto QueueDto
	err = json.Unmarshal(jsonData, &decodedQueueDto)
	assert.Nil(err, "Error al deserializar JSON a QueueDto")

	// Comprobar que los datos son iguales después de la serialización y deserialización
	assert.Equal(queueDto.Id, decodedQueueDto.Id, "El ID de QueueDto no coincide")
	assert.Equal(queueDto.Action, decodedQueueDto.Action, "La acción de QueueDto no coincide")
}

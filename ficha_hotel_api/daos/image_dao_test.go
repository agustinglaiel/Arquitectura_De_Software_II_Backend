package daos

import (
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define la interfaz para el DAO de imágenes
type ImageDaoInterface interface {
	InsertImage(image models.Image) (models.Image, errors.ApiError)
	GetImagesByHotelId(hotelId string) ([]models.Image, errors.ApiError)
}

// Implementa el mock basado en la interfaz
type MockImageDao struct {
	mock.Mock
}

func (m *MockImageDao) InsertImage(image models.Image) (models.Image, errors.ApiError) {
	args := m.Called(image)
	return args.Get(0).(models.Image), args.Get(1).(errors.ApiError)
}

func (m *MockImageDao) GetImagesByHotelId(hotelId string) ([]models.Image, errors.ApiError) {
	args := m.Called(hotelId)
	return args.Get(0).([]models.Image), args.Get(1).(errors.ApiError)
}

// Test para InsertImage
func TestInsertImage(t *testing.T) {
    mockDao := new(MockImageDao)
    testImage := models.Image{HotelID: primitive.NewObjectID()}
    createdImage := testImage
    createdImage.Id = primitive.NewObjectID()  // Simulando el ID asignado por la base de datos

    mockDao.On("InsertImage", testImage).Return(createdImage, noError())

    result, err := mockDao.InsertImage(testImage)
    assert.Equal(t, 200, err.Status()) // Verificar que el estado del error es 200, lo que significa "no error"
    assert.Equal(t, createdImage, result)
    mockDao.AssertExpectations(t)
}

// Test para GetImagesByHotelId
func TestGetImagesByHotelId(t *testing.T) {
    mockDao := new(MockImageDao)
    hotelId := "abc123"
    testImages := []models.Image{
        {Id: primitive.NewObjectID(), HotelID: primitive.NewObjectID()},
        {Id: primitive.NewObjectID(), HotelID: primitive.NewObjectID()},
    }

    mockDao.On("GetImagesByHotelId", hotelId).Return(testImages, noError())

    result, apiErr := mockDao.GetImagesByHotelId(hotelId)
    assert.Equal(t, 200, apiErr.Status()) // Asegurarse de que el estado es 200, indicando "no error"
    assert.Equal(t, testImages, result)
    mockDao.AssertExpectations(t)
}


func noError() errors.ApiError {
    // Usar un error existente con un mensaje neutro y estado 200
    return errors.NewApiError("No error", "no_error", 200, nil)
}


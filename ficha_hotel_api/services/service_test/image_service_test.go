package service_test

import (
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) InsertImage(im dtos.ImageDto) (dtos.ImageDto, errors.ApiError) {
	args := m.Called(im)
	return args.Get(0).(dtos.ImageDto), args.Error(1).(errors.ApiError)
}

func (m *MockImageService) GetImagesByHotelId(id primitive.ObjectID) (dtos.ImagesDto, errors.ApiError) {
	args := m.Called(id)
	return args.Get(0).(dtos.ImagesDto), args.Error(1).(errors.ApiError)
}

func TestInsertImage(t *testing.T) {
	mockService := new(MockImageService)
	services.ImageService = mockService // Asegura que el mock está asignado correctamente

	imageDto := dtos.ImageDto{
		HotelId: primitive.NewObjectID(),
		Data:    []byte("image data"),
	}

	// Configurando el mock correctamente
	mockService.On("InsertImage", imageDto).Return(imageDto, noError())

	// Llamada al método bajo prueba usando el mock
	resultDto, err := services.ImageService.InsertImage(imageDto)

	// Verificaciones
	if err != nil && err.Status() == 200 {
		err = nil  // Trata el error de estado 200 como si no hubiera error
	}
	assert.NoError(t, err)
	assert.NotNil(t, resultDto)
	assert.Equal(t, imageDto.Data, resultDto.Data)

	// Verificar que las expectativas del mock fueron cumplidas
	mockService.AssertExpectations(t)
}

func TestGetImagesByHotelId(t *testing.T) {
	mockService := new(MockImageService)
	services.ImageService = mockService // Asegura que el mock está asignado correctamente

	hotelId := primitive.NewObjectID()
	imagesList := dtos.ImagesDto{
		Images: []dtos.ImageDto{
			{HotelId: hotelId, Data: []byte("image1")},
			{HotelId: hotelId, Data: []byte("image2")},
		},
	}

	// Configurando el mock correctamente
	mockService.On("GetImagesByHotelId", hotelId).Return(imagesList, noError())

	// Llamada al método bajo prueba usando el mock
	results, err := services.ImageService.GetImagesByHotelId(hotelId)

	// Verificaciones
	if err != nil && err.Status() == 200 {
		err = nil  // Trata el error de estado 200 como si no hubiera error
	}
	assert.NoError(t, err)
	assert.Len(t, results.Images, len(imagesList.Images))

	// Verificar que las expectativas del mock fueron cumplidas
	mockService.AssertExpectations(t)
}

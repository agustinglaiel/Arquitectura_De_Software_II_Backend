package controller_test

import (
	"bytes"
	"encoding/json"
	"ficha_hotel_api/controllers"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definición del mock del servicio de imágenes
type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) InsertImage(img dtos.ImageDto) (dtos.ImageDto, errors.ApiError) {
	args := m.Called(img)
	if args.Get(0) == nil {
		return dtos.ImageDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dtos.ImageDto), nil
}

func (m *MockImageService) GetImagesByHotelId(hotelId primitive.ObjectID) (dtos.ImagesDto, errors.ApiError) {
	args := m.Called(hotelId)
	if args.Get(0) == nil {
		return dtos.ImagesDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dtos.ImagesDto), nil
}

// Test para InsertImage
func TestInsertImage(t *testing.T) {
	mockService := new(MockImageService)
	services.ImageService = mockService

	imageData := []byte("fake_image_data")
	hotelId := primitive.NewObjectID()

	mockService.On("InsertImage", mock.AnythingOfType("dtos.ImageDto")).Return(dtos.ImageDto{HotelId: hotelId}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: hotelId.Hex()}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/images/"+hotelId.Hex(), bytes.NewReader(imageData))
	c.Request.Header.Set("Content-Type", "image/jpeg")

	controllers.InsertImage(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Imagen insertada con éxito", response["message"])
	assert.Equal(t, hotelId.Hex(), response["imageId"])

	mockService.AssertExpectations(t)
}

// Test para GetImagesByHotelId
func TestGetImagesByHotelId(t *testing.T) {
	mockService := new(MockImageService)
	services.ImageService = mockService

	hotelId := primitive.NewObjectID()
	imagesList := dtos.ImagesDto{
		Images: []dtos.ImageDto{{HotelId: hotelId, Data: []byte("some_image_data")}},
	}

	mockService.On("GetImagesByHotelId", hotelId).Return(imagesList, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: hotelId.Hex()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/images/"+hotelId.Hex(), nil)

	controllers.GetImagesByHotelId(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var returnedImages dtos.ImagesDto
	err := json.NewDecoder(w.Body).Decode(&returnedImages)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(returnedImages.Images))
	assert.Equal(t, imagesList, returnedImages)

	mockService.AssertExpectations(t)
}

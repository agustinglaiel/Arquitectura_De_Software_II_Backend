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
)

type MockHotelService struct {
	mock.Mock
}

func (m *MockHotelService) GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return dtos.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dtos.HotelDto), nil
}

func (m *MockHotelService) InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
    args := m.Called(hotelDto)
	if args.Get(0) == nil {
		return dtos.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dtos.HotelDto), nil
}

func (m *MockHotelService) UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
    args := m.Called(id, hotelDto)
	if args.Get(0) == nil {
		return dtos.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dtos.HotelDto), nil
}

func (m *MockHotelService) GetHotels() ([]dtos.HotelDto, errors.ApiError) {
    args := m.Called()
	if args.Get(0) == nil {
		return args.Get(0).([]dtos.HotelDto), args.Error(1).(errors.ApiError)
	}
	return args.Get(0).([]dtos.HotelDto), nil
}

func (m *MockHotelService) DeleteHotelById(id string) errors.ApiError {
    args := m.Called(id)
	errArg := args.Get(0)
	if errArg != nil {
		if apiErr, ok := errArg.(errors.ApiError); ok {
			return apiErr
		} else {
			return errors.NewInternalServerApiError("Error casting to ApiError", nil)
		}
	}
	return nil
}

func TestGetHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    testId := "123" // Asume que es un ID válido
    expectedHotelDto := dtos.HotelDto{
        ID:          testId,
        Name:        "Test Hotel",
        Description: "A nice place",
        City:        "Test City",
        Amenities:   []string{"Pool", "Gym"},
    }

    // Usa nil para simular la ausencia de errores
    mockService.On("GetHotelById", testId).Return(expectedHotelDto, nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: testId}}

    controllers.GetHotelById(c)

    assert.Equal(t, http.StatusOK, w.Code)
    var hotel dtos.HotelDto
    err := json.Unmarshal(w.Body.Bytes(), &hotel)
    assert.NoError(t, err) // Verifica que no hay errores en la decodificación
    assert.Equal(t, expectedHotelDto, hotel)

    mockService.AssertExpectations(t)
}

func TestInsertHotel(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    hotelDto := dtos.HotelDto{
        Name:        "New Hotel",
        Description: "Brand new hotel",
        City:        "New City",
        Amenities:   []string{"Spa", "Parking"},
    }

    // Asumiendo que la inserción es exitosa y no hay errores
    mockService.On("InsertHotel", mock.AnythingOfType("dtos.HotelDto")).Return(hotelDto, nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest(http.MethodPost, "/hotels", bytes.NewBufferString(`{
        "name": "New Hotel",
        "description": "Brand new hotel",
        "city": "New City",
        "amenities": ["Spa", "Parking"]
    }`))
    c.Request.Header.Set("Content-Type", "application/json")

    controllers.InsertHotel(c)

    // Verificar que el código de estado HTTP es el esperado para una creación exitosa
    assert.Equal(t, http.StatusCreated, w.Code)

    var returnedHotel dtos.HotelDto
    err := json.NewDecoder(w.Body).Decode(&returnedHotel)
    assert.NoError(t, err)
    assert.Equal(t, hotelDto.Name, returnedHotel.Name) // Verifica que el nombre del hotel sea el esperado

    mockService.AssertExpectations(t)
}

func TestUpdateHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService

    hotelId := "123" // Asegúrate de que este es el mismo ID que se espera en la prueba
    updatedHotelDto := dtos.HotelDto{
        ID:          hotelId,
        Name:        "Updated Hotel",
        Description: "Recently updated features",
        City:        "Updated City",
        Amenities:   []string{"Spa", "Gym"},
    }

    // Asegúrate de que la configuración del mock coincida exactamente con lo que pasas en la prueba
    mockService.On("UpdateHotelById", hotelId, updatedHotelDto).Return(updatedHotelDto, nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: hotelId}}
    c.Request, _ = http.NewRequest(http.MethodPut, "/hotels/"+hotelId, bytes.NewBuffer(jsonEncode(updatedHotelDto)))
    c.Request.Header.Set("Content-Type", "application/json")

    controllers.UpdateHotelById(c)

    assert.Equal(t, http.StatusOK, w.Code)
    var returnedHotel dtos.HotelDto
    err := json.NewDecoder(w.Body).Decode(&returnedHotel)
    assert.NoError(t, err)
    assert.Equal(t, updatedHotelDto.Name, returnedHotel.Name)

    mockService.AssertExpectations(t)
}

func TestGetHotels(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService

    hotelsList := []dtos.HotelDto{
        {ID: "1", Name: "Hotel One", City: "City One"},
        {ID: "2", Name: "Hotel Two", City: "City Two"},
    }

    mockService.On("GetHotels").Return(hotelsList, nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest(http.MethodGet, "/hotels", nil)

    controllers.GetHotels(c)

    assert.Equal(t, http.StatusOK, w.Code)
    var returnedHotels []dtos.HotelDto
    err := json.NewDecoder(w.Body).Decode(&returnedHotels)
    assert.NoError(t, err)
    assert.Equal(t, len(hotelsList), len(returnedHotels))
    assert.Equal(t, hotelsList, returnedHotels)

    mockService.AssertExpectations(t)
}

func TestDeleteHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService

    hotelId := "123"

    mockService.On("DeleteHotelById", hotelId).Return(nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: hotelId}}
    c.Request, _ = http.NewRequest(http.MethodDelete, "/hotels/"+hotelId, nil)

    controllers.DeleteHotelById(c)

    assert.Equal(t, http.StatusNoContent, w.Code) // Assuming you're using 204 No Content for successful deletions

    mockService.AssertExpectations(t)
}

// Helper function to encode the hotel DTO to JSON
func jsonEncode(hotel dtos.HotelDto) []byte {
    data, _ := json.Marshal(hotel)
    return data
}


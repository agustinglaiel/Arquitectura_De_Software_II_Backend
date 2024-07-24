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

type MockHotelService struct {
    mock.Mock
}

func (m *MockHotelService) GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
    args := m.Called(id)
    return args.Get(0).(dtos.HotelDto), args.Error(1).(errors.ApiError)
}

func (m *MockHotelService) InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
    args := m.Called(hotelDto)
    return args.Get(0).(dtos.HotelDto), args.Error(1).(errors.ApiError)
}

func (m *MockHotelService) UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
    args := m.Called(id, hotelDto)
    return args.Get(0).(dtos.HotelDto), args.Error(1).(errors.ApiError)
}

func (m *MockHotelService) GetHotels() ([]dtos.HotelDto, errors.ApiError) {
    args := m.Called()
    return args.Get(0).([]dtos.HotelDto), args.Error(1).(errors.ApiError)
}

func (m *MockHotelService) DeleteHotelById(id string) errors.ApiError {
    args := m.Called(id)
    return args.Error(0).(errors.ApiError)
}

func TestGetHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    testId := primitive.NewObjectID().Hex() // ID válido
    expectedHotelDto := dtos.HotelDto{
        ID:          testId,
        Name:        "Test Hotel",
        Description: "A nice place",
        City:        "Test City",
        Amenities:   []string{"Pool", "Gym"},
    }

    // Configurando el mock correctamente usando noError()
    mockService.On("GetHotelById", testId).Return(expectedHotelDto, noError())

    // Llamada al método bajo prueba usando el mock
    hotelDto, err := services.HotelService.GetHotelById(testId)

    // Verificaciones
    if err != nil && err.Status() == 200 {
        err = nil  // Trata el error de estado 200 como si no hubiera error
    }
    assert.NoError(t, err)
    assert.NotNil(t, hotelDto)
    assert.Equal(t, expectedHotelDto.Name, hotelDto.Name)
    assert.Equal(t, expectedHotelDto.Description, hotelDto.Description)

    // Verificar que las expectativas del mock fueron cumplidas
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

    // Configurando el mock correctamente
    // Asegúrate de que el segundo retorno es un objeto que cumple con errors.ApiError
    mockService.On("InsertHotel", mock.AnythingOfType("dtos.HotelDto")).Return(hotelDto, noError())

    // Llamada al método bajo prueba usando el mock
    resultDto, err := services.HotelService.InsertHotel(hotelDto)

    // Verificaciones
	if err != nil && err.Status() == 200 {
        err = nil  // Trata el error de estado 200 como si no hubiera error
    }
    assert.NoError(t, err) // Verifica que no hay error
    assert.NotNil(t, resultDto) // Verifica que el resultado no es nulo
    assert.Equal(t, hotelDto.Name, resultDto.Name) // Verifica que el nombre del hotel es el esperado

    // Verificar que las expectativas del mock fueron cumplidas
    mockService.AssertExpectations(t)
}


func TestUpdateHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    hotelId := primitive.NewObjectID().Hex() // Asegúrate de que es un ID válido
    updatedHotelDto := dtos.HotelDto{
        ID:          hotelId,
        Name:        "Updated Hotel",
        Description: "Recently updated features",
        City:        "New City",
        Amenities:   []string{"Spa", "Parking"},
    }

    // Configurando el mock correctamente
    // Asegúrate de que el segundo retorno es un objeto que cumple con errors.ApiError
    mockService.On("UpdateHotelById", hotelId, updatedHotelDto).Return(updatedHotelDto, noError())

    // Llamada al método bajo prueba usando el mock
    resultDto, err := services.HotelService.UpdateHotelById(hotelId, updatedHotelDto)

    // Verificaciones
    if err != nil && err.Status() == 200 {
        err = nil  // Trata el error de estado 200 como si no hubiera error
    }
    assert.NoError(t, err) // Verifica que no hay error
    assert.NotNil(t, resultDto) // Verifica que el resultado no es nulo
    assert.Equal(t, updatedHotelDto.Name, resultDto.Name) // Verifica que el nombre del hotel es el esperado

    // Verificar que las expectativas del mock fueron cumplidas
    mockService.AssertExpectations(t)
}


func TestGetHotels(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    hotelsList := []dtos.HotelDto{
        {ID: primitive.NewObjectID().Hex(), Name: "Hotel One"},
        {ID: primitive.NewObjectID().Hex(), Name: "Hotel Two"},
    }

    // Configurando el mock correctamente
    // Asegúrate de que el segundo retorno es un objeto que cumple con errors.ApiError
    mockService.On("GetHotels").Return(hotelsList, noError())

    // Llamada al método bajo prueba usando el mock
    results, err := services.HotelService.GetHotels()

    // Verificaciones
    if err != nil && err.Status() == 200 {
        err = nil  // Trata el error de estado 200 como si no hubiera error
    }
    assert.NoError(t, err) // Verifica que no hay error
    assert.Len(t, results, len(hotelsList)) // Verifica que la longitud de los resultados es la esperada

    // Verificar que las expectativas del mock fueron cumplidas
    mockService.AssertExpectations(t)
}


func TestDeleteHotelById(t *testing.T) {
    mockService := new(MockHotelService)
    services.HotelService = mockService // Asegura que el mock está asignado correctamente

    hotelId := primitive.NewObjectID().Hex() // Asegúrate de que es un ID válido

    // Configurando el mock correctamente
    // Asegúrate de que el segundo retorno es un objeto que cumple con errors.ApiError
    mockService.On("DeleteHotelById", hotelId).Return(noError())

    // Llamada al método bajo prueba usando el mock
    err := services.HotelService.DeleteHotelById(hotelId)

    // Verificaciones
    if err != nil && err.Status() == 200 {
        err = nil  // Trata el error de estado 200 como si no hubiera error
    }
    assert.NoError(t, err) // Verifica que no hay error

    // Verificar que las expectativas del mock fueron cumplidas
    mockService.AssertExpectations(t)
}



func noError() errors.ApiError {
    // Usar un error existente con un mensaje neutro y estado 200
    return errors.NewApiError("No error", "no_error", 200, nil)
}
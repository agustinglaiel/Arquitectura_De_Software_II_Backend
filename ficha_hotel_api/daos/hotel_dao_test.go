package daos

import (
	"ficha_hotel_api/models"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definición de la interfaz para el DAO
type HotelDaoInterface interface {
	GetHotelById(id string) models.Hotel
	InsertHotel(hotel models.Hotel) models.Hotel
	UpdateHotel(hotel models.Hotel) models.Hotel
	GetHotels() ([]models.Hotel, error)
	DeleteHotelById(id string) error
}

// Implementación del mock basado en la interfaz
type MockHotelDao struct {
	mock.Mock
}

func (m *MockHotelDao) GetHotelById(id string) models.Hotel {
	args := m.Called(id)
	return args.Get(0).(models.Hotel)
}

func (m *MockHotelDao) InsertHotel(hotel models.Hotel) models.Hotel {
	args := m.Called(hotel)
	return args.Get(0).(models.Hotel)
}

func (m *MockHotelDao) UpdateHotel(hotel models.Hotel) models.Hotel {
	args := m.Called(hotel)
	return args.Get(0).(models.Hotel)
}

func (m *MockHotelDao) GetHotels() ([]models.Hotel, error) {
	args := m.Called()
	return args.Get(0).([]models.Hotel), args.Error(1)
}

func (m *MockHotelDao) DeleteHotelById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetHotelById(t *testing.T) {
	mockDao := new(MockHotelDao)
	expectedHotel := models.Hotel{
		ID:             primitive.NewObjectID(),
		Name:           "Test Hotel",
		Description:    "A test hotel",
		Amenities:      []string{"Pool", "Spa"},
		RoomCount:      100,
		City:           "Test City",
		AvailableRooms: 50,
	}

	mockDao.On("GetHotelById", mock.Anything).Return(expectedHotel)

	// Llamada al método que se quiere testear
	result := mockDao.GetHotelById("abc123")
	assert.Equal(t, expectedHotel, result)
	mockDao.AssertExpectations(t)
}

func TestInsertHotel(t *testing.T) {
	mockDao := new(MockHotelDao)
	newHotel := models.Hotel{
		Name:           "New Hotel",
		Description:    "Newly opened hotel",
		Amenities:      []string{"WiFi", "Gym"},
		RoomCount:      200,
		City:           "New City",
		AvailableRooms: 180,
	}

	// Se espera que el ID se genere automáticamente dentro de la función
	createdHotel := newHotel
	createdHotel.ID = primitive.NewObjectID() // Simulando el ID asignado por la base de datos

	mockDao.On("InsertHotel", newHotel).Return(createdHotel)

	// Llamada al método que se quiere testear
	result := mockDao.InsertHotel(newHotel)
	assert.Equal(t, createdHotel, result)
	mockDao.AssertExpectations(t)
}

func TestUpdateHotel(t *testing.T) {
	mockDao := new(MockHotelDao)
	existingHotel := models.Hotel{
		ID:             primitive.NewObjectID(),
		Name:           "Updated Hotel",
		Description:    "Updated description",
		Amenities:      []string{"Updated Amenities"},
		RoomCount:      150,
		City:           "Updated City",
		AvailableRooms: 120,
	}

	mockDao.On("UpdateHotel", existingHotel).Return(existingHotel)

	// Llamada al método que se quiere testear
	result := mockDao.UpdateHotel(existingHotel)
	assert.Equal(t, existingHotel, result)
	mockDao.AssertExpectations(t)
}

func TestGetHotels(t *testing.T) {
	assertion := assert.New(t)  // Creando una instancia de assert
	mockDao := new(MockHotelDao)
	hotelsList := []models.Hotel{
		{
			ID:             primitive.NewObjectID(),
			Name:           "Hotel One",
			Description:    "Description One",
			Amenities:      []string{"Pool", "Gym"},
			RoomCount:      100,
			City:           "City One",
			AvailableRooms: 80,
		},
		{
			ID:             primitive.NewObjectID(),
			Name:           "Hotel Two",
			Description:    "Description Two",
			Amenities:      []string{"Spa", "Restaurant"},
			RoomCount:      150,
			City:           "City Two",
			AvailableRooms: 140,
		},
	}

	mockDao.On("GetHotels").Return(hotelsList, nil)

	// Llamada al método que se quiere testear
	result, err := mockDao.GetHotels()
	assertion.Nil(err)
	assertion.Equal(hotelsList, result)
	mockDao.AssertExpectations(t)
}


func TestDeleteHotelById(t *testing.T) {
	assertion := assert.New(t)  // Creando una instancia de assert
	mockDao := new(MockHotelDao)
	hotelID := "abc123"
	mockDao.On("DeleteHotelById", hotelID).Return(nil)

	// Llamada al método que se quiere testear
	err := mockDao.DeleteHotelById(hotelID)
	assertion.Nil(err)
	mockDao.AssertExpectations(t)
}


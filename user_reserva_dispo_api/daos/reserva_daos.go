package daos

import (
	"log"
	"time"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"

	"github.com/jinzhu/gorm"
)

//var Db *gormDB TENEMOS QUE VER EL ERROR ESTE @JULI JAJJA

// CreateReservation agrega una nueva reserva en la base de datos
func CreateReservation(reservation models.Reservation) (models.Reservation, errors.ApiError) {
	result := Db.Create(&reservation)
	if result.Error != nil {
		return models.Reservation{}, errors.NewInternalServerApiError("Failed to create reservation", result.Error)
	}
	return reservation, nil
}

// UpdateReservation actualiza una reserva existente
func UpdateReservation(reservation models.Reservation) (models.Reservation, errors.ApiError) {
	result := Db.Save(&reservation)
	if result.Error != nil {
		return models.Reservation{}, errors.NewInternalServerApiError("Failed to update reservation", result.Error)
	}
	return reservation, nil
}

// DeleteReservation elimina una reserva por su ID
func DeleteReservation(reservationID int) errors.ApiError {
	result := Db.Delete(&models.Reservation{}, reservationID)
	if result.Error != nil {
		return errors.NewInternalServerApiError("Failed to delete reservation", result.Error)
	}
	return nil
}

// GetReservationsByUserId obtiene todas las reservas de un usuario
func GetReservationsByUserId(userID int) ([]models.Reservation, errors.ApiError) {
	var reservations []models.Reservation
	result := Db.Where("user_id = ?", userID).Find(&reservations)
	if result.Error != nil {
		return nil, errors.NewInternalServerApiError("Failed to fetch reservations", result.Error)
	}
	return reservations, nil
}

// GetReservationById obtiene una reserva por su ID
func GetReservationById(reservationID int) (models.Reservation, errors.ApiError) {
	var reservation models.Reservation
	result := Db.First(&reservation, reservationID)
	if result.Error != nil {
		return models.Reservation{}, errors.NewNotFoundApiError("Reservation not found")
	}
	return reservation, nil
}

// GetAllReservations obtiene todas las reservas (para admin)
func GetAllReservations() ([]models.Reservation, errors.ApiError) {
	var reservations []models.Reservation
	result := Db.Find(&reservations)
	if result.Error != nil {
		return nil, errors.NewInternalServerApiError("Failed to fetch all reservations", result.Error)
	}
	return reservations, nil
}

// CheckAvailability verifica la disponibilidad para un hotel en fechas específicas
func CheckAvailability(hotelID string, startDate, endDate time.Time) (models.Reservations, errors.ApiError) {

	var result models.Reservations
	Db.Where("? < end_date AND hotel_id = ? ", startDate, hotelID).Find(&result)
	if Db.Error != nil {
		log.Fatal(Db.Error.Error())
	}
	return result, nil
}

func CheckHotelExists(hotelID string) (bool, error) {
	var hotel models.Hotel
	result := Db.Where("id_amadeus = ?", hotelID).First(&hotel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // El hotel no existe
		}
		return false, result.Error // Error en la consulta
	}
	return true, nil // El hotel existe
}

func InsertHotel(hotel dtos.HotelDto) (models.Hotel, error) {
	var mHotel models.Hotel
	mHotel.IdMongo = hotel.ID
	mHotel.IdAmadeus = hotel.IdAmadeus

	result := Db.Create(&mHotel)
	if result.Error != nil {
		return models.Hotel{}, errors.NewInternalServerApiError("Failed to create Hotel", result.Error)
	}
	return mHotel, nil
}

func GetHotelById(id string) (models.Hotel, error) {

	var hotel models.Hotel
	result := Db.First(&hotel, hotel.IdMongo)
	if result.Error != nil {
		return models.Hotel{}, result.Error
	}
	return hotel, nil
}

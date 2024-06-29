package daos

import (
	"time"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"
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
func CheckAvailability(hotelID int, startDate, endDate time.Time) (bool, errors.ApiError) {
	var count int
	result := Db.Model(&models.Reservation{}).Where("hotel_id = ? AND start_date <= ? AND end_date >= ?", hotelID, endDate, startDate).Count(&count)
	if result.Error != nil {
		return false, errors.NewInternalServerApiError("Failed to check availability", result.Error)
	}
	return count == 0, nil // Retorna true si no hay reservas que se superpongan
}
package dtos

import (
	"time"
)

// ReservationDto representa una reserva con detalle completo.
type ReservationDto struct {
	ReservationID int       `json:"reservation_id"`
	UserID        int       `json:"user_id"`
	HotelID       int       `json:"hotel_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

type ReservationsDto []ReservationDto

// CreateReservationDto se utiliza para la creación de reservas.
type CreateReservationDto struct {
	UserID    int       `json:"user_id"`
	HotelID   string    `json:"hotel_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// UpdateReservationDto se utiliza para actualizar una reserva existente.
type UpdateReservationDto struct {
	ReservationID int       `json:"reservation_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

// ReservationAvailabilityDto se utiliza para verificar la disponibilidad de habitaciones.
type ReservationAvailabilityDto struct {
	HotelID   string    `json:"hotel_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

package models

import (
	"time"
)

type Reservation struct {
	ReservationID int       `gorm:"primaryKey;autoIncrement"`
	UserID        int       `gorm:"not null"`
	HotelID       string    `gorm:"not null"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
	Status        string    `gorm:"size:100;not null;default:'pending'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"` // Opcional, dependiendo de tus necesidades
	UpdatedAt     time.Time `gorm:"autoUpdateTime"` // Opcional, dependiendo de tus necesidades
}

type Reservations []Reservation

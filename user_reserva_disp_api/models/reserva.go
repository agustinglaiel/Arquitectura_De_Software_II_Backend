package models

import (
	"time"
)

type Reservation struct {
	ReservationID int       `gorm:"primaryKey;autoIncrement"`
	UserID        int       `gorm:"not null"`
	HotelID       int       `gorm:"not null"`
	RoomID        int       `gorm:"not null"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
	NumberOfGuests int      `gorm:"not null"`
	Status        string    `gorm:"size:100;not null;default:'pending'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

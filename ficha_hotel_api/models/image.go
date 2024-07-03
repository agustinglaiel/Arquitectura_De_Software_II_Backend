package models

type Image struct {
	Id      int    `gorm:"primaryKey"`
	Imagen  []byte `gorm:"type:longblob;not null"`
	HotelID string    `gorm:"not null"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID"`
}

type Images []Image
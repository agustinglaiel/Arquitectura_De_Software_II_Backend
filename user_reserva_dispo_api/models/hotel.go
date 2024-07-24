package models

type Hotel struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	IdMongo   string `gorm:"not null"`
	IdAmadeus string `gorm:"not null"`
}

type Hotels []Hotel

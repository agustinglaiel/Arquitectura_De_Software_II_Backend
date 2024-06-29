package models

type User struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	FirstName  string `gorm:"size:255;not null"`
	LastName   string `gorm:"size:255;not null"`
	Username   string `gorm:"size:255;not null;unique"`
	Password   string `gorm:"size:255;not null"`
	Email      string `gorm:"size:255;not null;unique"`
	//Phone      string `gorm:"size:50"`
	//Address    string `gorm:"size:255"`
	IsAdmin    bool   `gorm:"default:false"`
}

type Users []User 

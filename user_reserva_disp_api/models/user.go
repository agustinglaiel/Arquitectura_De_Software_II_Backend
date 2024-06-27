package models

type User struct {
	ID        int `json:"id" gorm:"_id,omitempty"`
	FirstName string `json:"first_name" gorm:"first_name"`
	LastName  string `json:"last_name" gorm:"last_name"`
	Email     string `json:"email" gorm:"email"`
	Password  string `json:"password" gorm:"password"`
	Type      bool   `json:"type" gorm:"type"` // True para admin
}

type Users []User

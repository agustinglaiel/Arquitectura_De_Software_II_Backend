package models

type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
	Type      bool   `json:"type" bson:"type"` // True for admin, False for regular user
}

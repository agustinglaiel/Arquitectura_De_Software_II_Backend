package models

type Reserva struct {
	ID             string `json:"id" bson:"_id,omitempty"`
	UserID         string `json:"user_id" bson:"user_id"`
	HotelID        string `json:"hotel_id" bson:"hotel_id"`
	CheckInDate    string `json:"check_in_date" bson:"check_in_date"`
	CheckOutDate   string `json:"check_out_date" bson:"check_out_date"`
	NumRooms       int    `json:"num_rooms" bson:"num_rooms"`
	NumGuests      int    `json:"num_guests" bson:"num_guests"`
	Status         string `json:"status" bson:"status"` // Ej: "confirmed", "cancelled"
	//CreatedAt      int64  `json:"created_at" bson:"created_at"`
	//UpdatedAt      int64  `json:"updated_at" bson:"updated_at"`
}

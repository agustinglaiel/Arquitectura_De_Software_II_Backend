package models

type Hotel struct {
    ID            string   `json:"id" bson:"_id,omitempty"`
    Name          string   `json:"name" bson:"name"`
    Description   string   `json:"description" bson:"description"`
    Photos        []string `json:"photos" bson:"photos"`
    Amenities     []string `json:"amenities" bson:"amenities"`
    RoomCount     int      `json:"room_count" bson:"room_count"`   
    City          string   `json:"city" bson:"city"`                
    AvailableRooms int     `json:"available_rooms" bson:"available_rooms"` 
}
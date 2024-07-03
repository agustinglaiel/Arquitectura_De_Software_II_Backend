package dtos

type HotelDto struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Amenities      []string `json:"amenities"`
	RoomCount      int      `json:"room_count"`
	City           string   `json:"city"`
	AvailableRooms int      `json:"available_rooms"`
}

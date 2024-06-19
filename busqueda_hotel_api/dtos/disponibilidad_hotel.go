package dtos

type DiponibilidadHotelDto struct{
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Photos         []string `json:"photos"`
	Amenities      []string `json:"amenities"`
	RoomCount      int      `json:"room_count"`
	City           string   `json:"city"`
	Availability   int      `json:"availability"`
}

type DisponibilidadHotelsDto struct{
	HotelsAvailability []DiponibilidadHotelDto `json:"hotels` 
}
package dtos

type HotelDTO struct {
    ID             string   `json:"id"`
    Name           string   `json:"name"`
    Description    string   `json:"description"`
    City           string   `json:"city"`
    Photos         []string `json:"photos"`
    Amenities      []string `json:"amenities"`
    RoomCount      int      `json:"room_count"`
    AvailableRooms int      `json:"available_rooms"`
}

type HotelsDTO []HotelDTO
/*
Para transferir datos entre los microservicios y Solr.
Campos similares al modelo Hotel para indexación y actualización en Solr.
*/
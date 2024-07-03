package dtos

type ImageDto struct {
	HotelId string `json:"hotel_id"`
	Data    []byte
}

type ImagesDto struct {
	Images []ImageDto `json:"images"`
}
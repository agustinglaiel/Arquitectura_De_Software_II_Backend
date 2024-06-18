package services

import (
	dao "ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
	"ficha_hotel_api/utils/queue"
)

type hotelService struct{}

type HotelServiceInterface interface {
	GetHotelById(id string) (dtos.HotelDto, errors.ApiError)
	InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	//DeleteHotel(id string) errors.ApiError
}

var (
	HotelService HotelServiceInterface
)

func init(){
	HotelService = &hotelService{}
}

func (s *hotelService) GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel = dao.GetHotelById(id)
	var hotelDto dtos.HotelDto

	if hotel.ID.Hex() == "000000000000000000000000" {
		return hotelDto, errors.NewBadRequestApiError("Hotel not found")
	}

	hotelDto.ID = hotel.ID.Hex()
	hotelDto.Name = hotel.Name
	hotelDto.Description = hotel.Description
	hotelDto.Photos = hotel.Photos
	hotelDto.Amenities = hotel.Amenities
	hotelDto.RoomCount = hotel.RoomCount
	hotelDto.City = hotel.City
	hotelDto.AvailableRooms = hotel.AvailableRooms

	return hotelDto, nil
}

func (s *hotelService) InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	hotel = dao.InsertHotel(hotel)

	if hotel.ID.Hex() == "000000000000000000000000" {
		return hotelDto, errors.NewBadRequestApiError("Error inserting new hotel")
	}

	hotelDto.ID = hotel.ID.Hex()

	queue.Send(hotelDto.ID)

	return hotelDto, nil
}

func (s *hotelService) UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel = dao.GetHotelById(id)

	if hotel.ID.Hex() == "000000000000000000000000" {
		return hotelDto, errors.NewBadRequestApiError("Hotel not found")
	}

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	dao.UpdateHotel(hotel)
	hotelDto.ID = hotel.ID.Hex()

	queue.Send(hotelDto.ID)

	return hotelDto, nil
}



package services

import (
	"fmt"
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"
)

type hotelService struct{}

type hotelServiceInterface interface {
	GetHotels() (dtos.HotelsDto, errors.ApiError)
	InsertHotel(hotelDto dtos.HotelPostDto, idAmadeus string) (dtos.HotelDto, errors.ApiError)
	GetHotelById(id int) (dtos.HotelDto, errors.ApiError)
	CheckHotelByIdAmadeus(id string) (bool, errors.ApiError)
	// UpdateHotel(updateHotelDto dto.HandleHotelDto) (dto.HotelDto, e.ApiError)
	// DeleteHotel(idHotel int, idUser int) (dto.DeleteHotelResponseDto, e.ApiError)

}

var (
	HotelService hotelServiceInterface
)

func init() {
	HotelService = &hotelService{}
}

func (s *hotelService) GetHotelById(id int) (dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel = daos.GetHotelById(id)
	var hotelDto dtos.HotelDto

	if hotel.Id == 0 {
		return hotelDto, errors.NewBadRequestApiError("Hotel no encontrado")
	}

	hotelDto.Id = hotel.Id
	hotelDto.HotelName = hotel.HotelName
	hotelDto.IdMongo = hotel.IdMongo
	hotelDto.IdAmadeus = hotel.IdAmadeus

	return hotelDto, nil
}

func (s *hotelService) CheckHotelByIdAmadeus(id string) (bool, errors.ApiError) {

	if daos.GetHotelByAmadeusId(id) == true {
		return false, errors.NewBadRequestApiError("Hotel ya en uso")
	}

	return true, nil
}

func (s *hotelService) GetHotels() (dtos.HotelsDto, errors.ApiError) {

	var hotels models.Hotels = daos.GetHotels()
	var hotelsDto dtos.HotelsDto

	for _, hotel := range hotels {
		var hotelDto dtos.HotelDto
		id := hotel.Id

		hotelDto, _ = s.GetHotelById(id)

		hotelsDto = append(hotelsDto, hotelDto)
	}

	return hotelsDto, nil
}

func (s *hotelService) InsertHotel(hotelDto dtos.HotelPostDto, idAmadeus string) (dtos.HotelDto, errors.ApiError) {
	fmt.Println("entro al service")
	var hotel models.Hotel
	var response dtos.HotelDto

	hotel.HotelName = hotelDto.HotelName
	hotel.IdAmadeus = idAmadeus
	hotel.IdMongo = hotelDto.IdMongo

	hotel = daos.InsertHotel(hotel)

	if hotel.Id == 0 {
		return response, errors.NewBadRequestApiError("Error al insertar hotel")
	}

	hotelDto.Id = hotel.Id

	return response, nil
}
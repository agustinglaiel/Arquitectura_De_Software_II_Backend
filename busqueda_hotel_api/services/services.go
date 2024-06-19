package services

import (
	dao "busqueda_hotel_api/daos"
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/models"
	"busqueda_hotel_api/utils/errors"
)

type hotelService struct {
	dao dao.HotelDao
}

type HotelServiceInterface interface {
	GetHotelById(id string) (dtos.HotelDTO, errors.ApiError)
	InsertHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	UpdateHotelById(id string, hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	GetHotels() ([]dtos.HotelDTO, errors.ApiError)
	GetHotelsByCiudad(ciudad string) ([]dtos.HotelDTO, errors.ApiError)
	GetDisponibilidad(ciudad, fechainicio, fechafinal string) ([]dtos.HotelDTO, errors.ApiError)
	DeleteHotelById(id string) errors.ApiError
}

var (
	HotelService HotelServiceInterface
)

func init() {
	HotelService = &hotelService{
		dao: dao.NewHotelSolrDao(),
	}
}

func (s *hotelService) GetHotelById(id string) (dtos.HotelDTO, errors.ApiError) {
	hotel, err := s.dao.Get(id)
	if err != nil {
		return dtos.HotelDTO{}, errors.NewInternalServerApiError("Error fetching hotel", err)
	}

	hotelDto := dtos.HotelDTO{
		ID:             hotel.ID,
		Name:           hotel.Name,
		Description:    hotel.Description,
		City:           hotel.City,
		Photos:         hotel.Photos,
		Amenities:      hotel.Amenities,
		RoomCount:      hotel.RoomCount,
		AvailableRooms: hotel.AvailableRooms,
	}

	return hotelDto, nil
}

func (s *hotelService) InsertHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError) {
	hotel := models.Hotel{
		ID:             hotelDto.ID,
		Name:           hotelDto.Name,
		Description:    hotelDto.Description,
		City:           hotelDto.City,
		Photos:         hotelDto.Photos,
		Amenities:      hotelDto.Amenities,
		RoomCount:      hotelDto.RoomCount,
		AvailableRooms: hotelDto.AvailableRooms,
	}

	err := s.dao.Create(&hotel)
	if err != nil {
		return dtos.HotelDTO{}, errors.NewInternalServerApiError("Error inserting new hotel", err)
	}

	hotelDto.ID = hotel.ID
	return hotelDto, nil
}

func (s *hotelService) UpdateHotelById(id string, hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError) {
	hotel := models.Hotel{
		ID:             id,
		Name:           hotelDto.Name,
		Description:    hotelDto.Description,
		City:           hotelDto.City,
		Photos:         hotelDto.Photos,
		Amenities:      hotelDto.Amenities,
		RoomCount:      hotelDto.RoomCount,
		AvailableRooms: hotelDto.AvailableRooms,
	}

	err := s.dao.Update(&hotel)
	if err != nil {
		return dtos.HotelDTO{}, errors.NewInternalServerApiError("Error updating hotel", err)
	}

	return hotelDto, nil
}

func (s *hotelService) GetHotels() ([]dtos.HotelDTO, errors.ApiError) {
	hotels, err := s.dao.GetAll()
	if err != nil {
		return nil, errors.NewInternalServerApiError("Error fetching hotels", err)
	}

	var hotelDtos []dtos.HotelDTO
	for _, hotel := range hotels {
		hotelDto := dtos.HotelDTO{
			ID:             hotel.ID,
			Name:           hotel.Name,
			Description:    hotel.Description,
			City:           hotel.City,
			Photos:         hotel.Photos,
			Amenities:      hotel.Amenities,
			RoomCount:      hotel.RoomCount,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos = append(hotelDtos, hotelDto)
	}

	return hotelDtos, nil
}

func (s *hotelService) GetHotelsByCiudad(ciudad string) ([]dtos.HotelDTO, errors.ApiError) {
	hotels, err := s.dao.GetByCiudad(ciudad)
	if err != nil {
		return nil, errors.NewInternalServerApiError("Error fetching hotels by city", err)
	}

	var hotelDtos []dtos.HotelDTO
	for _, hotel := range hotels {
		hotelDto := dtos.HotelDTO{
			ID:             hotel.ID,
			Name:           hotel.Name,
			Description:    hotel.Description,
			City:           hotel.City,
			Photos:         hotel.Photos,
			Amenities:      hotel.Amenities,
			RoomCount:      hotel.RoomCount,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos = append(hotelDtos, hotelDto)
	}

	return hotelDtos, nil
}

func (s *hotelService) GetDisponibilidad(ciudad, fechainicio, fechafinal string) ([]dtos.HotelDTO, errors.ApiError) {
	hotels, err := s.dao.GetDisponibilidad(ciudad, fechainicio, fechafinal)
	if err != nil {
		return nil, errors.NewInternalServerApiError("Error fetching hotel availability", err)
	}

	var hotelDtos []dtos.HotelDTO
	for _, hotel := range hotels {
		hotelDto := dtos.HotelDTO{
			ID:             hotel.ID,
			Name:           hotel.Name,
			Description:    hotel.Description,
			City:           hotel.City,
			Photos:         hotel.Photos,
			Amenities:      hotel.Amenities,
			RoomCount:      hotel.RoomCount,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos = append(hotelDtos, hotelDto)
	}

	return hotelDtos, nil
}

func (s *hotelService) DeleteHotelById(id string) errors.ApiError {
	err := s.dao.DeleteById(id)
	if err != nil {
		return errors.NewInternalServerApiError("Error deleting hotel", err)
	}

	return nil
}

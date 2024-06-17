package services

import (
	"context"
	"encoding/json"

	"ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
	"ficha_hotel_api/utils/queue"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelServiceInterface interface {
	CreateHotel(ctx context.Context, hotelDto dtos.HotelDto) (*dtos.HotelDto, errors.ApiError)
	DeleteHotel(ctx context.Context, id primitive.ObjectID) errors.ApiError
	UpdateHotel(ctx context.Context, id primitive.ObjectID, hotelDto dtos.HotelDto) (*dtos.HotelDto, errors.ApiError)
	GetHotelByID(ctx context.Context, id primitive.ObjectID) (*dtos.HotelDto, errors.ApiError)
}

type hotelService struct {
	dao *daos.HotelDAO
}

func NewHotelService(dao *daos.HotelDAO) HotelServiceInterface {
	return &hotelService{dao: dao}
}

func (s *hotelService) CreateHotel(ctx context.Context, hotelDto dtos.HotelDto) (*dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel
	hotel.ID = primitive.NewObjectID().Hex()
	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	err := s.dao.InsertHotel(ctx, hotel)
	if err != nil {
		return nil, errors.NewInternalServerApiError("error when trying to create hotel", err)
	}
	message, _ := json.Marshal(hotel)
	queue.Send(string(message))
	hotelDto.ID = hotel.ID
	return &hotelDto, nil
}

func (s *hotelService) DeleteHotel(ctx context.Context, id primitive.ObjectID) errors.ApiError {
	err := s.dao.DeleteHotel(ctx, id)
	if err != nil {
		return errors.NewInternalServerApiError("error when trying to delete hotel", err)
	}
	return nil
}

func (s *hotelService) UpdateHotel(ctx context.Context, id primitive.ObjectID, hotelDto dtos.HotelDto) (*dtos.HotelDto, errors.ApiError) {
	hotel, err := s.dao.GetHotelByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerApiError("error when trying to get hotel", err)
	}

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	err = s.dao.UpdateHotel(ctx, id, *hotel)
	if err != nil {
		return nil, errors.NewInternalServerApiError("error when trying to update hotel", err)
	}
	message, _ := json.Marshal(hotel)
	queue.Send(string(message))
	hotelDto.ID = hotel.ID
	return &hotelDto, nil
}

func (s *hotelService) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*dtos.HotelDto, errors.ApiError) {
	hotel, err := s.dao.GetHotelByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerApiError("error when trying to get hotel", err)
	}
	var hotelDto dtos.HotelDto
	hotelDto.ID = hotel.ID
	hotelDto.Name = hotel.Name
	hotelDto.Description = hotel.Description
	hotelDto.Photos = hotel.Photos
	hotelDto.Amenities = hotel.Amenities
	hotelDto.RoomCount = hotel.RoomCount
	hotelDto.City = hotel.City
	hotelDto.AvailableRooms = hotel.AvailableRooms
	return &hotelDto, nil
}

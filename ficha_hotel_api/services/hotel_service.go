package services

import (
	"context"
	"ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
	"ficha_hotel_api/utils/queue"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelServiceInterface interface {
	GetHotelById(id string) (dtos.HotelDto, errors.ApiError)
	InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	DeleteHotel(id string) errors.ApiError
}

type hotelService struct {
	db *mongo.Database
}

var (
	HotelService HotelServiceInterface
)

func NewHotelService(db *mongo.Database) HotelServiceInterface {
	return &hotelService{db: db}
}

func (s *hotelService) GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
	var hotelDto dtos.HotelDto

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return hotelDto, errors.NewBadRequestApiError("invalid hotel ID")
	}

	hotel, err := daos.NewHotelDAO(s.db).GetHotelByID(context.Background(), objectID)
	if err != nil {
		return hotelDto, errors.NewNotFoundApiError("hotel not found")
	}

	hotelDto.ID = hotel.ID
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

	hotel.ID = primitive.NewObjectID().Hex()

	err := daos.NewHotelDAO(s.db).InsertHotel(context.Background(), hotel)
	if err != nil {
		return hotelDto, errors.NewInternalServerApiError("error when trying to create hotel", err)
	}

	hotelDto.ID = hotel.ID
	queue.Send(hotelDto.ID)
	return hotelDto, nil
}

func (s *hotelService) UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return hotelDto, errors.NewBadRequestApiError("invalid hotel ID")
	}

	hotel, err := daos.NewHotelDAO(s.db).GetHotelByID(context.Background(), objectID)
	if err != nil {
		return hotelDto, errors.NewNotFoundApiError("hotel not found")
	}

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	err = daos.NewHotelDAO(s.db).UpdateHotel(context.Background(), objectID, *hotel)
	if err != nil {
		return hotelDto, errors.NewInternalServerApiError("error when trying to update hotel", err)
	}

	hotelDto.ID = hotel.ID
	queue.Send(hotelDto.ID)
	return hotelDto, nil
}

func (s *hotelService) DeleteHotel(id string) errors.ApiError {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBadRequestApiError("invalid hotel ID")
	}

	err = daos.NewHotelDAO(s.db).DeleteHotel(context.Background(), objectID)
	if err != nil {
		return errors.NewInternalServerApiError("error when trying to delete hotel", err)
	}

	return nil
}

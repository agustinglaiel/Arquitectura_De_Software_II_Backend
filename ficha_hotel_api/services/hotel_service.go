package services

import (
	"bytes"
	"encoding/json"
	dao "ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
	"ficha_hotel_api/utils/queue"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelService struct{}

type HotelServiceInterface interface {
	GetHotelById(id string) (dtos.HotelDto, errors.ApiError)
	InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError)
	GetHotels() ([]dtos.HotelDto, errors.ApiError)
	DeleteHotelById(id string) errors.ApiError
}

var (
	HotelService HotelServiceInterface
)

func init() {
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
	//hotelDto.Photos = hotel.Photos
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
	//hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	hotel = dao.InsertHotel(hotel)

	if hotel.ID.Hex() == "000000000000000000000000" {
		return hotelDto, errors.NewBadRequestApiError("Error inserting new hotel")
	}

	hotelDto.ID = hotel.ID.Hex()

	queue.Send(hotelDto.ID, "INSERT")
	url := "http://localhost:8060/insertHotel"
	jsonData, err := json.Marshal(hotelDto)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)

	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	return hotelDto, nil
}

func (s *hotelService) UpdateHotelById(id string, hotelDto dtos.HotelDto) (dtos.HotelDto, errors.ApiError) {
	var hotel models.Hotel = dao.GetHotelById(id)

	if hotel.ID.Hex() == "000000000000000000000000" {
		return hotelDto, errors.NewBadRequestApiError("Hotel not found")
	}

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	//hotel.Photos = hotelDto.Photos
	hotel.Amenities = hotelDto.Amenities
	hotel.RoomCount = hotelDto.RoomCount
	hotel.City = hotelDto.City
	hotel.AvailableRooms = hotelDto.AvailableRooms

	dao.UpdateHotel(hotel)
	hotelDto.ID = hotel.ID.Hex()

	queue.Send(hotelDto.ID, "UPDATE")

	return hotelDto, nil
}

func (s *hotelService) GetHotels() ([]dtos.HotelDto, errors.ApiError) {
	hotels, err := dao.GetHotels()
	if err != nil {
		return nil, errors.NewInternalServerApiError("Error fetching hotels", err)
	}

	var hotelDtos []dtos.HotelDto
	for _, hotel := range hotels {
		hotelDto := dtos.HotelDto{
			ID:          hotel.ID.Hex(),
			Name:        hotel.Name,
			Description: hotel.Description,
			//Photos:         hotel.Photos,
			Amenities:      hotel.Amenities,
			RoomCount:      hotel.RoomCount,
			City:           hotel.City,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos = append(hotelDtos, hotelDto)
	}

	return hotelDtos, nil
}

func (s *hotelService) DeleteHotelById(id string) errors.ApiError {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBadRequestApiError("Invalid hotel ID")
	}

	err = dao.DeleteHotelById(objectID)
	if err != nil {
		return errors.NewInternalServerApiError("Error deleting hotel", err)
	}

	return nil
}

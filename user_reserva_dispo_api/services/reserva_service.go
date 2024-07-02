package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/utils/errors"
)

type reservaService struct{}

type reservaServiceInterface interface {
	ParseDate(dateString string) (time.Time, errors.ApiError)
	ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) (time.Time, errors.ApiError)
}

var (
	ReservaService reservaServiceInterface
)
var layout = "02-01-2006"

func init() {
	ReservaService = &reservaService{}
}

func (s *reservaService) ParseDate(dateString string) (time.Time, errors.ApiError) {
	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return date, errors.NewBadRequestApiError(err.Error())
	}

	return date, nil
}

func (s *reservaService) ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) (time.Time, errors.ApiError) {
	hotel, err := GetHotelById(reservationDto.HotelID)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Now(), errors.NewBadRequestApiError(err.Error())
	}
	log.Println(hotel.ID)
	reservas, err := daos.CheckAvailability(reservationDto.HotelID, reservationDto.StartDate, reservationDto.EndDate)
	log.Println(reservas)

	return time.Now(), nil
}

func GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
	var hotel dtos.HotelDto
	url := fmt.Sprintf("http://localhost:8080/hotel/%s", id)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	if err := json.Unmarshal(body, &hotel); err != nil {
		log.Println("Error unmarshalling hotel data:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	return hotel, nil

}

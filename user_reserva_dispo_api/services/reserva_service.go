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
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"
)

type reservaService struct{}

type reservaServiceInterface interface {
	ParseDate(dateString string) (time.Time, errors.ApiError)
	ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) errors.ApiError
	PostReserva(reservationDto dtos.CreateReservationDto) (dtos.CreateReservationDto, errors.ApiError)
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

func (s *reservaService) ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) errors.ApiError {
	hotel, err := GetHotelById(reservationDto.HotelID)

	if err != nil {
		fmt.Println("Error parsing date:", err)
		return errors.NewBadRequestApiError(err.Error())
	}

	reservas, err := daos.CheckAvailability(reservationDto.HotelID, reservationDto.StartDate, reservationDto.EndDate)
	if err != nil {
		return err
	}
	var listaDias []time.Time
	for i := reservationDto.StartDate; i.Before(reservationDto.EndDate) || i.Equal(reservationDto.EndDate); i = i.AddDate(0, 0, 1) {
		listaDias = append(listaDias, i)
	}

	return comprobar(reservas, hotel.RoomCount, listaDias)
}

func (s *reservaService) PostReserva(reservationDto dtos.CreateReservationDto) (dtos.CreateReservationDto, errors.ApiError) {
	var compReserva dtos.ReservationAvailabilityDto
	compReserva.StartDate = reservationDto.StartDate
	compReserva.EndDate = reservationDto.EndDate
	compReserva.HotelID = reservationDto.HotelID
	err := s.ComprobaDispReserva(compReserva)

	if err != nil {
		return reservationDto, err
	}
	var reserva models.Reservation
	reserva.HotelID = reservationDto.HotelID
	reserva.EndDate = reservationDto.EndDate
	reserva.StartDate = reservationDto.StartDate
	reserva.UserID = reservationDto.UserID

	reserva, err = daos.CreateReservation(reserva)

	if err != nil {
		return reservationDto, err
	}
	reservationDto.ReservationID = reserva.ReservationID

	return reservationDto, nil
}

func comprobar(reservas models.Reservations, habitacionesTotales int, listaDias []time.Time) errors.ApiError {
	conteoDias := make([]int, len(listaDias))
	for c, dia := range listaDias {
		for _, reserva := range reservas {
			if reserva.StartDate.Before(dia.AddDate(0, 0, -1)) && reserva.EndDate.After(dia.AddDate(0, 0, 1)) {
				conteoDias[c]++
				if conteoDias[c] >= habitacionesTotales {
					return errors.NewBadRequestApiError(fmt.Sprintf("El dia %d/%d/%d no hay disponibilidad", listaDias[c].Day(), int(listaDias[c].Month()), listaDias[c].Year()))

				}
			}
		}
	}

	return nil
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

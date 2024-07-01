package services

import (
	"fmt"
	"time"
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

	return time.Now(), nil
}

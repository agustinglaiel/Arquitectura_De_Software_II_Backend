package controllers

import (
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"
	"user_reserva_dispo_api/utils/errors"

	"github.com/gin-gonic/gin"
)

func ComprobaDispReserva(c *gin.Context) {
	var avalabilityReservationDto dtos.ReservationAvailabilityDto
	var err errors.ApiError
	c.Param("StartDate")
	c.Param("EndDate")
	c.Param("HotelID")

	avalabilityReservationDto.StartDate, err = services.ReservaService.ParseDate(c.Param("StartDate"))

	if err != nil {
		//poner error

	}

	avalabilityReservationDto.StartDate, err = services.ReservaService.ParseDate(c.Param("EndDate"))

	if err != nil {
		//poner error

	}
	avalabilityReservationDto.HotelID = c.Param("HotelID")
	if avalabilityReservationDto.HotelID == "" {
		//poner error
	}

	_, err = services.ReservaService.ComprobaDispReserva(avalabilityReservationDto)

}

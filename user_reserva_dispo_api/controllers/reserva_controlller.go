package controllers

import (
	"strconv"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"
	"user_reserva_dispo_api/utils/errors"

	"github.com/gin-gonic/gin"
)

func ComprobaDispReserva(c *gin.Context) {
	var avalabilityReservationDto dtos.ReservationAvailabilityDto
	var err errors.ApiError

	avalabilityReservationDto.StartDate, err = services.ReservaService.ParseDate(c.Param("StartDate"))
	if err != nil {
		c.JSON(400, err)
	}
	avalabilityReservationDto.StartDate, err = services.ReservaService.ParseDate(c.Param("EndDate"))
	if err != nil {
		c.JSON(400, err)
	}
	avalabilityReservationDto.HotelID = c.Param("HotelID")
	if avalabilityReservationDto.HotelID == "" {
		c.JSON(400, err)
	}
	err = services.ReservaService.ComprobaDispReserva(avalabilityReservationDto)
	if err != nil {
		c.JSON(400, err)

	}
	c.JSON(200, nil)

}

func PostReserva(c *gin.Context) {
	var avalabilityReservationDto dtos.CreateReservationDto
	var err errors.ApiError

	avalabilityReservationDto.StartDate, err = services.ReservaService.ParseDate(c.Param("StartDate"))
	if err != nil {
		c.JSON(400, err)
	}

	avalabilityReservationDto.EndDate, err = services.ReservaService.ParseDate(c.Param("EndDate"))
	if err != nil {
		c.JSON(400, err)
	}
	avalabilityReservationDto.HotelID = c.Param("HotelID")

	if err != nil {
		c.JSON(400, err)
	}
	userId, errr := strconv.Atoi(c.Param("idUser"))
	if errr != nil {
		c.JSON(400, err)
	}

	avalabilityReservationDto.UserID = userId

	services.ReservaService.PostReserva(avalabilityReservationDto)

}

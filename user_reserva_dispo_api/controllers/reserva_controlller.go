package controllers

import (
	"net/http"
	"strconv"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"
	"user_reserva_dispo_api/utils/errors"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func ComprobaDispReserva(c *gin.Context) {
	var avalabilityReservationDto dtos.ReservationAvailabilityDto
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

func Test(c *gin.Context) {
	// Crear una solicitud HTTP GET
	c.JSON(200, services.ReservaService.GetAmadeustoken())
}
func InsertHotel(c *gin.Context) {
	var HotelDto dtos.HotelDto
	err := c.BindJSON(&HotelDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	HotelDto, err = services.ReservaService.InsertHotel(HotelDto)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(200, HotelDto)

}

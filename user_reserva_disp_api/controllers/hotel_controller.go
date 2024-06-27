package controllers

import (
	"net/http"
	"strconv"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func GetHotelById(c *gin.Context){
	log.Debug("Hotel id to load: " + c.Param("id"))
	
	id, _ := strconv.Atoi(c.Param("id"))
	var hotelDto dtos.HotelDto

	hotelDto, err := services.HotelService.GetHotelById(id)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelDto)
}	

func GetHotels(c *gin.Context) {
	var hotelsDto dtos.HotelsDto
	hotelsDto, err := services.HotelService.GetHotels()

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDto)
}

//queda lo del insert desde mongo digamos
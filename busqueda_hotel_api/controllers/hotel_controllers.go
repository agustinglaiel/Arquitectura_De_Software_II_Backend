package controllers

import (
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/services"
	"busqueda_hotel_api/utils/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiter = make(chan bool, 3)
)

func GetHotelById(c *gin.Context) {
	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelDTO, err := services.HotelService.GetHotelById(id)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelDTO)
}

func InsertHotel(c *gin.Context) {
	var hotelDTO dtos.HotelDTO
	err := c.BindJSON(&hotelDTO)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hotelDTO, er := services.HotelService.InsertHotel(hotelDTO)

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusCreated, hotelDTO)
}

func UpdateHotelById(c *gin.Context) {
	var hotelDTO dtos.HotelDTO
	err := c.BindJSON(&hotelDTO)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	updatedHotelDTO, er := services.HotelService.UpdateHotelById(id, hotelDTO)
	<-rateLimiter

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusOK, updatedHotelDTO)
}

func DeleteHotelById(c *gin.Context) {
	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	err := services.HotelService.DeleteHotelById(id)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetHotels(c *gin.Context) {
	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetHotels()
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

func GetHotelsByCiudad(c *gin.Context) {
	ciudad := c.Param("ciudad")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetHotelsByCiudad(ciudad)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

func GetDisponibilidad(c *gin.Context) {
	ciudad := c.Param("ciudad")
	fechainicio := c.Param("fechainicio")
	fechafinal := c.Param("fechafinal")

	if ciudad == "" || fechainicio == "" || fechafinal == "" {
		apiErr := errors.NewBadRequestApiError("Missing parameters")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetDisponibilidad(ciudad, fechainicio, fechafinal)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

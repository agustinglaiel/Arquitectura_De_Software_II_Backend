package controllers

import (
	"ficha_hotel_api/dtos"
	service "ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"
	"ficha_hotel_api/utils/queue"
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
	hotelDto, err := service.HotelService.GetHotelById(id)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelDto)
}

func InsertHotel(c *gin.Context){
	var hotelDto dtos.HotelDto
	err := c.BindJSON(&hotelDto)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hotelDto, er := service.HotelService.InsertHotel(hotelDto)

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusCreated, hotelDto)
}

func UpdateHotelById(c *gin.Context){
	var hotelDto dtos.HotelDto
	err := c.BindJSON(&hotelDto)

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
	updatedHotelDto, er := service.HotelService.UpdateHotelById(id, hotelDto)
	<-rateLimiter

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusOK, updatedHotelDto)
}

func GetHotels(c *gin.Context) {
	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotels, err := service.HotelService.GetHotels()
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotels)
}

func DeleteHotelById(c *gin.Context) {
    id := c.Param("id")

    if len(rateLimiter) == cap(rateLimiter) {
        apiErr := errors.NewTooManyRequestsError("too many requests")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    rateLimiter <- true
    err := service.HotelService.DeleteHotelById(id)
    <-rateLimiter

    if err != nil {
        c.JSON(err.Status(), err)
        return
    }

    queue.Send(id, "delete")

    c.JSON(http.StatusNoContent, nil)
}
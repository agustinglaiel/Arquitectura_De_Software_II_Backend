package controllers

import (
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiter = make(chan bool, 3)
)

func GetHotelById(service services.HotelServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if len(rateLimiter) == cap(rateLimiter) {
			apiErr := errors.NewTooManyRequestsError("too many requests")
			c.JSON(apiErr.Status(), apiErr)
			return
		}

		rateLimiter <- true
		hotelDto, err := service.GetHotelById(id)
		<-rateLimiter

		if err != nil {
			c.JSON(err.Status(), err)
			return
		}

		c.JSON(http.StatusOK, hotelDto)
	}
}

func InsertHotel(service services.HotelServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var hotelDto dtos.HotelDto
		err := c.BindJSON(&hotelDto)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		hotelDto, er := service.InsertHotel(hotelDto)

		if er != nil {
			c.JSON(er.Status(), er)
			return
		}

		c.JSON(http.StatusCreated, hotelDto)
	}
}

func UpdateHotelById(service services.HotelServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		updatedHotelDto, er := service.UpdateHotelById(id, hotelDto)
		<-rateLimiter

		if er != nil {
			c.JSON(er.Status(), er)
			return
		}

		c.JSON(http.StatusOK, updatedHotelDto)
	}
}

func DeleteHotel(service services.HotelServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if len(rateLimiter) == cap(rateLimiter) {
			apiErr := errors.NewTooManyRequestsError("too many requests")
			c.JSON(apiErr.Status(), apiErr)
			return
		}

		rateLimiter <- true
		err := service.DeleteHotel(id)
		<-rateLimiter

		if err != nil {
			c.JSON(err.Status(), err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

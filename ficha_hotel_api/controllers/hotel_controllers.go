package controllers

import (
	"context"
	"net/http"

	"ficha_hotel_api/dtos"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	rateLimiter = make(chan bool, 3)
)

type HotelController struct {
	service services.HotelServiceInterface
}

func NewHotelController(service services.HotelServiceInterface) *HotelController {
	return &HotelController{service: service}
}

func (c *HotelController) CreateHotel(ctx *gin.Context) {
	var hotelDTO dtos.HotelDto
	if err := ctx.BindJSON(&hotelDTO); err != nil {
		apiErr := errors.NewBadRequestApiError("invalid JSON body")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	result, err := c.service.CreateHotel(context.Background(), hotelDTO)
	<-rateLimiter

	if err != nil {
		apiErr := errors.NewInternalServerApiError("error when trying to create hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *HotelController) DeleteHotel(ctx *gin.Context) {
	hotelID := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		apiErr := errors.NewBadRequestApiError("invalid hotel ID")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	err = c.service.DeleteHotel(context.Background(), id)
	<-rateLimiter

	if err != nil {
		apiErr := errors.NewInternalServerApiError("error when trying to delete hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *HotelController) UpdateHotel(ctx *gin.Context) {
	hotelID := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		apiErr := errors.NewBadRequestApiError("invalid hotel ID")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	var hotelDTO dtos.HotelDto
	if err := ctx.BindJSON(&hotelDTO); err != nil {
		apiErr := errors.NewBadRequestApiError("invalid JSON body")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	result, err := c.service.UpdateHotel(context.Background(), id, hotelDTO)
	<-rateLimiter

	if err != nil {
		apiErr := errors.NewInternalServerApiError("error when trying to update hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *HotelController) GetHotelByID(ctx *gin.Context) {
	hotelID := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		apiErr := errors.NewBadRequestApiError("invalid hotel ID")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotel, err := c.service.GetHotelByID(context.Background(), id)
	<-rateLimiter

	if err != nil {
		apiErr := errors.NewInternalServerApiError("error when trying to get hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, hotel)
}
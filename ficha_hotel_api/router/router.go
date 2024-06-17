package router

import (
	"ficha_hotel_api/controllers"

	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine, hotelController *controllers.HotelController) {
	router.GET("/hotel/:id", hotelController.GetHotelByID)
	router.POST("/hotel", hotelController.CreateHotel)
	router.DELETE("/hotel/:id", hotelController.DeleteHotel)
	router.PUT("/hotel/:id", hotelController.UpdateHotel)
}

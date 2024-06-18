package router

import (
	"ficha_hotel_api/controllers"
	"ficha_hotel_api/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine, hotelService services.HotelServiceInterface) {
	router.GET("/hotel/:id", controllers.GetHotelById(hotelService))
	router.POST("/hotel", controllers.InsertHotel(hotelService))
	router.PUT("/hotel/:id", controllers.UpdateHotelById(hotelService))
	router.DELETE("/hotel/:id", controllers.DeleteHotel(hotelService))

	fmt.Println("Finishing mappings configurations")
}

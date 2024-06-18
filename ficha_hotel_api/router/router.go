package router

import (
	"ficha_hotel_api/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine) {
	router.GET("/hotel/:id", controllers.GetHotelById)
	router.POST("/hotel", controllers.InsertHotel)
	router.PUT("/hotel/:id", controllers.UpdateHotelById)
	//router.DELETE("/hotel/:id", controllers.DeleteHotel)

	fmt.Println("Finishing mappings configurations")
}

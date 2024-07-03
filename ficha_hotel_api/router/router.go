package router

import (
	"ficha_hotel_api/controllers"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine) {

	router.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://localhost:3000", "https://mydomain.com"}, // Cambia esto a los orígenes que necesitas permitir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/hotel/:id", controllers.GetHotelById)
	router.POST("/hotel", controllers.InsertHotel)
	router.PUT("/hotel/:id", controllers.UpdateHotelById)
	router.DELETE("/hotel/:id", controllers.DeleteHotelById)
	router.GET("/hotels", controllers.GetHotels)

	fmt.Println("Finishing mappings configurations")
}

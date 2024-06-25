package router

import (
	"busqueda_hotel_api/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine) {
	// Rutas para hoteles
	router.GET("/hotel/:id", controllers.GetHotel)
	router.POST("/hotel", controllers.CreateHotel)
	router.PUT("/hotel/:id", controllers.UpdateHotel)
	//router.DELETE("/hotel/:id", controllers.DeleteHotelById)
	router.GET("/hotels", controllers.GetHotels)
	router.GET("/hotels/ciudad/:ciudad", controllers.GetHotelsByCiudad)
	router.GET("/hotels/disponibilidad/:ciudad/:fechainicio/:fechafinal", controllers.GetDisponibilidad)

	fmt.Println("Finishing mappings configurations")
}

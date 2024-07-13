package router

import (
	solrController "busqueda_hotel_api/controllers"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine) {

	//routerAdmin := router.Group("/admin")
	router.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://localhost:3000", "https://mydomain.com"}, // Cambia esto a los orígenes que necesitas permitir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	/*
		// Rutas para hoteles
		router.GET("/hotel/:id", controllers.GetHotel)
		router.POST("/hotel", controllers.CreateHotel)
		router.PUT("/hotel/:id", controllers.UpdateHotel)
		routerAdmin.DELETE("/hotel/:id", controllers.DeleteHotel)
		router.GET("/hotels", controllers.GetHotels)
		router.GET("/hotels/ciudad/:ciudad", controllers.GetHotelsByCiudad)
		router.GET("/hotels/disponibilidad/:ciudad/:fechainicio/:fechafinal", controllers.GetDisponibilidad)
	*/
	router.GET("/busqueda_hotel_api/search=:searchQuery", solrController.GetQuery)
	router.GET("/busqueda_hotel_api/searchAll=:searchQuery", solrController.GetQueryAllFields)
	router.GET("/ciudades", solrController.GetCiudades)
	router.GET("/hotelesByCiudad/:ciudad", solrController.GetHotelesByCiudad)
	fmt.Println("Finishing mappings configurations")
}

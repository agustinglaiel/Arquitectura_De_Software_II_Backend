package main

import (
	"ficha_hotel_api/controllers"
	"ficha_hotel_api/daos"
	"ficha_hotel_api/router"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar la base de datos
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.DisconnectDB()

	// Inicializar DAO
	hotelDAO := daos.NewHotelDAO(database)

	// Inicializar Servicio
	hotelService := services.NewHotelService(hotelDAO)

	// Inicializar Controlador
	hotelController := controllers.NewHotelController(hotelService)

	// Inicializar Router
	r := gin.Default()
	router.MapUrls(r, hotelController)
	r.Run(":8080")
}

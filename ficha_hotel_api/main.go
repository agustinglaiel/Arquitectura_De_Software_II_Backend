package main

import (
	"ficha_hotel_api/router"
	"ficha_hotel_api/services"
	"ficha_hotel_api/utils/db"
	"ficha_hotel_api/utils/queue"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ginRouter *gin.Engine
)

func main() {
	// Initialize RabbitMQ
	queue.Init()

	// Initialize the database
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Cannot init db")
		fmt.Println(err)
		return
	}
	defer db.DisconnectDB()

	// Initialize Hotel Service
	hotelService := services.NewHotelService(database)

	// Initialize the router
	ginRouter = gin.Default()
	router.MapUrls(ginRouter, hotelService)

	// Start the server
	fmt.Println("Starting server")
	if err := ginRouter.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}

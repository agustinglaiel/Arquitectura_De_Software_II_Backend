package main

import (
	"ficha_hotel_api/router"
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

	ginRouter = gin.Default()
	router.MapUrls(ginRouter)
	err := db.InitDB()
	defer db.DisconnectDB()

	if err != nil {
		fmt.Println("Cannot init Data base")
		fmt.Println(err)
		return
	}
	
	// Start the server
	fmt.Println("Starting server")
	ginRouter.Run(":8080")
}

package main

import (
	"busqueda_hotel_api/router"
	"busqueda_hotel_api/utils/db"
	"busqueda_hotel_api/utils/queue"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ginRouter *gin.Engine
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("Cannot init db")
		fmt.Println(err)
		return
	}
	defer db.DisconnectDB()

	ginRouter = gin.Default()
	router.MapUrls(ginRouter)

	go queue.InitConsumer()

	fmt.Println("Starting server")
	ginRouter.Run(":8080")
}

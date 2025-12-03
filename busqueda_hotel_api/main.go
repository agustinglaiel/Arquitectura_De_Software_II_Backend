package main

import (
	"busqueda_hotel_api/router"
	"busqueda_hotel_api/utils/queue"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Solr connection

	// Initialize the RabbitMQ consumer
	go queue.StartReceiving()

	// Create a new Gin router
	r := gin.Default()

	// Map the routes
	router.MapUrls(r)

	// Start the server
	port := "8070"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

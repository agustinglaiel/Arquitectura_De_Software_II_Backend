package main

import (
	"log"
	"user_reserva_dispo_api/router"
	"user_reserva_dispo_api/utils/db"

	"github.com/gin-gonic/gin"
)

func main() {

	// Inicializar la base de datos
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	db.StartDbEngine()
	// Crea el router con Gin
	r := gin.Default()

	// Mapa de URLs
	router.MapUrls(r)

	// Define el puerto y arranca el servidor
	port := ":8060"
	log.Printf("Starting server on port %s...\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package router

import (
	"fmt"
	"user_reserva_dispo_api/controllers"

	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine) {
	// User routes
	router.POST("/users/register", controllers.RegisterUser)  // Registro de usuarios
	router.POST("/users/login", controllers.LoginUser)        // Login de usuarios
	router.GET("/users/:id", controllers.GetUserById)         // Obtener usuario por ID
	router.GET("/users", controllers.GetUsers)                // Obtener todos los usuarios
	router.PUT("/users/:id", controllers.UpdateUser)          // Actualizar usuario
	router.DELETE("/users/:id", controllers.DeleteUser)       // Eliminar usuario

	fmt.Println("Finishing mappings configurations")
}

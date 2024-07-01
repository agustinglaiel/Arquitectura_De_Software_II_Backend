package router

import (
	"fmt"
	"user_reserva_dispo_api/controllers"
	"user_reserva_dispo_api/utils/auth"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica si el usuario es un administrador

func MapUrls(router *gin.Engine) {
	routerUsuario := router.Group("/usuario")
	routerUsuario.Use(auth.TokenMiddleware())

	routerAdmin := router.Group("/admin")
	routerAdmin.Use(auth.AdminTokenMiddleware())

	router.GET("/reserva/:StartDate/:EndDate/:HotelID", controllers.ComprobaDispReserva)
	router.POST("/users/register", controllers.RegisterUser) // Registro de usuarios
	router.POST("/users/login", controllers.LoginUser)       // Login de usuarios
	router.DELETE("/:id", controllers.DeleteUser)

	///////////Admin Rutas///////////////
	{
		routerAdmin.GET("/getUserById/:id", controllers.GetUserById)  // Obtener usuario por ID
		routerAdmin.GET("/getUsers", controllers.GetUsers)            // Obtener todos los usuarios
		routerAdmin.PUT("/updateUser/:id", controllers.UpdateUser)    // Actualizar usuario
		routerAdmin.DELETE("/deleteUser/:id", controllers.DeleteUser) // Eliminar usuario
		routerAdmin.HEAD("/")
	}

	fmt.Println("Finishing mappings configurations")
}

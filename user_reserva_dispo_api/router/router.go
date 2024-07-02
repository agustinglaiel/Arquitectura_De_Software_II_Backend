package router

import (
	"fmt"
	"time"
	"user_reserva_dispo_api/controllers"
	"user_reserva_dispo_api/utils/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica si el usuario es un administrador

func MapUrls(router *gin.Engine) {

	router.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://localhost:3000", "https://mydomain.com"}, // Cambia esto a los orígenes que necesitas permitir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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

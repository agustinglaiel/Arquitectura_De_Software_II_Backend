package router

import (
	"fmt"
	"net/http"
	"user_reserva_dispo_api/controllers"
	"user_reserva_dispo_api/utils/auth"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica si el usuario es un administrador
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if !auth.IsAdmin(token) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
            return
        }
        c.Next()
    }
}

func MapUrls(router *gin.Engine) {
	router.POST("/users/register", controllers.RegisterUser)  // Registro de usuarios
	router.POST("/users/login", controllers.LoginUser)        // Login de usuarios

	// Rutas protegidas para acciones de administrador
	adminRoutes := router.Group("/users")
	adminRoutes.Use(auth.AuthMiddleware(), AdminMiddleware())
	{
		adminRoutes.GET("/:id", controllers.GetUserById)         // Obtener usuario por ID
		adminRoutes.GET("/", controllers.GetUsers)               // Obtener todos los usuarios
		adminRoutes.PUT("/:id", controllers.UpdateUser)          // Actualizar usuario
		adminRoutes.DELETE("/:id", controllers.DeleteUser)       // Eliminar usuario
	}

	fmt.Println("Finishing mappings configurations")
}

package middelwers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")

	url := "https://api.ejemplo.com/endpoint"

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req.Header.Set("Authorization", tokenString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(http.StatusForbidden, "Servidor de autorizacion no disponible")
		return
	}
	if resp.StatusCode == http.StatusForbidden {
		ctx.JSON(http.StatusForbidden, "No tenes permiso de administrador")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusForbidden, "Token Invalido o inexistente")
	}
	if resp.StatusCode == http.StatusOK {
		ctx.Status(http.StatusOK)
	}
	if resp.StatusCode == http.StatusBadRequest {
		ctx.JSON(http.StatusForbidden, "Servidor de autorizacion no disponible")
	}
}

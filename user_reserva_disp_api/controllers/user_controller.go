package controllers

import (
	"log" // Usa el paquete de log estándar de Go
	"net/http"
	"os"
	"strconv"
	"time"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func InsertUser(c *gin.Context) {
	var userDto dtos.UserDto
	err := c.BindJSON(&userDto)

	if err != nil {
		log.Printf("Error parsing JSON params: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userDto, apiErr := services.UserService.InsertUser(userDto)
	if apiErr != nil {
		log.Printf("Insert error: %s", apiErr.Message())
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusCreated, userDto)
}

func GetUserById(c *gin.Context) {
	log.Printf("ID de usuario para cargar: %s", c.Param("id"))

	id, _ := strconv.Atoi(c.Param("id"))
	userDto, apiErr := services.UserService.GetUserById(id)

	if apiErr != nil {
		log.Printf("Error getting user by ID: %s", apiErr.Message())
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusOK, userDto)
}

func LoginUser(c *gin.Context) {
    var loginDto dtos.LoginRequestDto
    if err := c.BindJSON(&loginDto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    userDto, apiErr := services.UserService.AuthenticateUser(loginDto.Email, loginDto.Password)
    if apiErr != nil {
        c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
        return
    }

    token := generateToken(userDto)
    if token == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    response := dtos.LoginResponseDto{
        Token: token,
        ID:    userDto.ID,
        Type:  userDto.Type,
    }

    c.JSON(http.StatusOK, response)
}

func GetUserByEmail(c *gin.Context) {
    email := c.Param("email")
    if email == "" {
        log.Printf("No email provided for user retrieval")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
        return
    }

    userDto, apiErr := services.UserService.GetUserByEmail(email)
    if apiErr != nil {
        log.Printf("Error retrieving user by email: %s", apiErr.Message())
        c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
        return
    }

    token := generateToken(userDto)
    if token == "" {
        log.Printf("Failed to generate token after retrieving user by email")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    response := struct {
        Token   string       `json:"token"`
        Usuario dtos.UserDto `json:"usuario"`
    }{
        Token:   token,
        Usuario: userDto,
    }

    c.JSON(http.StatusOK, response)
}


func generateToken(userDto dtos.UserDto) string {
    claims := jwt.MapClaims{
        "id":        userDto.ID,
        "admin":     userDto.Type, // Type es un booleano que indica si el usuario es admin
        "exp":       time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        log.Printf("Failed to sign token: %s", err.Error())
        return ""
    }
    return tokenString
}




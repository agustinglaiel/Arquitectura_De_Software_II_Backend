package controllers

import (
	"log"
	"net/http"
	"strconv"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"

	"github.com/gin-gonic/gin"
)

// RegisterUser handles the user registration
func RegisterUser(c *gin.Context) {
	var userDto dtos.UserDto
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(userDto)
	result, err := services.UserService.RegisterUser(userDto)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var loginDto dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.UserService.LoginUser(loginDto.Username, loginDto.Password)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserById handles fetching a user by their ID
func GetUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	userDto, apiErr := services.UserService.GetUserById(userId)
	if apiErr != nil {
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusOK, userDto)
}

// GetUsers handles fetching all users
func GetUsers(c *gin.Context) {
	usersDto, apiErr := services.UserService.GetUsers()
	if apiErr != nil {
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusOK, usersDto)
}

// UpdateUser handles updating user data
func UpdateUser(c *gin.Context) {
	var userDto dtos.UserDto
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, apiErr := services.UserService.UpdateUser(userDto)
	if apiErr != nil {
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles the deletion of a user
func DeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	apiErr := services.UserService.DeleteUser(userId)
	if apiErr != nil {
		c.JSON(apiErr.Status(), gin.H{"error": apiErr.Message()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

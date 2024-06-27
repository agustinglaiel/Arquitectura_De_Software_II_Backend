package controllers

import (
	"net/http"
	"strconv"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/services"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func GetUserById(c *gin.Context){
	log.Debug("User id to load: " + c.Param("id"))

	id, _ := strconv.Atoi(c.Param("id"))
	var userDto dtos.UserDto

	userDto, err := services.UserService.GetUserById(id)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, userDto)
}

func GetUsers(c *gin.Context) {
	var usersDto dtos.UsersDto
	usersDto, err := services.UserService.GetUsers()

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, usersDto)
}

func UserInsert(c *gin.Context) {
	var userDto dtos.UserDto
	err := c.BindJSON(&userDto)

	// Error Parsing json param
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userDto, er := services.UserService.InsertUser(userDto)
	// Error del Insert
	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusCreated, userDto)
}

func Login(c *gin.Context) {
	var loginDto dtos.LoginDto
	er := c.BindJSON(&loginDto)

	if er != nil {
		log.Error(er.Error())
		c.JSON(http.StatusBadRequest, er.Error())
		return
	}
	log.Debug(loginDto)

	var loginResponseDto dtos.LoginResponseDto
	loginResponseDto, err := services.UserService.Login(loginDto)
	if err != nil {
		if err.Status() == 400 {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, loginResponseDto)
}

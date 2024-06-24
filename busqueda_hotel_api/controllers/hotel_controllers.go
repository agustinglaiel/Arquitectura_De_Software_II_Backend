package controllers

import (
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/services"
	"busqueda_hotel_api/utils/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiter = make(chan bool, 3)
)

func GetHotelById(c *gin.Context) {
	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelDTO, err := services.HotelService.GetHotelById(id)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelDTO)
}

func InsertHotel(c *gin.Context) {
	var hotelDTO dtos.HotelDTO
	err := c.BindJSON(&hotelDTO)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hotelDTO, er := services.HotelService.InsertHotel(hotelDTO)

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusCreated, hotelDTO)
}

func UpdateHotelById(c *gin.Context) {
	var hotelDTO dtos.HotelDTO
	err := c.BindJSON(&hotelDTO)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	updatedHotelDTO, er := services.HotelService.UpdateHotelById(id, hotelDTO)
	<-rateLimiter

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusOK, updatedHotelDTO)
}

func DeleteHotelById(c *gin.Context) {
	id := c.Param("id")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	err := services.HotelService.DeleteHotelById(id)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetHotels(c *gin.Context) {
	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetHotels()
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

func GetHotelsByCiudad(c *gin.Context) {
	ciudad := c.Param("ciudad")

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetHotelsByCiudad(ciudad)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

func GetDisponibilidad(c *gin.Context) {
	ciudad := c.Param("ciudad")
	fechainicio := c.Param("fechainicio")
	fechafinal := c.Param("fechafinal")

	if ciudad == "" || fechainicio == "" || fechafinal == "" {
		apiErr := errors.NewBadRequestApiError("Missing parameters")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if len(rateLimiter) == cap(rateLimiter) {
		apiErr := errors.NewTooManyRequestsError("too many requests")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	rateLimiter <- true
	hotelsDTO, err := services.HotelService.GetDisponibilidad(ciudad, fechainicio, fechafinal)
	<-rateLimiter

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, hotelsDTO)
}

func GetOrInsertByID(id string) {
	// Hago una request a hotel-api pidiendo todos los datos del hotel
	url := fmt.Sprintf("http://localhost:8070/hotel/%s", id)

	// Realiza la solicitud HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al hacer la solicitud HTTP:", err)
		return
	}
	defer resp.Body.Close()

	// Verifica si la respuesta fue exitosa (código 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("La solicitud no fue exitosa. Código de respuesta:", resp.StatusCode)
		return
	}

	// Lee el cuerpo de la respuesta HTTP
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta HTTP:", err)
		return
	}

	// Deserializa la respuesta en un objeto HotelDTO
	var hotelResponse dtos.HotelDTO
	if err := json.Unmarshal(body, &hotelResponse); err != nil {
		fmt.Println("Error al deserializar la respuesta:", err)
		return
	}

	// Me fijo si ya tengo cargado el hotel en solr
	_, err = services.HotelService.GetHotelById(id)
	if err != nil {
		// Si no lo tengo cargado entonces lo agrego
		_, err := services.HotelService.InsertHotel(hotelResponse)
		if err != nil {
			// Maneja el error de creación
			fmt.Println("Error al crear el hotel:", err)
			return
		}
		fmt.Println("Hotel nuevo agregado:", id)
		return
	}

	// Si ya lo tengo cargado, le hago el update
	_, err = services.HotelService.UpdateHotelById(id, hotelResponse)
	if err != nil {
		// Maneja el error de actualización
		fmt.Println("Error al actualizar el hotel:", err)
		return
	}
	return
}
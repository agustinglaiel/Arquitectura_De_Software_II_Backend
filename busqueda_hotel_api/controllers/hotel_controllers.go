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
    url := fmt.Sprintf("http://localhost:8080/hotel/%s", id)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error al hacer la solicitud HTTP:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("La solicitud no fue exitosa. Código de respuesta:", resp.StatusCode)
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error al leer la respuesta HTTP:", err)
        return
    }

    var hotelResponse dtos.HotelDTO
    if err := json.Unmarshal(body, &hotelResponse); err != nil {
        fmt.Println("Error al deserializar la respuesta:", err)
        return
    }

    hotelSolr, err := services.HotelService.GetHotelById(id)
    if err != nil {
        if apiErr, ok := err.(errors.ApiError); ok && apiErr.Status() == http.StatusNotFound {
            _, err := services.HotelService.InsertHotel(hotelResponse)
            if err != nil {
                fmt.Println("Error al crear el hotel:", err)
                return
            }
            fmt.Println("Hotel nuevo agregado:", id)
            return
        }
        fmt.Println("Error fetching hotel by id:", err)
        return
    }

    hotelSolr.Name = hotelResponse.Name
    hotelSolr.Description = hotelResponse.Description
    hotelSolr.City = hotelResponse.City
    hotelSolr.Photos = hotelResponse.Photos
    hotelSolr.Amenities = hotelResponse.Amenities
    hotelSolr.RoomCount = hotelResponse.RoomCount
    hotelSolr.AvailableRooms = hotelResponse.AvailableRooms

    _, err = services.HotelService.UpdateHotelById(id, hotelResponse)
    if err != nil {
        fmt.Println("Error al actualizar el hotel:", err)
        return
    }
    fmt.Println("Hotel actualizado:", id)
    return
}
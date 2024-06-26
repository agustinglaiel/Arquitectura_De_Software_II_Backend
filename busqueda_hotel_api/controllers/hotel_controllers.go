package controllers

import (
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrInsertByID(id string) {
    log.Printf("Recibido ID del hotel CONTROLLER: %s", id)
    url := fmt.Sprintf("http://localhost:8080/hotel/%s", id)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error al hacer la solicitud HTTP: %s", err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("La solicitud no fue exitosa. Código de respuesta: %d", resp.StatusCode)
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error al leer la respuesta HTTP: %s", err.Error())
        return
    }

    var hotelResponse dtos.HotelDTO
    if err := json.Unmarshal(body, &hotelResponse); err != nil {
        log.Printf("Error al deserializar la respuesta: %s", err.Error())
        return
    }

    log.Printf("Datos del hotel obtenidos de la API de ficha: %+v", hotelResponse)

    hotelSolr, err := services.HotelService.GetHotel(id)
    if err != nil {
        log.Printf("Error al obtener el hotel de Solr: %s", err.Error())
        _, err := services.HotelService.CreateHotel(hotelResponse)
		if err != nil {
			fmt.Println("Error al crear el hotel CONTROLLER: ", err)
			return
		}
		fmt.Println("Hotel nuevo agregado CONTROLLER con id: ", id)
		return
    }

    log.Printf("Hotel encontrado en Solr con ID: %s. Procediendo a actualizar.", id)
    hotelSolr.Name = hotelResponse.Name
    hotelSolr.Description = hotelResponse.Description
    hotelSolr.City = hotelResponse.City
    hotelSolr.Photos = hotelResponse.Photos
    hotelSolr.Amenities = hotelResponse.Amenities
    hotelSolr.RoomCount = hotelResponse.RoomCount
    hotelSolr.AvailableRooms = hotelResponse.AvailableRooms

    _, err = services.HotelService.UpdateHotel(hotelSolr)
    if err != nil {
        log.Printf("Error al actualizar el hotel en Solr CONTROLLER: %s", err.Error())
        return
    }
    log.Printf("Hotel actualizado en Solr CONTROLLER con ID: %s", id)
    return
}

func GetHotels(ctx *gin.Context) {
	hotelsDto, err := services.HotelService.GetAllHotels()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, hotelsDto)
}

func GetHotelsByCiudad(ctx *gin.Context) {
	ciudad := ctx.Param("ciudad")
	hotelsDto, err := services.HotelService.GetHotelsByCiudad(ciudad)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, hotelsDto)
}

func GetDisponibilidad(ctx *gin.Context) {
	ciudad := ctx.Param("ciudad")
	fechainicio := ctx.Param("fechainicio")
	fechafinal := ctx.Param("fechafinal")

	if ciudad == "" || fechainicio == "" || fechafinal == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Las fechas de inicio y final son obligatorias"})
		return
	}

	searchRequest := dtos.SearchRequestDTO{
		City:     ciudad,
		DateFrom: fechainicio,
		DateTo:   fechafinal,
	}

	hotelsDto, err := services.HotelService.GetDisponibilidad(searchRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, hotelsDto)
}


func GetHotel(ctx *gin.Context) {
	hotelID := ctx.Param("id")
	hotelDto, err := services.HotelService.GetHotel(hotelID)
	if err != nil {
		fmt.Println("Error al buscar el hotel:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Hotel no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, hotelDto)
}

func CreateHotel(ctx *gin.Context) {
	var hotelDto dtos.HotelDTO

	if err := ctx.ShouldBindJSON(&hotelDto); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Request Body"})
		return
	}

	hotel, err := services.HotelService.CreateHotel(hotelDto)
	if err != nil {
		fmt.Println("Error al crear el hotel:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el hotel"})
		return
	}

	ctx.JSON(http.StatusCreated, hotel)
}

func UpdateHotel(ctx *gin.Context) {
	hotelID := ctx.Param("id")

	existingHotel, err := services.HotelService.GetHotel(hotelID)
	if err != nil {
		fmt.Println("Error al buscar el hotel:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Hotel no encontrado"})
		return
	}

	var hotelDto dtos.HotelDTO
	if err := ctx.ShouldBindJSON(&hotelDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	existingHotel.Name = hotelDto.Name
	existingHotel.Description = hotelDto.Description
	existingHotel.City = hotelDto.City
	existingHotel.Photos = hotelDto.Photos
	existingHotel.Amenities = hotelDto.Amenities
	existingHotel.RoomCount = hotelDto.RoomCount
	existingHotel.AvailableRooms = hotelDto.AvailableRooms

	_, err = services.HotelService.UpdateHotel(existingHotel)
	if err != nil {
		fmt.Println("Error al actualizar el hotel:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el hotel"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Hotel actualizado con éxito"})
}

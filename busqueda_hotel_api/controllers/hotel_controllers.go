package controllers

import (
	"busqueda_hotel-api/dtos"
	service "busqueda_hotel-api/services"
	"busqueda_hotel-api/utils/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	//"busqueda_hotel-api/utils/validate"
)

func GetOrInsertByID(id string) {
	//Hago una request a hotel-api pidiendo todos los datos del hotel
	url := fmt.Sprintf("http://hotel-api:8000/hotel/%s", id)

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

	// Deserializa la respuesta en un objeto HotelDto
	var hotelResponse dtos.HotelDto
	if err := json.Unmarshal(body, &hotelResponse); err != nil {
		fmt.Println("Error al deserializar la respuesta:", err)
		return
	}

	//Me fijo si ya tengo cargado el hotel en solr
	hotelSolr, err := service.HotelService.GetHotel(id)
	if err != nil {
		// Si no lo tengo cargado entonces lo agrego
		_, err := service.HotelService.CreateHotel(hotelResponse)
		if err != nil {
			// Maneja el error de creación
			fmt.Println("Error al crear el hotel:", err)
			return
		}
		fmt.Println("Hotel nuevo agregado:", id)
		return
	}
	// Si ya lo tengo cargado, le hago el update
	// Actualiza los campos del hotel existente con los nuevos valores
	hotelSolr.Nombre = hotelResponse.Nombre
	hotelSolr.Descripcion = hotelResponse.Descripcion
	hotelSolr.Email = hotelResponse.Email
	hotelSolr.Ciudad = hotelResponse.Ciudad
	hotelSolr.Images = hotelResponse.Images
	hotelSolr.CantHab = hotelResponse.CantHab
	hotelSolr.Amenities = hotelResponse.Amenities

	// Actualiza el hotel en Solr
	_, err = service.HotelService.UpdateHotel(hotelSolr)
	if err != nil {
		// Maneja el error de actualización
		fmt.Println("Error al actualizar el hotel:", err)
		return
	}
	return
}

func GetHotels(ctx *gin.Context) {

	hotelsDto, err := service.HotelService.GetAllHotels()
	if err != nil {
		// Maneja el error de búsqueda del hotel
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusOK, hotelsDto)
	return
}

func GetHotelsByCiudad(ctx *gin.Context) {

	ciudad := ctx.Param("ciudad")
	hotelsDto, err := service.HotelService.GetHotelsByCiudad(ciudad)
	if err != nil {
		// Maneja el error de búsqueda del hotel
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusOK, hotelsDto)
	return
}

func GetDisponibilidad(ctx *gin.Context) {

	fechainicio := ctx.Param("fechainicio")
	fechafinal := ctx.Param("fechafinal")
	ciudad := ctx.Param("ciudad")

	// Validar que las fechas no estén vacías
	if fechainicio == "" || fechafinal == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Las fechas de check-in y check-out son obligatorias"})
		return
	}

	// Validar el formato de las fechas
	/*if !validate.IsValidDateFormat(checkin) || !validate.IsValidDateFormat(checkout) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El formato de las fechas debe ser 'YYYY-MM-DD'"})
		return
	}*/

	hotelsDto, err := service.HotelService.GetDisponibilidad(ciudad, fechainicio, fechafinal)
	if err != nil {
		// Maneja el error de búsqueda del hotel
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusOK, hotelsDto)
	return
}

func GetHotel(ctx *gin.Context) {
	// Obtiene el ID del hotel desde los parámetros de la URL
	hotelID := ctx.Param("id")

	// Verifica si el hotel existe antes de actualizarlo
	hotelDto, err := service.HotelService.GetHotel(hotelID)
	if err != nil {
		// Maneja el error de búsqueda del hotel
		fmt.Println("Error al buscar el hotel:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Hotel no encontrado"})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusOK, hotelDto)
	return
}

// CreateHotel maneja la solicitud POST para crear un nuevo hotel
func CreateHotel(ctx *gin.Context) {

	var hotelDto dtos.HotelDto

	if err := ctx.ShouldBindJSON(&hotelDto); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Request Body"})
		return
	}

	// Crea el hotel en Solr
	hotel, err := service.HotelService.CreateHotel(hotelDto)
	if err != nil {
		// Maneja el error de creación
		fmt.Println("Error al crear el hotel:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el hotel"})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusCreated, hotel)
	return
}

// UpdateHotel maneja la solicitud PUT para actualizar un hotel existente
func UpdateHotel(ctx *gin.Context) {
	// Obtiene el ID del hotel desde los parámetros de la URL
	hotelID := ctx.Param("id")

	// Verifica si el hotel existe antes de actualizarlo
	existingHotel, err := service.HotelService.GetHotel(hotelID)
	if err != nil {
		// Maneja el error de búsqueda del hotel
		fmt.Println("Error al buscar el hotel:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Hotel no encontrado"})
		return
	}

	var hotelDto dtos.HotelDto

	if err := ctx.ShouldBindJSON(&hotelDto); err != nil {
		errors.NewBadRequestApiError("Invalid request body")
		return
	}

	// Actualiza los campos del hotel existente con los nuevos valores
	existingHotel.Nombre = hotelDto.Nombre
	existingHotel.Descripcion = hotelDto.Descripcion
	existingHotel.Email = hotelDto.Email
	existingHotel.Ciudad = hotelDto.Ciudad
	existingHotel.Images = hotelDto.Images
	existingHotel.CantHab = hotelDto.CantHab
	existingHotel.Amenities = hotelDto.Amenities

	// Actualiza el hotel en Solr
	_, err = service.HotelService.UpdateHotel(existingHotel)
	if err != nil {
		// Maneja el error de actualización
		fmt.Println("Error al actualizar el hotel:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el hotel"})
		return
	}

	// Devuelve una respuesta de éxito
	ctx.JSON(http.StatusOK, gin.H{"message": "Hotel actualizado con éxito"})
}

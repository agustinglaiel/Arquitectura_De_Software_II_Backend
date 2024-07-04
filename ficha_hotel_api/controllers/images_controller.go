package controllers

import (
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/services"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertImage(c *gin.Context) {
    contentType := c.GetHeader("Content-Type")
    if contentType != "image/jpeg" && contentType != "image/png" {
        errMsg := "Content-Type debe ser image/jpeg o image/png"
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }

    imagenBytes, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo de la solicitud: " + err.Error()})
        return
    }

    hotelId, err := primitive.ObjectIDFromHex(c.Param("idHotel"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de hotel inválido: " + err.Error()})
        return
    }

    var img dtos.ImageDto
    img.Data = imagenBytes
    img.HotelId = hotelId
    img, err = services.ImageService.InsertImage(img)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error al insertar la imagen: " + err.Error()})
        return
    }
    
    // Devuelve un mensaje de éxito y el ID de la imagen
    c.JSON(http.StatusOK, gin.H{
        "message": "Imagen insertada con éxito",
        "imageId": img.HotelId.Hex(),  // Asumiendo que quieres devolver el ID de la imagen
    })
}

func GetImagesByHotelId(c *gin.Context) {
	hotelId, err := primitive.ObjectIDFromHex(c.Param("idHotel"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de hotel inválido: " + err.Error()})
		return
	}

	imagesDto, err := services.ImageService.GetImagesByHotelId(hotelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener imágenes por ID de hotel: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, imagesDto)
}
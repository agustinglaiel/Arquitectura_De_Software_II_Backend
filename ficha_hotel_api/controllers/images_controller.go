package controllers

import (
	"ficha_hotel_api/dtos"
	se "ficha_hotel_api/services"
	"ficha_hotel_api/utils/errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func InsertImage(ctx *gin.Context) {

	contentType := ctx.GetHeader("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		errMsg := "Content-Type debe ser image/jpeg o image/png"
		log.Error(errMsg)
		apiErr := errors.NewBadRequestApiError(errMsg)
		
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	imagenBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewBadRequestApiError("Error al leer el cuerpo de la solicitud")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	var img dtos.ImageDto
	img.Data = imagenBytes


	
	img.HotelId = ctx.Param("idHotel")

	img, err = se.ImageService.InsertImage(img)

	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewBadRequestApiError("Error al insertar la imagen",)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.Data(http.StatusOK, "image/jpeg", img.Data)

} //TOKEN ADMIN

func GetImagesByHotelId(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("idHotel"))
	if err != nil {
		errMsg := "ID de hotel inválido"
		log.Error(errMsg)
		apiErr := errors.NewBadRequestApiError(errMsg)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	imagesDto, err := se.ImageService.GetImagesByHotelId(id)

	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewInternalServerApiError("Error al obtener imágenes por ID de hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, imagesDto)

}
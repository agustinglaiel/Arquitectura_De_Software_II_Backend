package controllers

import (
	"busqueda_hotel_api/utils/errors"
	"ficha_hotel_api/dtos"
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
		apiErr := errors.NewBadRequestErrorApi(errMsg)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	imagenBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewInternalServerErrorApi("Error al leer el cuerpo de la solicitud", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	var img dtos.ImageDto
	img.Data = imagenBytes
	hotelId, err := strconv.Atoi(ctx.Param("idHotel"))

	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewInternalServerErrorApi("Error al convertir el ID del hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}
	img.HotelId = hotelId

	img, err = se.ImageService.InsertImage(img)

	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewInternalServerErrorApi("Error al insertar la imagen", err)
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
		apiErr := errors.NewBadRequestErrorApi(errMsg)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	imagesDto, err := HotelService.GetImagesByHotelId(id)

	if err != nil {
		log.Error(err.Error())
		apiErr := errors.NewInternalServerErrorApi("Error al obtener imágenes por ID de hotel", err)
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, imagesDto)

}
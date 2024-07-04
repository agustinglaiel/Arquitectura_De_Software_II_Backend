package services

import (
	"ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type imageService struct{}

type imageServicesInterface interface {
    InsertImage(im dtos.ImageDto) (dtos.ImageDto, errors.ApiError)
    GetImagesByHotelId(id primitive.ObjectID) (dtos.ImagesDto, errors.ApiError)
}

var (
	ImageService imageServicesInterface
)

func init() {
	ImageService = &imageService{}
}

func (*imageService) InsertImage(im dtos.ImageDto) (dtos.ImageDto, errors.ApiError) {
    // Convertir ObjectID a string para la llamada al servicio
    _, err := HotelService.GetHotelById(im.HotelId.Hex())
    if err != nil {
        return im, err
    }

    var mImage models.Image
    mImage.HotelID = im.HotelId // Ya es un ObjectID, no necesita conversión
    mImage.Imagen = im.Data

    mImage, err = daos.InsertImage(mImage)
    if err != nil {
        return im, err
    }

    return im, nil
}

func (*imageService) GetImagesByHotelId(id primitive.ObjectID) (dtos.ImagesDto, errors.ApiError) {
    images, err := daos.GetImagesByHotelId(id.Hex()) // Necesita ser convertido a string para pasar a DAO
    if err != nil {
        return dtos.ImagesDto{}, err
    }
    imagesList := make([]dtos.ImageDto, 0)
    for _, img := range images {
        var dto dtos.ImageDto
        dto.Data = img.Imagen
        dto.HotelId = img.HotelID // No necesita conversión a Hex ya que el DTO usa ObjectID
        imagesList = append(imagesList, dto)
    }
    return dtos.ImagesDto{
        Images: imagesList,
    }, nil
}
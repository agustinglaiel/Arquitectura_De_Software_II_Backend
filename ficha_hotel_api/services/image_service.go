package services

import (
	"ficha_hotel_api/daos"
	"ficha_hotel_api/dtos"
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"
)

type imageService struct{}

type imageServicesInterface interface {
	InsertImage(im dtos.ImageDto) (dtos.ImageDto, errors.ApiError)
	GetImagesByHotelId(id int) (dtos.ImagesDto, errors.ApiError)
}

var (
	ImageService imageServicesInterface
)

func init() {
	ImageService = &imageService{}
}

func (*imageService) InsertImage(im dtos.ImageDto) (dtos.ImageDto, errors.ApiError) {

	_, err := HotelService.GetHotelById(im.HotelId)
	if err != nil {
		return im, err
	}
	var mImage models.Image
	mImage.HotelID = im.HotelId
	mImage.Imagen = im.Data

	mImage, err = daos.InsertImage(mImage)
	if err != nil {
		return im, err
	}

	return im, nil
}

func (*imageService) GetImagesByHotelId(id int) (dtos.ImagesDto, errors.ApiError) {
	images, _ := daos.GetImagesByHotelId(id)
	imagesList := make([]dtos.ImageDto, 0)
	for _, im := range images {
		var dto dtos.ImageDto
		dto.Data = im.Imagen
		dto.HotelId = im.HotelID
		imagesList = append(imagesList, dto)
	}
	return dtos.ImagesDto{
		Images: imagesList,
	}, nil
}
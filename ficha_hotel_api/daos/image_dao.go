package daos

import (
	"ficha_hotel_api/models"
	"ficha_hotel_api/utils/errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
)



var Db *gorm.DB

func InsertImage(im models.Image) (models.Image, errors.ApiError) {
	img := Db.Create(&im)

	if img.Error != nil {
		return im, errors.NewBadRequestApiError("Error al insertar la imagen")
	}

	log.Debug("imagen para hotel. Id= ", im.HotelID)

	return im, nil
}

func GetImagesByHotelId(id int) (models.Images, errors.ApiError) {
	var imgs models.Images
	Db.Where("hotel_id = ?", id).Find(&imgs)
	if Db.Error != nil {
		return nil, errors.NewBadRequestApiError("No se pudo obtener el id")
	}

	return imgs, nil
}
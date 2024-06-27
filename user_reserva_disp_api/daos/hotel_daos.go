package daos

import (
	"user_reserva_dispo_api/models"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func GetHotelById(id int) models.Hotel{
	var hotel models.Hotel

	Db.Where("id = ?", id).First(&hotel)
	log.Debug("Hotel: ", hotel)

	return hotel
}

func GetHotelByIdMongo(id string) models.Hotel{
	var hotel models.Hotel

	Db.Where("id_mongo = ?", id).First(&hotel)
	log.Debug("Hotel: ", hotel)

	return hotel
}

func GetHotelByAmadeusId(idam string) bool {
	var hotel models.Hotel

	result := Db.Where("id_amadeus = ?", idam).First(&hotel)

	if result.Error != nil {
		return false 
	}

	return true //O sea ya existe un hotel con ese id en amadeus
}

func CheckHotelById(id int) bool {
	var hotel models.Hotel

	result := Db.Where("id = ?", id).First(&hotel)

	if result.Error != nil {
		return false
	}

	return true
}

func GetHotels() models.Hotels {
	var hotels models.Hotels
	Db.Find(&hotels)

	log.Debug("Hotels: ", hotels)

	return hotels
}

func InsertHotel(hotel models.Hotel) models.Hotel {
	result := Db.Create(&hotel)

	if result.Error != nil {
		//TODO Manage Errors
		log.Error("")
		hotel.Id = 0
	}
	log.Debug("Hotel Created: ", hotel.Id)
	return hotel
}
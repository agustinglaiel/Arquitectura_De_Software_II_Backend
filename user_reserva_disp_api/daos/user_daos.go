package daos

import (
	"user_reserva_dispo_api/models"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func GetUserByUsername(username string)(models.User, error){
	var user models.User
	result := Db.Where("user_name = ?", username).First(&user)

	log.Debug("User: ", user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func GetUserByEmail(email string) bool {
	var user models.User
	result := Db.Where("email = ?", email).First(&user)

	log.Debug("User: ", user)

	if result.Error != nil {
		return true // Si no lo encuntra es porque no existe
	}

	return false
}

func GetUserById(id int) models.User {
	var user models.User

	Db.Where("id = ?", id).First(&user)
	log.Debug("User: ", user)

	return user
}

func CheckUserById(id int) bool {
	var user models.User 

	// realza consulta a la base de datos: (con el id proporcionado como parametro)
	result := Db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return false
	}

	return true
}

func GetUsers() models.Users {
	var users models.Users
	Db.Find(&users)

	log.Debug("Users: ", users)

	return users
}

func InsertUser(user models.User) models.User {
	result := Db.Create(&user)

	if result.Error != nil {
		//TODO Manage Errors
		log.Error("")
	}
	log.Debug("User Created: ", user.Id)
	return user
}
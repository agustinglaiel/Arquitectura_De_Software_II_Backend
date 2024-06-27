package daos

import (
	"user_reserva_dispo_api/models"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
)


var Db *gorm.DB

func InsertUser(user models.User) models.User {
	result := Db.Create(&user)

	if result.Error != nil {
		log.Error("")
	}

	log.Debug("Usuario Creado: ", user.ID)
	return user
}

func GetUserById(id int) models.User {
	var user models.User

	Db.Where("id = ?", id).First(&user)
	log.Debug("Usuario: ", user)

	return user
}

/*func GetUserByUsername(username string) models.User {
	var user models.User

	Db.Where("user_name = ?", username).First(&user)
	log.Debug("Usuario: ", user)

	return user
}*/

func GetUserByEmail(email string) (models.User, error) {
    var user models.User
    result := Db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return models.User{}, result.Error // Retorna un usuario vacío y el error
    }
    return user, nil // Retorna el usuario encontrado y nil como error
}


func GetUsers() models.Users {
	var users models.Users

	Db.Find(&users)
	log.Debug("Usuarios: ", users)

	return users
}

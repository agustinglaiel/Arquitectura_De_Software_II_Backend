package daos

import (
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB //ACA VAMOS A TENER UN PROBLEMA (VAMOS A TENER QUE PONER EL DE RESERVA Y EL DE USUARIO EN CARPETAS SINO TIRA ERROR)

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := Db.Where("user_name = ?", username).First(&user)
	if result.Error != nil {
		log.Debug("Error retrieving user by username: %s", result.Error)
		return user, errors.NewApiError("User not found", "not_found", 404, nil)
	}
	return user, nil
}

func GetUserByEmail(email string) (bool, error) {
	var user models.User
	result := Db.Where("email = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return false, nil // No user found
	}
	if result.Error != nil {
		log.Debug("Error checking user email: %s", result.Error)
		return false, errors.NewInternalServerApiError("Database error", result.Error)
	}
	return true, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User
	result := Db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		log.Debug("Error retrieving user by ID: %s", result.Error)
		return user, errors.NewNotFoundApiError("User not found")
	}
	return user, nil
}


func CheckUserById(id int) (bool, error) {
	var user models.User
	result := Db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		log.Debug("Error checking user ID: %s", result.Error)
		return false, errors.NewApiError("User check failed", "check_failed", 400, nil)
	}
	return result.RowsAffected > 0, nil
}

func GetUsers() (models.Users, error) {
	var users models.Users
	result := Db.Find(&users)
	if result.Error != nil {
		log.Debug("Error retrieving users: %s", result.Error)
		return nil, errors.NewInternalServerApiError("Error fetching users", result.Error)
	}
	return users, nil
}

func InsertUser(user models.User) (models.User, error) {
	result := Db.Create(&user)
	if result.Error != nil {
		log.Debug("Error inserting user: %s", result.Error)
		return user, errors.NewBadRequestApiError("Error creating user")
	}
	return user, nil
}

func UpdateUser(user models.User) (models.User, error) {
    result := Db.Save(&user)
    if result.Error != nil {
        log.Debug("Error updating user: %s", result.Error)
        return user, errors.NewInternalServerApiError("Error updating user", result.Error)
    }
    return user, nil
}

func DeleteUser(userID int) error {
    result := Db.Delete(&models.User{}, userID)
    if result.Error != nil {
        log.Debug("Error deleting user: %s", result.Error)
        return errors.NewInternalServerApiError("Error deleting user", result.Error)
    }
    return nil
}

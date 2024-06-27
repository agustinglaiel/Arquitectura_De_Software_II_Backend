package services

import (
	"encoding/json"
	"fmt"
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/cache"
	"user_reserva_dispo_api/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct{}

type userServiceInterface interface{
	InsertUser(userDto dtos.UserDto)(dtos.UserDto, errors.ApiError)
	GetUserById(id int)(dtos.UserDto, errors.ApiError)
	//GetUserByUsername(username string)(dtos.UserDto, errors.ApiError)
	GetUserByEmail(email string)(dtos.UserDto, errors.ApiError)
	AuthenticateUser(email, password string) (dtos.UserDto, errors.ApiError)  // Asegúrate de añadir esta línea
}

var (
	UserService userServiceInterface
)

func init(){
	UserService = &userService{}
}

func (s *userService) InsertUser(userDto dtos.UserDto)(dtos.UserDto, errors.ApiError){
    var user models.User

    user.FirstName = userDto.FirstName
    user.LastName = userDto.LastName
    user.Email = userDto.Email

    // Hash de la contraseña
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
    if err != nil {
        return dtos.UserDto{}, errors.NewInternalServerApiError("Error al hacer hash de la contraseña", err)
    }
    user.Password = string(hashedPassword)

    user = daos.InsertUser(user)

	//Guardar en caché
	userBytes, _ := json.Marshal(userDto)
	cache.Set(fmt.Sprintf("user_%d", userDto.ID), userBytes)
	cache.Set(fmt.Sprintf("user_email_%s", userDto.Email), userBytes)
	fmt.Println("Usuario Guardado en Caché")

	return userDto, nil
}

func (s *userService) GetUserById(id int)(dtos.UserDto, errors.ApiError){
	cacheKey := fmt.Sprintf("user_%d", id)

	var cacheDTO dtos.UserDto
	cacheBytes := cache.Get(cacheKey)
	if cacheBytes != nil {
		fmt.Println("Usuario encontrado en caché")
		_ = json.Unmarshal(cacheBytes, &cacheDTO)
		return cacheDTO, nil
	}

	var user models.User = daos.GetUserById(id)
	var userDto dtos.UserDto

	if user.ID == 0 {
		return userDto, errors.NewBadRequestApiError("Usuario No Encontrado")
	}

	userDto.ID = user.ID
	userDto.FirstName = user.FirstName
	userDto.LastName = user.LastName
	userDto.Email = user.Email
	userDto.Password = user.Password

	//guardamos en caché
	userBytes, _ := json.Marshal(userDto)
	cache.Set(cacheKey, userBytes)
	fmt.Println("Usuario guardado en caché")

	return userDto, nil
}

func (s *userService) GetUserByEmail(email string) (dtos.UserDto, errors.ApiError) {

    // Genera una clave de caché única para el usuario
    cacheKey := fmt.Sprintf("user_email_%s", email)

    // Intentar obtener datos de la caché primero
    var cacheDTO dtos.UserDto
    cacheBytes := cache.Get(cacheKey)
    if cacheBytes != nil {
        fmt.Println("Found user in cache!")
        if err := json.Unmarshal(cacheBytes, &cacheDTO); err == nil {
            return cacheDTO, nil
        }
        // Si hay un error en deserialización, continua para refetch desde la base de datos
    }

    // Obtener el usuario de la base de datos
    user, err := daos.GetUserByEmail(email)
    if err != nil {
        return dtos.UserDto{}, errors.NewNotFoundApiError("Usuario no encontrado")
    }

    // Preparar DTO para la respuesta, excluyendo la contraseña por razones de seguridad
    userDto := dtos.UserDto{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Email:     user.Email,
        Type:      user.Type,
    }

    // Serializar y guardar en caché
    userBytes, err := json.Marshal(userDto)
    if err == nil {
        cache.Set(cacheKey, userBytes)
        fmt.Println("Saved user in cache!")
    } else {
        fmt.Println("Error marshalling user:", err)
    }

    return userDto, nil
}


func (s *userService) AuthenticateUser(email, password string) (dtos.UserDto, errors.ApiError) {
    user, err := daos.GetUserByEmail(email)
    if err != nil {
        return dtos.UserDto{}, errors.NewNotFoundApiError("Usuario no encontrado")
    }

    // Comprobamos la contraseña
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return dtos.UserDto{}, errors.NewBadRequestApiError("Contraseña incorrecta")
    }

    // Si la contraseña es correcta, devolvemos el DTO sin la contraseña
    return dtos.UserDto{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Email:     user.Email,
        Type:      user.Type,
    }, nil
}

/*
// Genera una clave de caché única para usuarios
func generateUserCacheKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}*/
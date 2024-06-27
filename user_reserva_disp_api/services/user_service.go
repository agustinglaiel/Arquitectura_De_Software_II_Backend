package services

import (
	"fmt"
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userService struct{}

type userServiceInterface interface {
	GetUsers()(dtos.UsersDto, errors.ApiError)
	InsertUser(userDto dtos.UserDto)(dtos.UserDto, errors.ApiError)
	GetUserById(id int)(dtos.UserDto, errors.ApiError)
	Login(loginDto dtos.LoginDto)(dtos.LoginResponseDto, errors.ApiError)
}

var(
	UserService userServiceInterface
)

func init() {
	UserService = &userService{}
}

func (s *userService) GetUserById(id int)(dtos.UserDto, errors.ApiError){
	var user models.User = daos.GetUserById(id)
	var userDto dtos.UserDto

	if user.Id == 0 {
		return userDto, errors.NewBadRequestApiError("Usuario no encontrado")
	}

	userDto.Name = user.Name
	userDto.LastName = user.LastName
	userDto.UserName = user.UserName
	userDto.Phone = user.Phone
	userDto.Address = user.Address
	userDto.Email = user.Email
	userDto.Id = user.Id
	userDto.Type = user.Type

	return userDto, nil
}

func (s *userService) GetUsers()(dtos.UsersDto, errors.ApiError){
	var users models.Users = daos.GetUsers()
	var usersDto dtos.UsersDto

	for _, user := range users {
		var userDto dtos.UserDto

		if !userDto.Type{
			userDto.Name = user.Name
			userDto.LastName = user.LastName
			userDto.UserName = user.UserName
			userDto.Phone = user.Phone
			userDto.Address = user.Address
			userDto.Email = user.Email
			userDto.Id = user.Id
			userDto.Type = user.Type
		}

		usersDto = append(usersDto, userDto)
	}

	return usersDto, nil
}

func (s *userService) InsertUser(userDto dtos.UserDto) (dtos.UserDto, errors.ApiError) {
	var user models.User

	if !daos.GetUserByEmail(userDto.Email){
		return userDto, errors.NewBadRequestApiError("El mail ya esta registrado")
	}

	user.Name = userDto.Name
	user.LastName = userDto.LastName
	user.UserName = userDto.UserName

	var hashedPassword, err = s.HashPassword(userDto.Password)

	if err != nil {
		return userDto, errors.NewBadRequestApiError("No se puede utilizar esa contraseña")
	}

	user.Password = hashedPassword //Ver como hasheo la pass
	user.Phone = userDto.Phone
	user.Address = userDto.Address
	user.Email = userDto.Email
	user.Type = userDto.Type

	user = daos.InsertUser(user)

	if user.Id == 0 {
		return userDto, errors.NewBadRequestApiError("Nombre de usuario no disponible")
	}

	userDto.Id = user.Id
	return userDto, nil
}

func (s *userService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("No se pudo hashear la password %w", err)
	}

	return string(hashedPassword), nil
}

func (s *userService) VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func (s *userService) Login(loginDto dtos.LoginDto) (dtos.LoginResponseDto, errors.ApiError){
	var user models.User
	user, err := daos.GetUserByUsername(loginDto.Username)
	var loginResponseDto dtos.LoginResponseDto
	loginResponseDto.UserId = -1

	if err != nil {
		return loginResponseDto, errors.NewBadRequestApiError("Usuario no existente")
	}

	var comparison error = s.VerifyPassword(user.Password, loginDto.Password)

	if loginDto.Username == user.UserName {
		if comparison != nil {
			return loginResponseDto, errors.NewUnauthorizedApiError("Contraseña incorrecta 2")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginDto.Username,
		"password": loginDto.Password,
	})
	var jwtKey = []byte("secret_key")
	tokenString, _ := token.SignedString(jwtKey)

	var verifyToken error = s.VerifyPassword(user.Password, tokenString)

	if loginDto.Username != user.UserName {
		if verifyToken != nil {
			return loginResponseDto, errors.NewUnauthorizedApiError("Contraseña incorrecta 3")
		}
	}

	loginResponseDto.UserId = user.Id
	loginResponseDto.Token = tokenString
	loginResponseDto.Name = user.Name
	loginResponseDto.LastName = user.LastName
	loginResponseDto.UserName = user.UserName
	loginResponseDto.Email = user.Email
	loginResponseDto.Type = user.Type
	log.Debug(loginResponseDto)
	
	return loginResponseDto, nil
}
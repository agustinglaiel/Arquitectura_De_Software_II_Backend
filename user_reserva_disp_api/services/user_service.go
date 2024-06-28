package services

import (
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	auth "user_reserva_dispo_api/utils/auth"
	"user_reserva_dispo_api/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct{}

type userServiceInterface interface {
	RegisterUser(userDto dtos.UserDto) (dtos.UserDto, errors.ApiError)
	LoginUser(username, password string) (dtos.LoginResponseDto, errors.ApiError) 
	GetUserById(userID int) (dtos.UserDto, errors.ApiError)
	GetUsers() (dtos.UsersDto, errors.ApiError)
	UpdateUser(userDto dtos.UserDto) (dtos.UserDto, errors.ApiError)
	DeleteUser(userID int) errors.ApiError
}

var(
	UserService userServiceInterface
)

func init() {
	UserService = &userService{}
}

func (s *userService) RegisterUser(userDto dtos.UserDto) (dtos.UserDto, errors.ApiError) {
	// Primero, verificamos si el correo electrónico ya está registrado
	exists, err := daos.GetUserByEmail(userDto.Email)
	if err != nil {
		return dtos.UserDto{}, errors.NewInternalServerApiError("Database error during email check", err)
	}
	if exists {
		return dtos.UserDto{}, errors.NewConflictApiError("Email already registered")
	}

	// Hasheamos la contraseña antes de guardarla
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.UserDto{}, errors.NewInternalServerApiError("Failed to hash password", err)
	}
	user := models.User{
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Username:  userDto.Username,
		Password:  string(hashedPassword),
		Email:     userDto.Email,
		Phone:     userDto.Phone,
		Address:   userDto.Address,
		IsAdmin:   userDto.Type,
	}

	// Insertamos el usuario en la base de datos
	registeredUser, err := daos.InsertUser(user)
	if err != nil {
		return dtos.UserDto{}, errors.NewInternalServerApiError("Failed to register user", err)
	}

	// Actualizamos el ID del DTO para reflejar el valor generado por la base de datos
	userDto.ID = registeredUser.ID
	return userDto, nil
}

func (s *userService) LoginUser(username, password string) (dtos.LoginResponseDto, errors.ApiError) {
    // Recuperar usuario por nombre de usuario
    user, err := daos.GetUserByUsername(username)
    if err != nil {
        return dtos.LoginResponseDto{}, errors.NewNotFoundApiError("User not found")
    }

    // Verificar la contraseña
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return dtos.LoginResponseDto{}, errors.NewUnauthorizedApiError("Invalid password")
    }

    // Generar token JWT
    token, err := auth.GenerateToken(user.ID, user.IsAdmin)
    if err != nil {
        return dtos.LoginResponseDto{}, errors.NewInternalServerApiError("Failed to generate token", err)
    }

    // Preparar y retornar la respuesta
    response := dtos.LoginResponseDto{
        UserID:    user.ID,
        Token:     token,
        Name:      user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
        Email:     user.Email,
        Type:      user.IsAdmin,
    }

    return response, nil
}

func (s *userService) GetUserById(userID int) (dtos.UserDto, errors.ApiError) {
    user, err := daos.GetUserById(userID)
    if err != nil {
        return dtos.UserDto{}, errors.NewNotFoundApiError("User not found")
    }

    userDto := dtos.UserDto{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
        Phone:     user.Phone,
        Address:   user.Address,
        Email:     user.Email,
        Type:      user.IsAdmin,
    }

    return userDto, nil
}

func (s *userService) GetUsers() (dtos.UsersDto, errors.ApiError) {
    users, err := daos.GetUsers()
    if err != nil {
        return nil, errors.NewInternalServerApiError("Error fetching users", err)
    }

    var usersDto dtos.UsersDto
    for _, user := range users {
        userDto := dtos.UserDto{
            ID:        user.ID,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Username:  user.Username,
            Phone:     user.Phone,
            Address:   user.Address,
            Email:     user.Email,
            Type:      user.IsAdmin,
        }
        usersDto = append(usersDto, userDto)
    }

    return usersDto, nil
}


func (s *userService) UpdateUser(userDto dtos.UserDto) (dtos.UserDto, errors.ApiError) {
    // Crear un modelo de usuario desde el DTO
    userToUpdate := models.User{
        ID:        userDto.ID,
        FirstName: userDto.FirstName,
        LastName:  userDto.LastName,
        Username:  userDto.Username,
        Password:  userDto.Password, // Considera hashear la contraseña antes de actualizar
        Email:     userDto.Email,
        Phone:     userDto.Phone,
        Address:   userDto.Address,
        IsAdmin:   userDto.Type,
    }

    // Actualizar usuario en la base de datos
    updatedUser, err := daos.UpdateUser(userToUpdate)
    if err != nil {
        return dtos.UserDto{}, errors.NewBadRequestApiError("Failed to update user")
    }

    // Actualizar el DTO con los valores del modelo actualizado
    updatedUserDto := dtos.UserDto{
        ID:        updatedUser.ID,
        FirstName: updatedUser.FirstName,
        LastName:  updatedUser.LastName,
        Username:  updatedUser.Username,
        Phone:     updatedUser.Phone,
        Address:   updatedUser.Address,
        Email:     updatedUser.Email,
        Type:      updatedUser.IsAdmin,
    }

    return updatedUserDto, nil
}


func (s *userService) DeleteUser(userID int) errors.ApiError {
    err := daos.DeleteUser(userID)
    if err != nil {
        return errors.NewBadRequestApiError("Failed to delete user")
    }
    return nil
}
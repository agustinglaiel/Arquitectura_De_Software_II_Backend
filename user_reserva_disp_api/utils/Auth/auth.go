package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key") // Utiliza una clave más segura y almacénala fuera del código

// GenerateToken genera un token JWT para un usuario
func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expirationTime.Unix(),
		// Puedes añadir más campos personalizados
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// VerifyToken verifica la validez del token proporcionado
func VerifyToken(tokenStr string) (*jwt.Token, error) {
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, err
}

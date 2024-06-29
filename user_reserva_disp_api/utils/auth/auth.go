package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key") // Utiliza una clave más segura y almacénala fuera del código

// GenerateToken genera un token JWT para un usuario
func GenerateToken(userID int, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expirationTime.Unix(),
		// Puedes añadir más campos personalizados aquí
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyToken verifica la validez del token proporcionado
func VerifyToken(tokenStr string) (*jwt.Token, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, err
}

// AuthMiddleware es un middleware de Gin para autenticar usando JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token, err := VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}

// IsAdmin verifica si el token pertenece a un administrador
func IsAdmin(tokenStr string) bool {
	token, err := VerifyToken(tokenStr)
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return false
	}
	// Aquí asumimos que el claim "isAdmin" es un booleano que indica si el usuario es administrador o no
	isAdmin, _ := strconv.ParseBool(claims.Issuer) // Ajustar según cómo se almacena este dato
	return isAdmin
}

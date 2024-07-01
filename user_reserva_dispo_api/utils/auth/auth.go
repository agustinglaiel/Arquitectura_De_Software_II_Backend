package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key") // Utiliza una clave más segura y almacénala fuera del código

// Claims standars + el es admin
type Claims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}

// GenerateToken genera un token JWT para un usuario
func GenerateToken(userID int, isAdmin bool) (string, error) {
	numAdmin := 0
	if isAdmin {
		numAdmin = 1
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": numAdmin,
		"fecha": time.Now().Unix(),
		"id":    userID,
	})

	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AdminTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("No token")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización faltante"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil // Clave secreta para verificar el token
		})

		if err != nil {
			log.Println("Invalido")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización inválido"})
			c.Abort()
			return
		}

		if token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				// Obtener la fecha de creación del token como Unix timestamp
				creationTime, ok := claims["fecha"].(float64)
				if !ok {
					log.Println("No fecha")

					c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener la fecha de creación del token"})
					c.Abort()
					return
				}

				// Convertir el Unix timestamp a time.Time
				creationDate := time.Unix(int64(creationTime), 0)

				// Verificar si ha pasado más de un día
				if time.Since(creationDate).Hours() <= 6 {
					isAdmin, ok := claims["admin"].(float64)
					if ok && isAdmin == 1 {
						text := strconv.FormatFloat(claims["id"].(float64), 'f', -1, 64)
						c.AddParam("idUser", text)
						c.AddParam("admin", "1")
						c.Next()

						c.Next()
					} else {
						log.Println("No sos admin")
						c.JSON(http.StatusForbidden, gin.H{"error": "Tienes que ser administrador para esta tarea"})
						c.Abort()
					}
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "El token ha expirado"})
					c.Abort()
				}
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización inválido"})
			c.Abort()
			return
		}
	}
}

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización faltante"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil // Clave secreta para verificar el token
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización inválido"})
			c.Abort()
			return
		}

		if token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				// Obtener la fecha de creación del token como Unix timestamp
				creationTime, ok := claims["fecha"].(float64)
				if !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener la fecha de creación del token"})
					c.Abort()
					return
				}

				// Convertir el Unix timestamp a time.Time
				creationDate := time.Unix(int64(creationTime), 0)

				// Verificar si ha pasado más de un día
				if time.Since(creationDate).Hours() <= 24 {
					text := strconv.FormatFloat(claims["id"].(float64), 'f', -1, 64)
					c.AddParam("idUser", text)
					c.Next()
					return
				}
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expiro"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización inválido"})
			c.Abort()
			return
		}
	}
}

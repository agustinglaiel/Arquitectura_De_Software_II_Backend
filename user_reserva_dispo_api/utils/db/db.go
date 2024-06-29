package db

import (
	"fmt"

	reservaDao "user_reserva_dispo_api/daos"
	userDao "user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	db  *gorm.DB
	err error
)

// InitDB inicializa la conexión a la base de datos.
func InitDB() error {
	// Parámetros de conexión a la base de datos
	DBName := "tpintegrador"
	DBUser := "tpintegrador"
	DBPass := "tpintegrador"
	DBHost := "localhost" // Cambia a la IP correcta si es necesario
	DBPort := "3307"      // Puerto mapeado del contenedor MySQL en el host

	// Formatea la cadena de conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DBUser, DBPass, DBHost, DBPort, DBName)

	// Abre la conexión a la base de datos
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
		return err
	}

	log.Println("Conexión establecida correctamente")

	// Asigna la conexión a los DAOs
	userDao.Db = db
	reservaDao.Db = db

	return nil
}

func StartDbEngine() {
	// We need to migrate all classes model.
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Reservation{})
	log.Info("Finalizacion de las tablas de la base de datos de migracion")
}

// CloseDB cierra la conexión de la base de datos.
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %s", err)
	}
	fmt.Println("Database connection closed.")
}

package db

import (
	"fmt"

	reservaDao "user_reserva_dispo_api/daos"
	userDao "user_reserva_dispo_api/daos"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	db  *gorm.DB
	err error
)

// InitDB inicializa la conexión a la base de datos.
func InitDB() {
	// DB Connections Paramters
	DBName := ""
	DBUser := ""
	DBPass := ""
	DBHost := ""
	// ------------------------

	db, err = gorm.Open("mysql", DBUser+":"+DBPass+"@tcp("+DBHost+":3306)/"+DBName+"?charset=utf8&parseTime=True")

	if err != nil{
		log.Info("La conexión no se pudo abrir")
		log.Fatal(err)
	} else {
		log.Info("Conexión establecida correctamente")
	}

	userDao.Db = db
	reservaDao.Db = db
}

// CloseDB cierra la conexión de la base de datos.
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %s", err)
	}
	fmt.Println("Database connection closed.")
}

package config

import (
	"fmt"
)

var (
	SOLRHOST       = "solr"
	SOLRPORT       = 8983
	SOLRCOLLECTION = "hotelSearch"

	HOTELSHOST = "ficha_hotel_api"
	HOTELSPORT = 8080

	QUEUENAME = "ficha_hotel-api"
	EXCHANGE  = "hotels"

	LBHOST = "lbbusqueda"
	LBPORT = 80

	RABBITUSER     = "guest"
	RABBITPASSWORD = "guest"
	RABBITHOST     = "rabbit"
	RABBITPORT     = 5672

	AMPQConnectionURL = fmt.Sprintf("amqp://%s:%s@%s:%d/", RABBITUSER, RABBITPASSWORD, RABBITHOST, RABBITPORT)

	USERAPIHOST = "user_reserva_dispo_api"
	USERAPIPORT = 8060
)

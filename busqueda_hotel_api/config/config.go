package config

import (
	"fmt"
)

var (
	SOLRHOST       = "localhost"
	SOLRPORT       = 8983
	SOLRCOLLECTION = "hotels"

	HOTELSHOST = "localhost"
	HOTELSPORT = 8080

	QUEUENAME = "ficha_hotel-api"
	EXCHANGE  = "hotels"

	LBHOST = "lbbusqueda"
	LBPORT = 80

	RABBITUSER     = "user"
	RABBITPASSWORD = "password"
	RABBITHOST     = "localhost"
	RABBITPORT     = 5672

	AMPQConnectionURL = fmt.Sprintf("amqp://%s:%s@%s:%d/", RABBITUSER, RABBITPASSWORD, RABBITHOST, RABBITPORT)

	USERAPIHOST = "user-res-api"
	USERAPIPORT = 8060
)

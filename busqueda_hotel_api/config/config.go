package config

import (
	"fmt"
)

var (
	SOLRHOST       = "solr"
	SOLRPORT       = 8983
	SOLRCOLLECTION = "hotels"

	HOTELSHOST = "localhost"
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

	USERAPIHOST = "user-res-api"
	USERAPIPORT = 8060
)

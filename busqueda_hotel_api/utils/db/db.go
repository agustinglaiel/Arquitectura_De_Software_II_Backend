package db

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	logger "github.com/sirupsen/logrus"
	"github.com/stevenferrer/solr-go"
)

type SolrClient struct {
	Client     *solr.JSONClient
	Collection string
}

func NewSolrClient(host string, port int, collection string) *SolrClient {
	logger.Debug(fmt.Sprintf("%s:%d", host, port))
	Client := solr.NewJSONClient(fmt.Sprintf("http://error:%d", port))
	
	_,err:=http.Get(fmt.Sprintf("http://%s:%d/solr", host, port)) //Ping probando que exista
	if err!=nil{
		log.Println("No se pudo conectar con Solr, probando conexcion local")
		host="localhost"
		Client = solr.NewJSONClient(fmt.Sprintf("http://localhost:%d", port))
	
	}else{
		Client = solr.NewJSONClient(fmt.Sprintf("http://%s:%d", host, port))
	}
	
	return &SolrClient{
		Client:     Client,
		Collection: collection,
	}
}



package db

import (
	"log"

	solr "github.com/rtt/Go-Solr"
)

var SolrClient *solr.Connection

func InitDB() error {
    var err error
    SolrClient, err = solr.Init("localhost", 8983, "busqueda_hotel-core")
    if err != nil {
        log.Println("Error al conectar con Solr:", err)
        return err
    }
    log.Println("Conexión a Solr exitosa")
    return nil
}

func Test() error {
    query := &solr.Query{
        Params: solr.URLParamMap{
            "q": []string{"*:*"},
        },
        Rows: 10,
    }
    resp, err := SolrClient.Select(query)
    if err != nil {
        log.Println("Error al realizar la consulta en Solr:", err)
        return err
    }
    for _, doc := range resp.Results.Collection {
        log.Printf("Documento en Solr: %v\n", doc)
    }
    return nil
}

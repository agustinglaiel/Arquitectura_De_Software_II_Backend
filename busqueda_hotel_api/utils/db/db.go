package db

import (
	"fmt"

	solr "github.com/rtt/Go-Solr"
)

var SolrClient *solr.Connection

func InitDB() error {
	var err error
	SolrClient, err = solr.Init("solr", 8983, "busqueda_hotel-core")
	if err != nil {
		return err
	}
	return nil
}

func Test() error {
	// Realiza la consulta a Solr
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{"*:*"}, // Consulta que selecciona todos los documentos
		},
		Rows: 10, // Número de filas a recuperar (ajusta según tus necesidades)
	}
	resp, err := SolrClient.Select(query)
	if err != nil {
		return err
	}

	// Itera a través de los resultados e imprímelos
	for _, doc := range resp.Results.Collection {
		// Aquí puedes acceder a los campos del documento y mostrar la información que necesites
		// Por ejemplo, si tienes un campo "nombre" en tus documentos Solr, puedes imprimirlo así:
		fmt.Printf("Nombre: %s\n", doc.Fields["test"])
		// Repite esto para otros campos que desees imprimir
	}
	return nil
}
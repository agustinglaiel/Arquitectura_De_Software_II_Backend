package daos

import (
	"busqueda_hotel_api/config"
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/utils/errors"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	logger "github.com/sirupsen/logrus"
	"github.com/stevenferrer/solr-go"
)

type SolrClient struct {
	Client     *solr.JSONClient
	Collection string
}

func (sc *SolrClient) GetCiudades() ([]string, error) {
	// Realizar la solicitud HTTP a Solr
	response, err := http.Get("http://localhost:8983/solr/hotels/select?facet=true&facet.field=city&q=*:*&rows=0&facet.limit=-1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()

	// Verificar el estado de la respuesta
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: estado de respuesta no válido. Código: %d", response.StatusCode)
		return nil, err
	}

	// Decodificar la respuesta JSON en un mapa genérico
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Acceder a los valores dentro del campo 'city'
	facetFields := result["facet_counts"].(map[string]interface{})["facet_fields"].(map[string]interface{})
	cityValues := facetFields["city"].([]interface{})

	// Iterar sobre los valores de 'city' y sus recuentos
	var cities []string
	for i := 0; i < len(cityValues); i += 2 {
		value := cityValues[i].(string) // Convertir a string
		cities = append(cities, value)
	}
	return cities, nil
}

func (sc *SolrClient) GetHotelCiudad(ciudad string) {
	log.Println(ciudad)
}
func (sc *SolrClient) GetQuery(query string, field string) (dtos.HotelsDTO, errors.ApiError) {
	var response dtos.SolrResponseDto
	var hotelsDto dtos.HotelsDTO
	q, err := http.Get(fmt.Sprintf("http://%s:%d/solr/hotels/select?q=%s%s%s", config.SOLRHOST, config.SOLRPORT, field, "%3A", query))
	log.Println(fmt.Sprintf("http://%s:%d/solr/hotels/select?q=%s%s%s", config.SOLRHOST, config.SOLRPORT, field, "%3A", query))
	if err != nil {
		return hotelsDto, errors.NewBadRequestApiError("Error getting from solr")
	}

	defer q.Body.Close()
	println(q.Body)
	err = json.NewDecoder(q.Body).Decode(&response)
	if err != nil {
		log.Printf("Response Body: %s", q.Body)
		log.Printf("Error: %s", err.Error())
		return hotelsDto, errors.NewBadRequestApiError("Error in unmarshal")
	}
	hotelsDto = response.Response.Docs
	log.Printf("Hotels: ", hotelsDto)
	return hotelsDto, nil
}
func (sc *SolrClient) GetQueryAllFields(query string) (dtos.HotelsDTO, errors.ApiError) {
	var response dtos.SolrResponseDto
	var hotelsDto dtos.HotelsDTO

	q, err := http.Get(fmt.Sprintf("http://%s:%d/solr/hotel/select?q=*:*", config.SOLRHOST, config.SOLRPORT))
	if err != nil {
		return hotelsDto, errors.NewBadRequestApiError("error getting from solr")
	}

	var body []byte
	body, err = io.ReadAll(q.Body)
	if err != nil {
		return hotelsDto, errors.NewBadRequestApiError("Error reading body")
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return hotelsDto, errors.NewBadRequestApiError("Error in unmarshal")
	}
	hotelsDto = response.Response.Docs
	return hotelsDto, nil
}

func (sc *SolrClient) Add(hotelDto dtos.Hotel2DTO) errors.ApiError {
	var addHotelDto dtos.AddDto
	addHotelDto.Add = dtos.DocDto{Doc: hotelDto}
	data, err := json.Marshal(addHotelDto)

	reader := bytes.NewReader(data)
	if err != nil {
		return errors.NewBadRequestApiError("Error getting json")
	}
	resp, err := sc.Client.Update(context.TODO(), sc.Collection, solr.JSON, reader)
	logger.Debug(resp)
	if err != nil {
		return errors.NewBadRequestApiError("Error in solr")
	}

	er := sc.Client.Commit(context.TODO(), sc.Collection)
	if er != nil {
		logger.Debug("Error committing load")
		return errors.NewInternalServerApiError("Error committing to solr", er)
	}
	return nil
}

func (sc *SolrClient) Delete(id string) errors.ApiError {
	var deleteDto dtos.DeleteDto
	deleteDto.Delete = dtos.DeleteDoc{Query: fmt.Sprintf("id:%s", id)}
	data, err := json.Marshal(deleteDto)
	reader := bytes.NewReader(data)
	if err != nil {
		return errors.NewBadRequestApiError("Error getting json")
	}
	resp, err := sc.Client.Update(context.TODO(), sc.Collection, solr.JSON, reader)
	logger.Debug(resp)
	if err != nil {
		return errors.NewBadRequestApiError("Error in solr")
	}
	er := sc.Client.Commit(context.TODO(), sc.Collection)
	if er != nil {
		logger.Debug("Error committing load")
		return errors.NewInternalServerApiError("Error committing to Solr", er)
	}
	return nil
}

/*
type HotelDao interface {
	Get(id string) (*models.Hotel, error)
	Create(hotel *models.Hotel) error
	Update(hotel *models.Hotel) error
	GetAll() ([]*models.Hotel, error)
	GetByCity(city string) ([]*models.Hotel, error)
}

type HotelSolrDao struct{}

func NewHotelSolrDAO() HotelDao {
	return &HotelSolrDao{}
}

func (dao *HotelSolrDao) Create(hotel *models.Hotel) error {
	hotelDocument := map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":             hotel.ID,
				"name":           hotel.Name,
				"description":    hotel.Description,
				"city":           hotel.City,
				"photos":         hotel.Photos,
				"room_count":     hotel.RoomCount,
				"amenities":      hotel.Amenities,
				"available_rooms": hotel.AvailableRooms,
			},
		},
	}

	_, err := db.SolrClient.Update(hotelDocument, true)
	if err != nil {
		log.Printf("Error al crear hotel en Solr: %s", err.Error())
		return err
	}
	log.Printf("Hotel creado en Solr: %+v", hotel)
	return nil
}

func (dao *HotelSolrDao) Update(hotel *models.Hotel) error {
	hotelDocument := map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":             hotel.ID,
				"name":           hotel.Name,
				"description":    hotel.Description,
				"city":           hotel.City,
				"photos":         hotel.Photos,
				"room_count":     hotel.RoomCount,
				"amenities":      hotel.Amenities,
				"available_rooms": hotel.AvailableRooms,
			},
		},
	}

	updateURL := "http://localhost:8983/solr/busqueda_hotel-core/update?commit=true"
	requestBody, err := json.Marshal(hotelDocument)
	if err != nil {
		return err
	}

	resp, err := http.Post(updateURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error al actualizar el hotel en Solr. Código de respuesta: %d", resp.StatusCode)
	}
	log.Printf("Hotel actualizado en Solr: %+v", hotel)
	return nil
}

func (dao *HotelSolrDao) Get(id string) (*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{fmt.Sprintf("id:%s", id)},
			"rows": []string{"1"},
		},
	}

	//log.Printf("Realizando consulta a Solr con ID: %s", id)

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		log.Printf("Error al realizar consulta a Solr: %s", err.Error())
		return nil, err
	}

	if len(resp.Results.Collection) == 0 {
		//log.Printf("No se encontró el hotel con ID %s en Solr", id)
		return nil, fmt.Errorf("No se encontró el hotel con id: %s en Solr", id)
	}

	doc := resp.Results.Collection[0]
	hotel := &models.Hotel{
		ID:             doc.Fields["id"].(string),
		Name:           getStringField(doc, "name"),
		Description:    getStringField(doc, "description"),
		City:           getStringField(doc, "city"),
		Photos:         getStringSliceFromInterface(doc.Field("photos")),
		RoomCount:      int(doc.Field("room_count").([]interface{})[0].(float64)),
		Amenities:      getStringSliceFromInterface(doc.Field("amenities")),
		AvailableRooms: int(doc.Field("available_rooms").([]interface{})[0].(float64)),
	}

	//log.Printf("Hotel obtenido de Solr: %+v", hotel)
	return hotel, nil
}

func (dao *HotelSolrDao) GetAll() ([]*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{"*:*"},
			"rows": []string{"1000"},
		},
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		return nil, err
	}

	var hotels []*models.Hotel
	for _, doc := range resp.Results.Collection {
		hotel := &models.Hotel{
			ID:             doc.Fields["id"].(string),
			Name:           getStringField(doc, "name"),
			Description:    getStringField(doc, "description"),
			City:           getStringField(doc, "city"),
			Photos:         getStringSliceFromInterface(doc.Field("photos")),
			RoomCount:      int(doc.Field("room_count").([]interface{})[0].(float64)),
			Amenities:      getStringSliceFromInterface(doc.Field("amenities")),
			AvailableRooms: int(doc.Field("available_rooms").([]interface{})[0].(float64)),
		}
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}

func (dao *HotelSolrDao) GetByCity(city string) ([]*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{fmt.Sprintf("city:\"%s\"", city)},
			"rows": []string{"1000"},
		},
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		return nil, err
	}

	var hotels []*models.Hotel
	for _, doc := range resp.Results.Collection {
		hotel := &models.Hotel{
			ID:             doc.Fields["id"].(string),
			Name:           getStringField(doc, "name"),
			Description:    getStringField(doc, "description"),
			City:           getStringField(doc, "city"),
			Photos:         getStringSliceFromInterface(doc.Field("photos")),
			RoomCount:      int(doc.Field("room_count").([]interface{})[0].(float64)),
			Amenities:      getStringSliceFromInterface(doc.Field("amenities")),
			AvailableRooms: int(doc.Field("available_rooms").([]interface{})[0].(float64)),
		}
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}

func getStringField(doc solr.Document, field string) string {
	if val, ok := doc.Field(field).([]interface{}); ok && len(val) > 0 {
		return val[0].(string)
	}
	return ""
}

func getStringSliceFromInterface(i interface{}) []string {
	var result []string
	if slice, ok := i.([]interface{}); ok {
		for _, v := range slice {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}
*/

package daos

import (
	"busqueda_hotel_api/models"
	"busqueda_hotel_api/utils/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	solr "github.com/rtt/Go-Solr"
)

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
		return err
	}
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

	return nil
}

func (dao *HotelSolrDao) Get(id string) (*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{fmt.Sprintf("id:%s", id)},
			"rows": []string{"1"},
		},
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		return nil, err
	}

	if len(resp.Results.Collection) == 0 {
		return nil, fmt.Errorf("hotel not found")
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

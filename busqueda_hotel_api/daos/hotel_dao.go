package daos

import (
	"busqueda_hotel_api/models"
	"busqueda_hotel_api/utils/db"
	"encoding/json"
	"fmt"
	"strings"

	solr "github.com/rtt/Go-Solr"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelDao interface {
	Get(id string) (*models.Hotel, error)
	Create(hotel *models.Hotel) error
	Update(hotel *models.Hotel) error
	GetAll() ([]*models.Hotel, error)
	GetByCiudad(city string) ([]*models.Hotel, error)
	GetDisponibilidad(ciudad, fechainicio, fechafinal string) ([]*models.Hotel, error)
	DeleteById(id string) error
}

type HotelSolrDao struct{}

func NewHotelSolrDao() HotelDao {
	return &HotelSolrDao{}
}

func (dao *HotelSolrDao) Get(id string) (*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{fmt.Sprintf("id:%s", id)},
		},
		Rows: 1,
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil || len(resp.Results.Collection) == 0 {
		fmt.Println("Error fetching hotel by id:", err)
		return nil, err
	}

	var hotel models.Hotel
	data, err := json.Marshal(resp.Results.Collection[0].Fields)
	if err != nil {
		fmt.Println("Error marshalling hotel data:", err)
		return nil, err
	}
	err = json.Unmarshal(data, &hotel)
	if err != nil {
		fmt.Println("Error unmarshalling hotel data:", err)
		return nil, err
	}

	return &hotel, nil
}

func (dao *HotelSolrDao) Create(hotel *models.Hotel) error {
	hotel.ID = primitive.NewObjectID().Hex()

	doc := solr.Document{
        Fields: map[string]interface{}{
            "id":              hotel.ID,
            "name":            hotel.Name,
            "description":     hotel.Description,
            "city":            hotel.City,
            "photos":          hotel.Photos,
            "amenities":       strings.Join(hotel.Amenities, ", "), // Convierte la lista a una cadena
            "room_count":      hotel.RoomCount,
            "available_rooms": hotel.AvailableRooms,
        },
    }

	_, err := db.SolrClient.Update(doc, true)
    if err != nil {
        fmt.Println("Error inserting hotel:", err)
        return models.Hotel{}
    }

    return hotel
}

func (dao *HotelSolrDao) Update(hotel *models.Hotel) error {
	doc := map[string]interface{}{
		"id":             hotel.ID,
		"name":           hotel.Name,
		"description":    hotel.Description,
		"city":           hotel.City,
		"photos":         hotel.Photos,
		"amenities":      hotel.Amenities,
		"room_count":     hotel.RoomCount,
		"available_rooms": hotel.AvailableRooms,
	}

	_, err := db.SolrClient.Update(doc, true)
	if err != nil {
		fmt.Println("Error updating hotel:", err)
		return err
	}

	return nil
}

func (dao *HotelSolrDao) GetAll() ([]*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{"*:*"},
		},
		Rows: 100,
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		fmt.Println("Error fetching all hotels:", err)
		return nil, err
	}

	var hotels []*models.Hotel
	for _, doc := range resp.Results.Collection {
		var hotel models.Hotel
		data, err := json.Marshal(doc.Fields)
		if err != nil {
			fmt.Println("Error marshalling hotel data:", err)
			return nil, err
		}
		err = json.Unmarshal(data, &hotel)
		if err != nil {
			fmt.Println("Error unmarshalling hotel data:", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (dao *HotelSolrDao) GetByCiudad(ciudad string) ([]*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{fmt.Sprintf("city:%s", ciudad)},
		},
		Rows: 100,
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		fmt.Println("Error fetching hotels by city:", err)
		return nil, err
	}

	var hotels []*models.Hotel
	for _, doc := range resp.Results.Collection {
		var hotel models.Hotel
		data, err := json.Marshal(doc.Fields)
		if err != nil {
			fmt.Println("Error marshalling hotel data:", err)
			return nil, err
		}
		err = json.Unmarshal(data, &hotel)
		if err != nil {
			fmt.Println("Error unmarshalling hotel data:", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (dao *HotelSolrDao) GetDisponibilidad(ciudad, fechainicio, fechafinal string) ([]*models.Hotel, error) {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{fmt.Sprintf("city:%s AND available_rooms:[1 TO *]", ciudad)},
		},
		Rows: 100,
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil {
		fmt.Println("Error fetching hotel availability:", err)
		return nil, err
	}

	var hotels []*models.Hotel
	for _, doc := range resp.Results.Collection {
		var hotel models.Hotel
		data, err := json.Marshal(doc.Fields)
		if err != nil {
			fmt.Println("Error marshalling hotel data:", err)
			return nil, err
		}
		err = json.Unmarshal(data, &hotel)
		if err != nil {
			fmt.Println("Error unmarshalling hotel data:", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (dao *HotelSolrDao) DeleteById(id string) error {
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q": []string{fmt.Sprintf("id:%s", id)},
		},
		Rows: 1,
	}

	resp, err := db.SolrClient.Select(query)
	if err != nil || len(resp.Results.Collection) == 0 {
		fmt.Println("Error fetching hotel by id for deletion:", err)
		return err
	}

	hotelID := resp.Results.Collection[0].Fields["id"].(string)
	_, err = db.SolrClient.Update(map[string]interface{}{
		"id":        hotelID,
		"_version_": -1, // Setting a negative version number to mark for deletion
	}, true)
	if err != nil {
		fmt.Println("Error deleting hotel:", err)
		return err
	}

	return nil
}

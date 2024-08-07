package services

import (
	"busqueda_hotel_api/config"
	client "busqueda_hotel_api/daos"
	"busqueda_hotel_api/dtos"

	"busqueda_hotel_api/utils/errors"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"fmt"
	"strconv"
	"strings"
)

type SolrService struct {
	solr *client.SolrClient
}

func NewSolrServiceImpl(solr *client.SolrClient) *SolrService {
	return &SolrService{
		solr: solr,
	}
}
func (s *SolrService) GetCiudades() ([]string, error) {
	ciudades, err := s.solr.GetCiudades()
	if err != nil {
		return nil, err
	}
	return ciudades, nil
}

func (s *SolrService) GetHotelesByCiudad(ciudad string) (dtos.HotelsDTO, error) {
	return nil, nil
}
func (s *SolrService) GetQuery(query string) (dtos.HotelsDTO, errors.ApiError) {
	var hotelsDto dtos.HotelsDTO
	queryParams := strings.Split(query, "_")
	numParams := len(queryParams)
	log.Printf("Paramas: %d", numParams)
	field, query := queryParams[0], queryParams[1]
	log.Printf("%s and %s", field, query)
	hotelsDto, err := s.solr.GetQuery(query, field)
	if err != nil {
		return hotelsDto, errors.NewBadRequestApiError("Solr failed")
	}

	if numParams == 4 {

		startdateQuery, enddateQuery := queryParams[2], queryParams[3]
		startdateSplit := strings.Split(startdateQuery, "-")
		enddateSplit := strings.Split(enddateQuery, "-")
		startdate := fmt.Sprintf("%s%s%s", startdateSplit[0], startdateSplit[1], startdateSplit[2])
		enddate := fmt.Sprintf("%s%s%s", enddateSplit[0], enddateSplit[1], enddateSplit[2])

		sDate, _ := strconv.Atoi(startdate)
		eDate, _ := strconv.Atoi(enddate)

		log.Debug(sDate)
		log.Debug(eDate)

		resultChan := make(chan dtos.HotelDTO, len(hotelsDto))

		var wg sync.WaitGroup
		var hotel dtos.HotelDTO

		for _, hotel = range hotelsDto {
			wg.Add(1)
			go func(hotel dtos.HotelDTO) {
				defer wg.Done()
				result, err := s.GetHotelInfo(hotel.ID, sDate, eDate)
				if err != nil {
					result = false
				}
				var response dtos.HotelDTO
				log.Debug("Adentro")
				log.Debug(result)
				log.Debug(response)

				if result == true {
					response = hotel
					resultChan <- response
				}
			}(hotel)
		}

		var hotelResults dtos.HotelsDTO
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for response := range resultChan {
			hotelResults = append(hotelResults, response)
		}

		return hotelResults, nil
	}

	return hotelsDto, nil
}

func (s *SolrService) GetHotelInfo(id string, startdate int, enddate int) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/user-res-api/hotel/availability/%s/%d/%d", config.USERAPIHOST, config.USERAPIPORT, id, startdate, enddate))
	if err != nil {
		return false, errors.NewBadRequestApiError("user_reserva_disp_api failed")
	}
	var body []byte
	body, _ = io.ReadAll(resp.Body)
	var responseDto dtos.AvailabilityResponse
	err = json.Unmarshal(body, &responseDto)

	if err != nil {
		log.Debugf("error in unmarshal")
		return false, errors.NewBadRequestApiError("Get hotel info failed")
	}
	status := responseDto.Status
	return status, nil
}

func (s *SolrService) GetQueryAllFields(query string) (dtos.HotelsDTO, errors.ApiError) {
	var hotelsDto dtos.HotelsDTO
	hotelsDto, err := s.solr.GetQueryAllFields(query)
	if err != nil {
		log.Debug(err)
		return hotelsDto, errors.NewBadRequestApiError("Solr Failed")
	}
	return hotelsDto, nil
}

func (s *SolrService) AddFromId(id string) errors.ApiError {
	var hotelDto dtos.Hotel2DTO
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/hotel/%s", config.HOTELSHOST, config.HOTELSPORT, id))

	if err != nil {
		log.Debugf("error getting item %s", id)
		return errors.NewBadRequestApiError("error getting hotel " + id)
	}
	log.Println("144")
	var body []byte
	body, err = io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &hotelDto)

	if err != nil {
		log.Println("151")
		log.Debugf("error in unmarshal of hotel %s", id)
		return errors.NewBadRequestApiError("error in unmarshal of hotel")
	}

	er := s.solr.Add(hotelDto)
	log.Debug(hotelDto)
	if er != nil {
		log.Println("158")

		log.Debugf("error adding to solr")
		return errors.NewInternalServerApiError("Adding to Solr error", err)
	}

	return nil
}

func (s *SolrService) Delete(id string) errors.ApiError {
	err := s.solr.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

/*type hotelService struct {
	dao daos.HotelDao
}

type HotelServiceInterface interface {
	GetHotel(id string) (dtos.HotelDTO, errors.ApiError)
	CreateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	UpdateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	GetAllHotels() (dtos.HotelsDTO, errors.ApiError)
	GetHotelsByCiudad(ciudad string) (dtos.HotelsDTO, errors.ApiError)
	GetDisponibilidad(searchRequest dtos.SearchRequestDTO) ([]dtos.SearchResultDTO, errors.ApiError)
}

var (
	HotelService HotelServiceInterface
)

func init() {
	HotelService = &hotelService{
		dao: daos.NewHotelSolrDAO(),
	}
}

func (s *hotelService) GetAllHotels() (dtos.HotelsDTO, errors.ApiError) {
	var hotelDtos dtos.HotelsDTO
	hotelDtos.Hotels = []dtos.HotelDTO{}
	hotels, err := s.dao.GetAll()
	if err != nil {
		return hotelDtos, errors.NewBadRequestApiError("error al obtener hoteles")
	}

	for _, hotel := range hotels {
		hotelDto := dtos.HotelDTO{
			ID:             hotel.ID,
			Name:           hotel.Name,
			Description:    hotel.Description,
			City:           hotel.City,
			Photos:         hotel.Photos,
			RoomCount:      hotel.RoomCount,
			Amenities:      hotel.Amenities,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos.Hotels = append(hotelDtos.Hotels, hotelDto)
	}

	return hotelDtos, nil
}

func (s *hotelService) GetHotelsByCiudad(ciudad string) (dtos.HotelsDTO, errors.ApiError) {
	var hotelDtos dtos.HotelsDTO
	hotelDtos.Hotels = []dtos.HotelDTO{}
	hotels, err := s.dao.GetByCity(ciudad)
	if err != nil {
		return hotelDtos, errors.NewBadRequestApiError("error al obtener hoteles")
	}

	for _, hotel := range hotels {
		hotelDto := dtos.HotelDTO{
			ID:             hotel.ID,
			Name:           hotel.Name,
			Description:    hotel.Description,
			City:           hotel.City,
			Photos:         hotel.Photos,
			RoomCount:      hotel.RoomCount,
			Amenities:      hotel.Amenities,
			AvailableRooms: hotel.AvailableRooms,
		}
		hotelDtos.Hotels = append(hotelDtos.Hotels, hotelDto)
	}

	return hotelDtos, nil
}

type DisponibilidadResult struct {
	HotelID        string
	Disponibilidad bool
}

func (s *hotelService) GetDisponibilidad(searchRequest dtos.SearchRequestDTO) ([]dtos.SearchResultDTO, errors.ApiError) {
	var searchResults []dtos.SearchResultDTO
	var hotels []*models.Hotel
	var err error

	if searchRequest.City == "" {
		hotels, err = s.dao.GetAll()
	} else {
		hotels, err = s.dao.GetByCity(searchRequest.City)
	}

	if err != nil {
		return searchResults, errors.NewBadRequestApiError("error al obtener hoteles")
	}

	disponibilidadCh := make(chan DisponibilidadResult, len(hotels))
	var wg sync.WaitGroup

	for _, hotel := range hotels {
		searchResult := dtos.SearchResultDTO{
			ID:          hotel.ID,
			Name:        hotel.Name,
			Description: hotel.Description,
			City:        hotel.City,
			Thumbnail:   hotel.Photos[0],
		}

		wg.Add(1)
		go func(hotel *models.Hotel, searchResult dtos.SearchResultDTO) {
			defer wg.Done()
			disponibilidad, err := checkDisponibilidad(hotel.ID, searchRequest.DateFrom, searchRequest.DateTo)
			if err != nil {
				disponibilidadCh <- DisponibilidadResult{HotelID: hotel.ID, Disponibilidad: false}
				return
			}
			disponibilidadCh <- DisponibilidadResult{HotelID: hotel.ID, Disponibilidad: disponibilidad}
		}(hotel, searchResult)

		searchResults = append(searchResults, searchResult)
	}

	wg.Wait()
	close(disponibilidadCh)

	disponibilidadMap := make(map[string]bool)
	for result := range disponibilidadCh {
		disponibilidadMap[result.HotelID] = result.Disponibilidad
	}

	for i, result := range searchResults {
		disponibilidad := disponibilidadMap[result.ID]
		searchResults[i].Availability = disponibilidad
	}

	return searchResults, nil
}

func checkDisponibilidad(hotelID string, fechainicio string, fechafinal string) (bool, error) {
	url := fmt.Sprintf("http://user-res-api:8002/hotel/%s/disponibilidad?fecha-inicio=%s&fecha-final=%s", hotelID, fechainicio, fechafinal)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("La solicitud de disponibilidad no fue exitosa. Código de respuesta: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var disponibilidadResponse struct {
		Disponibilidad bool `json:"disponibilidad"`
	}
	if err := json.Unmarshal(body, &disponibilidadResponse); err != nil {
		return false, err
	}

	return disponibilidadResponse.Disponibilidad, nil
}

func (s *hotelService) GetHotel(id string) (dtos.HotelDTO, errors.ApiError) {
	var hotelDto dtos.HotelDTO
	//log.Printf("Obteniendo hotel con ID: %s", id)

	hotel, err := s.dao.Get(id)
	if err != nil {
		log.Printf("Error al obtener hotel con ID %s: %s", id, err.Error())
		return hotelDto, errors.NewInternalServerApiError("Error fetching hotel", err)
	}

	if hotel == nil || hotel.ID == "" {
		log.Printf("No se encontró el hotel con ID %s", id)
		return hotelDto, errors.NewNotFoundApiError("Hotel not found")
	}

	hotelDto.ID = hotel.ID
	hotelDto.Name = hotel.Name
	hotelDto.Description = hotel.Description
	hotelDto.City = hotel.City
	hotelDto.Photos = hotel.Photos
	hotelDto.RoomCount = hotel.RoomCount
	hotelDto.Amenities = hotel.Amenities
	hotelDto.AvailableRooms = hotel.AvailableRooms

	//log.Printf("Hotel obtenido exitosamente con ID: %s", hotel.ID)
	return hotelDto, nil
}

func (s *hotelService) CreateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError) {
	var hotel models.Hotel

	hotel.ID = hotelDto.ID
	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.City = hotelDto.City
	hotel.Photos = hotelDto.Photos
	hotel.RoomCount = hotelDto.RoomCount
	hotel.Amenities = hotelDto.Amenities
	hotel.AvailableRooms = hotelDto.AvailableRooms

	//log.Printf("Intentando crear el hotel en Solr con datos: %+v", hotel)

	err := s.dao.Create(&hotel)
	if err != nil {
		log.Printf("Error al crear el hotel en Solr: %s", err.Error())
		return hotelDto, errors.NewBadRequestApiError(err.Error())
	}
	hotelDto.ID = hotel.ID

	log.Printf("Hotel creado exitosamente en Solr con ID: %s", hotel.ID)
	return hotelDto, nil
}

func (s *hotelService) UpdateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError) {
	var hotel models.Hotel

	hotel.ID = hotelDto.ID
	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.City = hotelDto.City
	hotel.Photos = hotelDto.Photos
	hotel.RoomCount = hotelDto.RoomCount
	hotel.Amenities = hotelDto.Amenities
	hotel.AvailableRooms = hotelDto.AvailableRooms

	err := s.dao.Update(&hotel)
	if err != nil {
		return hotelDto, errors.NewBadRequestApiError("error in update")
	}
	hotelDto.ID = hotel.ID

	return hotelDto, nil
}
*/

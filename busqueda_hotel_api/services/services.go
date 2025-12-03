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

var hostLocal string



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
	if hostLocal ==""{
		_, err := http.Get(fmt.Sprintf("http://%s:%d/hotels", config.HOTELSHOST, config.HOTELSPORT, id))
		if err!=nil{
			log.Println("Aca entre QLIA")

			hostLocal ="localhost"
		}else{
			hostLocal =config.HOTELSHOST
		}
	}
	
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/user-res-api/hotel/availability/%s/%d/%d", hostLocal, config.USERAPIPORT, id, startdate, enddate))
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
	if hostLocal ==""{
		_, err := http.Get(fmt.Sprintf("http://%s:%d/hotels", config.HOTELSHOST, config.HOTELSPORT, id))
		if err!=nil{
			log.Println("Aca entre QLIA")

			hostLocal ="localhost"
		}else{
			hostLocal =config.HOTELSHOST
		}
	}
	log.Printf(fmt.Sprintf("http://%s:%d/hotel/%s", hostLocal, config.HOTELSPORT, id))
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/hotel/%s", hostLocal, config.HOTELSPORT, id))

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


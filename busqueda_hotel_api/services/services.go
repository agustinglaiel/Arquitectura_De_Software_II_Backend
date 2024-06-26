package services

import (
	"busqueda_hotel_api/daos"
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/models"
	"busqueda_hotel_api/utils/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type hotelService struct {
	dao daos.HotelDao
}

type HotelServiceInterface interface {
	GetHotel(id string) (dtos.HotelDTO, errors.ApiError)
	CreateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	UpdateHotel(hotelDto dtos.HotelDTO) (dtos.HotelDTO, errors.ApiError)
	GetAllHotels() (dtos.HotelsDTO, errors.ApiError)
	GetHotelsByCiudad(ciudad string) (dtos.HotelsDTO, errors.ApiError)
	GetDisponibilidad(searchRequest dtos.SearchRequestDTO) ([]dtos.SearchResultDTO, errors.ApiError)
	DeleteHotel(id string) errors.ApiError
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
    log.Printf("Obteniendo hotel con ID: %s", id)

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

    log.Printf("Hotel obtenido exitosamente con ID: %s", hotel.ID)
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

    log.Printf("Intentando crear el hotel en Solr con datos: %+v", hotel)

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

func (s *hotelService) DeleteHotel(id string) errors.ApiError {
    err := s.dao.Delete(id)
    if err != nil {
        return errors.NewInternalServerApiError("Error deleting hotel from Solr", err)
    }
    return nil
}

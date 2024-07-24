package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"user_reserva_dispo_api/daos"
	"user_reserva_dispo_api/dtos"
	"user_reserva_dispo_api/models"
	"user_reserva_dispo_api/utils/errors"
)

type reservaService struct{}

type reservaServiceInterface interface {
	ParseDate(dateString string) (time.Time, errors.ApiError)
	ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) errors.ApiError
	PostReserva(reservationDto dtos.CreateReservationDto) (dtos.CreateReservationDto, errors.ApiError)
	GetAmadeustoken() string
}

var (
	ReservaService reservaServiceInterface
)
var layout = "02-01-2006"

func init() {
	ReservaService = &reservaService{}
}

func (s *reservaService) ParseDate(dateString string) (time.Time, errors.ApiError) {
	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return date, errors.NewBadRequestApiError(err.Error())
	}

	return date, nil
}

func (s *reservaService) ComprobaDispReserva(reservationDto dtos.ReservationAvailabilityDto) errors.ApiError {
	hotel, err := GetHotelById(reservationDto.HotelID)

	if err != nil {
		fmt.Println("Error parsing date:", err)
		return errors.NewBadRequestApiError(err.Error())
	}

	reservas, err := daos.CheckAvailability(reservationDto.HotelID, reservationDto.StartDate, reservationDto.EndDate)
	if err != nil {
		return err
	}
	var listaDias []time.Time
	for i := reservationDto.StartDate; i.Before(reservationDto.EndDate) || i.Equal(reservationDto.EndDate); i = i.AddDate(0, 0, 1) {
		listaDias = append(listaDias, i)
	}

	return comprobar(reservas, hotel.RoomCount, listaDias)
}

func (s *reservaService) PostReserva(reservationDto dtos.CreateReservationDto) (dtos.CreateReservationDto, errors.ApiError) {
	var compReserva dtos.ReservationAvailabilityDto
	compReserva.StartDate = reservationDto.StartDate
	compReserva.EndDate = reservationDto.EndDate
	compReserva.HotelID = reservationDto.HotelID
	err := s.ComprobaDispReserva(compReserva)

	if err != nil {
		return reservationDto, err
	}
	var reserva models.Reservation
	reserva.HotelID = reservationDto.HotelID
	reserva.EndDate = reservationDto.EndDate
	reserva.StartDate = reservationDto.StartDate
	reserva.UserID = reservationDto.UserID

	reserva, err = daos.CreateReservation(reserva)

	if err != nil {
		return reservationDto, err
	}
	reservationDto.ReservationID = reserva.ReservationID

	return reservationDto, nil
}

func comprobar(reservas models.Reservations, habitacionesTotales int, listaDias []time.Time) errors.ApiError {
	conteoDias := make([]int, len(listaDias))
	for c, dia := range listaDias {
		for _, reserva := range reservas {
			if reserva.StartDate.Before(dia.AddDate(0, 0, -1)) && reserva.EndDate.After(dia.AddDate(0, 0, 1)) {
				conteoDias[c]++
				if conteoDias[c] >= habitacionesTotales {
					return errors.NewBadRequestApiError(fmt.Sprintf("El dia %d/%d/%d no hay disponibilidad", listaDias[c].Day(), int(listaDias[c].Month()), listaDias[c].Year()))

				}
			}
		}
	}

	return nil
}

func GetHotelById(id string) (dtos.HotelDto, errors.ApiError) {
	var hotel dtos.HotelDto
	url := fmt.Sprintf("http://localhost:8080/hotel/%s", id)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	if err := json.Unmarshal(body, &hotel); err != nil {
		log.Println("Error unmarshalling hotel data:", err)
		return hotel, errors.NewBadRequestApiError(err.Error())
	}

	return hotel, nil

}

// funcion para generar un token de amadeus cada vez que voy a hacer la consulta
func (s *reservaService) GetAmadeustoken() string {

	fmt.Printf("entro al f d token")
	// Define los datos que deseas enviar en el cuerpo de la solicitud.
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "xKPS4fEY8izDKbYvllOLZAWJ5eYPDzBh")
	data.Set("client_secret", "jGGh1F3xjwLBXQ7Q")

	// Realiza la solicitud POST a la API externa.
	resp, err := http.Post("https://test.api.amadeus.com/v1/security/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return ""
	}
	defer resp.Body.Close()

	/* // Custom HTTP client with TLS configuration to skip certificate verification.
		customTransport := &http.Transport{
	        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	    }

	    client := &http.Client{
	        Transport: customTransport,
	    } */

	/*    // Create a new request.
		req, err := http.NewRequest("POST", "https://test.api.amadeus.com/v1/security/oauth2/token", strings.NewReader(data.Encode()))
		   if err != nil {
			   fmt.Println("Error al crear la solicitud:", err)
			   return ""
		   }
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform the request.
	    resp, err := client.Do(req)

		if err != nil {
	        fmt.Println("Error al hacer la solicitud:", err)
	        return ""
	    }*/

	// Lee la respuesta de la API.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return ""
	}

	// Parsea la respuesta JSON para obtener el token (asumiendo que la respuesta es JSON).
	// Si la respuesta es en otro formato, ajusta esto en consecuencia.
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return ""
	}
	token, ok := response["access_token"].(string)
	if !ok {
		return ""
	}
	fmt.Println("token:", token)
	return token

}

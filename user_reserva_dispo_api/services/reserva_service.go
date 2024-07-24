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
	InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, error)
	GetAmadeustoken() string
	AvailabilityAmadeus(startdateconguiones string, enddateconguiones string, idAm string) bool
}

var (
	ReservaService reservaServiceInterface
)
var layout = "02-01-2006"
var layoutAmadeus = "2006-01-02"

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

	h, er := daos.GetHotelById(reservationDto.HotelID)
	if er != nil {
		return errors.NewBadRequestApiError(er.Error())
	}
	if s.AvailabilityAmadeus(reservationDto.StartDate.Format(layoutAmadeus), reservationDto.EndDate.Format(layoutAmadeus), h.IdAmadeus) {

		return comprobar(reservas, hotel.RoomCount, listaDias)
	}
	return errors.NewBadRequestApiError("No hay disponibilidad en esa fecha")

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
	return token
}

func (s *reservaService) InsertHotel(hotelDto dtos.HotelDto) (dtos.HotelDto, error) {

	apiUrl := "https://test.api.amadeus.com/v1/reference-data/locations/hotels/by-city?cityCode=MIA&radius=5&radiusUnit=KM&hotelSource=ALL"
	// Crear una solicitud HTTP GET
	solicitud, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud dentro de insert hotel:", err)
		return hotelDto, err
	}

	token := s.GetAmadeustoken()
	solicitud.Header.Set("Authorization", "Bearer "+token)

	// Realiza la solicitud HTTP
	cliente := &http.Client{}

	respuesta, err := cliente.Do(solicitud)
	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return hotelDto, err

	}

	defer respuesta.Body.Close()

	// Leer y manejar la respuesta de la API externa
	var response struct {
		Data []struct {
			HotelID string `json:"hotelId"`
		} `json:"data"`
	}

	// Decodificar la respuesta JSON
	decoder := json.NewDecoder(respuesta.Body)
	if err := decoder.Decode(&response); err != nil {
		return hotelDto, err

	}
	log.Println(response)
	for _, hotel := range response.Data {
		fmt.Printf("Id amadeus: %s\n", hotel.HotelID)
		hotelM, err := daos.CheckHotelExists(hotel.HotelID)
		if err != nil {
			fmt.Println(err)
			return hotelDto, err
		} else if !hotelM {
			hotelDto.IdAmadeus = hotel.HotelID
			_, er := daos.InsertHotel(hotelDto)
			// Error del Insert
			if er != nil {
				return hotelDto, er
			}

			break // Se encontró el ID, sal del bucle
		}
	}

	return hotelDto, nil

}

func (s *reservaService) AvailabilityAmadeus(startdateconguiones string, enddateconguiones string, idAm string) bool {
	fmt.Println("entro a availability")
	apiUrl := "https://test.api.amadeus.com/v3/shopping/hotel-offers"
	apiUrl += "?hotelIds=HHMIA500"

	apiUrl += "&checkInDate=" + startdateconguiones
	apiUrl += "&checkOutDate=" + enddateconguiones

	fmt.Println(apiUrl)

	// Crear una solicitud HTTP
	solicitud, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("ERROR CREANDO SOLICITUD:", err)
		return false
	}

	// Agregar el encabezado de autorización Bearer con tu token
	token := s.GetAmadeustoken()
	solicitud.Header.Set("Authorization", "Bearer "+token)

	fmt.Println(solicitud)

	cliente := &http.Client{}

	respuesta, err := cliente.Do(solicitud)

	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return false
	}

	// Verifica el código de estado de la respuesta
	if respuesta.StatusCode != http.StatusOK {
		//fmt.Printf("La solicitud a la API de Amadeus no fue exitosa. Código de estado: %d\n", respuesta)
		return true
	}
	defer respuesta.Body.Close() // Mover defer aquí para cerrar el cuerpo correctamente

	// Lee el cuerpo de la respuesta
	responseBody, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return false
	}

	// Crear una estructura para deserializar el JSON de la respuesta
	var responseStruct struct {
		Data []struct {
			Type                   string `json:"type"`
			ID                     string `json:"id"`
			ProviderConfirmationID string `json:"providerConfirmationId"`
		} `json:"data"`
	}

	// Decodificar el JSON y extraer el campo "id"
	if err := json.Unmarshal(responseBody, &responseStruct); err != nil {
		fmt.Println("Error al decodificar el JSON de la respuesta:", err)
		return false
	}

	// Obtén el ID del hotel del primer elemento en "data"
	if len(responseStruct.Data) > 0 {
		// si el largo de la respuesta es mayor q cero es pq hay disponibilidad --> llamo al service
		fmt.Println("Amadeus nos dice que hay disponibilidad")
		return true
	}

	fmt.Println("No hay disponibilidad en esas fechas")
	return false
}

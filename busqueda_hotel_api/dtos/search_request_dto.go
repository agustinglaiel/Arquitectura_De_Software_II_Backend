package dtos

type SearchRequestDTO struct {
    City      string `json:"city"`
    DateFrom  string `json:"date_from"`
    DateTo    string `json:"date_to"`
}

/*
Para recibir los parámetros de búsqueda desde el frontend.
Campos: ciudad, fecha desde, fecha hasta.
*/
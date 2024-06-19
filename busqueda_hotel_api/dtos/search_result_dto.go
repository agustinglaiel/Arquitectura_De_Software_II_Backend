package dtos

type SearchResultDTO struct {
    ID           string `json:"id"`
    Name         string `json:"name"`
    Description  string `json:"description"`
    Thumbnail    string `json:"thumbnail"`
    City         string `json:"city"`
    Availability bool   `json:"availability"`
}

/*
Para devolver los resultados de la búsqueda al frontend.
Campos: ID del hotel, nombre, descripción, miniatura, ciudad, disponibilidad.
*/
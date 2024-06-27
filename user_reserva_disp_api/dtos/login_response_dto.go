package dtos

type LoginResponseDto struct {
    Token string `json:"token"`
    ID    int    `json:"id"`
    Type  bool   `json:"type"` // True para admin
}
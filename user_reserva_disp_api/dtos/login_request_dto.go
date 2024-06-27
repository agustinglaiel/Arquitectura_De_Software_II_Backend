package dtos

type LoginRequestDto struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
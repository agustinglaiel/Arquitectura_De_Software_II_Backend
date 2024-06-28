package dtos

type UserDto struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Phone     string    `json:"phone"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Type      bool   `json:"type"` // True para admin
}

type UsersDto []UserDto

type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	UserID    int    `json:"user_id"`
	Token     string `json:"token"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Type      bool   `json:"type"` // True para admin
}

